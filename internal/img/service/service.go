package service

import (
	"scaffold/internal/common/reskit/codes"
	"scaffold/internal/img/domain"
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"time"
)

type service struct {
	repo		domain.ImgRepository
	s3Client	*s3.Client
	msgQueue	domain.ImgMsgQueue
	publicBucket	bucket
	deleteBucket	bucket
}

func loadS3() (*s3.Client, string, string) {
	// 加载环境变量
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	publicBucket := os.Getenv("R2_PUBLIC_BUCKET_NAME")
	deleteBucket := os.Getenv("R2_DELETE_BUCKET_NAME")

	if accountID == "" || accessKeyID == "" || secretAccessKey == "" || publicBucket == "" || deleteBucket == "" {
		log.Fatal("Missing one or more R2 environment variables (R2_ACCOUNT_ID, R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_PUBLIC_BUCKET_NAME, R2_DELETE_BUCKET_NAME)")
	}

	// 配置 S3 客户端
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(
		service, region string, options ...interface{},
	) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithRegion("auto"),	// R2 不使用区域，但 SDK 需要
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	return s3.NewFromConfig(cfg), publicBucket, deleteBucket
}

func NewImgService(repo domain.ImgRepository, msgQueue domain.ImgMsgQueue) domain.ImgService {
	client, publicBucket, deleteBucket := loadS3()

	return &service{
		repo:		repo,
		s3Client:	client,
		publicBucket:	bucket(publicBucket),
		deleteBucket:	bucket(deleteBucket),
		msgQueue:	msgQueue,
	}
}

const compressQuality = 60

// Compress 压缩图片质量，返回压缩后的图片数据
func (s *service) Compress(src io.Reader) (io.Reader, error) {
	// 注册 PNG 解码器
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	// 解码图片
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("解码输入图片失败: %v", err)
	}

	// 创建一个缓冲区来存储压缩后的图片数据
	output := &bytes.Buffer{}

	// 将图片编码为 JPEG 格式并以指定质量写入缓冲区
	err = jpeg.Encode(output, img, &jpeg.Options{Quality: compressQuality})
	if err != nil {
		return nil, fmt.Errorf("编码图片失败: %v", err)
	}

	return bytes.NewReader(output.Bytes()), nil
}

func (s *service) genNewPath(oldPath string) string {
	return fmt.Sprintf("deleted/%d_%s", time.Now().Unix(), oldPath)
}

func (s *service) Upload(src io.Reader, img *domain.Img, categoryID int64) (*domain.Img, error) {
	// 压缩图片
	compressed, err := s.Compress(src)
	if err != nil {
		return nil, codes.ErrImgCompress.WithCause(err)
	}

	// 1. 若有 categoryID 则需要先检查分类是否存在
	var category *domain.Category
	if categoryID != 0 {
		var err error
		category, err = s.repo.FindCategoryByID(categoryID)
		if err != nil {
			return nil, err
		}
	}

	// 2.查询是否有相同路径
	nowPath := ""
	if category != nil {
		nowPath = category.Prefix + "/" + img.Path
	}
	exist, err := s.repo.ExistByPath(nowPath)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, codes.ErrImgPathRepeat
	}

	// 3.入库
	res, err := s.repo.Create(img, categoryID)
	if err != nil {
		return nil, err
	}

	// 后续不要再使用 img 使用res！
	// 4.上传s3
	uploadOk := true
	if err = s.UploadFile(compressed, res.Path); err != nil {
		uploadOk = false
		err = codes.ErrImgUploadToR3Failed.WithCause(err)
	}

	// 5.如果第4步发生错误 则删除已入库的记录
	if !uploadOk {
		if err := s.repo.Delete(res.ID, true); err != nil {
			zap.L().Error("数据库入库成功但图片上传失败，尝试回滚删除数据库记录时出错",
				zap.Int64("id:", res.ID),
				zap.Int64("category_id:", categoryID),
				zap.String("path", res.Path),
				zap.Error(err),
			)
		}
	}

	return res, err
}

// Delete 删除逻辑
// 硬删除 -> 直接删除 publicBucket 中的对象
// 软删除 -> 复制原有对象到不可公共访问的 deleteBucket 删除 publicBucket 中的对象 -> 类似于回收站功能
func (s *service) Delete(id int64, hard ...bool) error {
	img, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	var ifHard bool

	if len(hard) == 0 {
		ifHard = false
	} else if hard[0] == true {
		ifHard = true
	}

	if ifHard {
		// 1.删除r3
		if err := s.DeleteFile(img.Path, s.publicBucket); err != nil {
			return err
		}
		// 2.删除记录
		if err := s.repo.Delete(img.ID, true); err != nil {
			return err
		}
	} else {
		// 1.从 publicBucket 移到 deleteBucket
		if err := s.CopyFileToAnotherBucket(img.Path, s.publicBucket, s.deleteBucket); err != nil {
			return err
		}
		//2.删除 publicBucket 中的对象
		if err := s.DeleteFile(img.Path, s.publicBucket); err != nil {
			return err
		}

		// 3.软删除记录
		if err := s.repo.Delete(img.ID, false); err != nil {
			return err
		}

		// 4.将id记录到消息队列
		if err := s.msgQueue.AddToDeleteQueue(img.ID); err != nil {
			zap.L().Error("图片软删除：添加到定时删除队列失败",
				zap.Int64("imgID", img.ID),
				zap.String("path", img.Path),
				zap.Error(err),
			)
		}
	}

	return nil
}

// ClearRecycleBin 删除被软删除的数据
// 此时删除 deleteBucket对象 数据库记录 消息队列key
func (s *service) ClearRecycleBin(id int64) error {
	// 1.先查询图片信息
	img, err := s.repo.FindByID(id, true)
	if err != nil {
		return err
	}

	// 2.删除 deleteBucket 中的文件
	if err := s.DeleteFile(img.Path, s.deleteBucket); err != nil {
		return err
	}

	// 3.硬删除数据库记录
	if err := s.repo.Delete(id, true); err != nil {
		return err
	}

	// 4.从删除队列中移除（防止定时器重复删除）
	if err := s.msgQueue.RemoveFromDeleteQueue(id); err != nil {
		zap.L().Error("清空回收站：移除定时删除任务失败",
			zap.Int64("imgID", id),
			zap.Error(err),
		)
	}

	return nil
}

func (s *service) List(query *domain.ImgQuery) (*domain.ImgPages, error) {
	return s.repo.List(query)
}

func (s *service) ListenDeleteQueue() {
	s.msgQueue.ListenDeleteQueue(func(imgID int64) {
		//1.先查询img
		img, err := s.repo.FindByID(imgID, true)
		if err != nil {
			zap.L().Error("定时删除队列：查询图片失败",
				zap.Int64("imgID", imgID),
				zap.Error(err),
			)
			return
		}

		//2.删除R3
		if err := s.repo.Delete(imgID, true); err != nil {
			zap.L().Error("定时删除队列：删除数据库记录失败",
				zap.Int64("imgID", imgID),
				zap.String("path", img.Path),
				zap.Error(err),
			)
			return
		}

		//	3.删除记录
		if err := s.DeleteFile(img.Path, s.deleteBucket); err != nil {
			zap.L().Error("定时删除队列：删除存储文件失败",
				zap.Int64("imgID", imgID),
				zap.String("path", img.Path),
				zap.String("bucket", s.deleteBucket.string()),
				zap.Error(err),
			)
		}

		zap.L().Info("定时删除队列：图片删除成功",
			zap.Int64("imgID", imgID),
			zap.String("path", img.Path),
		)
	})
}

func (s *service) RestoreFromRecycleBin(id int64) (*domain.Img, error) {
	// 1.查询已软删除的图片信息
	img, err := s.repo.FindByID(id, true)
	if err != nil {
		return nil, err
	}

	// 2.从 deleteBucket 复制回 publicBucket
	if err := s.CopyFileToAnotherBucket(img.Path, s.deleteBucket, s.publicBucket); err != nil {
		return nil, err
	}

	// 3.删除 deleteBucket 中的文件
	if err := s.DeleteFile(img.Path, s.deleteBucket); err != nil {
		return nil, err
	}

	// 4.恢复数据库记录（取消软删除）
	res, err := s.repo.Restore(id)
	if err != nil {
		return nil, err
	}

	// 5.从删除队列中移除
	if err := s.msgQueue.RemoveFromDeleteQueue(id); err != nil {
		zap.L().Error("从回收站恢复：移除定时删除任务失败",
			zap.Int64("imgID", id),
			zap.Error(err),
		)
	}

	return res, nil
}

const maxCategory = 15

func (s *service) CreateCategory(category *domain.Category) (*domain.Category, error) {
	exist, err := s.repo.CategoryExistByTitle(category.Title)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if exist {
		return nil, codes.ErrImgCategoryTitleRepeat
	}

	count, err := s.repo.CountCategory()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if count >= maxCategory {
		return nil, codes.ErrImgCategoryToMany
	}

	return s.repo.CreateCategory(category)
}

func (s *service) isCategoryExistImg(id int64) error {
	existing, err := s.repo.IsCategoryExistImg(id)
	if err != nil {
		return err
	}

	if existing {
		return codes.ErrImgCategoryExistImg
	}
	return nil
}

func (s *service) UpdateCategory(category *domain.Category) (*domain.Category, error) {
	// 1.除开自己以外 是否有与修改之后title相同的数据
	stored, err := s.repo.FindCategoryByTitle(category.Title)
	if err != nil && !errors.Is(err, codes.ErrImgCategoryNotFound) {
		return nil, err
	}
	if stored != nil && stored.ID != category.ID {
		return nil, codes.ErrImgCategoryTitleRepeat
	}

	// 2.再去查询原先数据 比对path是否一致
	// 若一致则允许更新
	// 若不一致 则需该分类下无图片关联方可进行修改
	old, err := s.repo.FindCategoryByID(category.ID)
	if err != nil {
		return nil, err
	}
	if old.Prefix == category.Prefix {
		return s.repo.UpdateCategory(category)
	}
	// 若不一致
	if err := s.isCategoryExistImg(category.ID); err != nil {
		return nil, err
	}
	return s.repo.UpdateCategory(category)
}

func (s *service) DeleteCategory(id int64) error {
	// 检验当前分类下是否存在图片
	if err := s.isCategoryExistImg(id); err != nil {
		return err
	}

	return s.repo.DeleteCategory(id)
}

func (s *service) ListCategories() (categories []*domain.Category, err error) {
	return s.repo.ListCategories()
}
