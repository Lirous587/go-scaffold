package service

import (
	"scaffold/internal/captcha/domain"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/wenlng/go-captcha-assets/bindata/chars"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha/v2/click"
	"log"
	"strconv"
	"strings"
)

type imageClickService struct {
	captcha click.Captcha
}

func NewImageClickCaptchaService() domain.CaptchaService {
	builder := click.NewBuilder()

	// fonts
	fonts, err := fzshengsksjw.GetFont()
	if err != nil {
		log.Fatalln(err)
	}

	// background images
	imgs, err := imagesv2.GetImages()
	if err != nil {
		log.Fatalln(err)
	}

	builder.SetResources(
		click.WithChars(chars.GetChineseChars()),
		click.WithFonts([]*truetype.Font{fonts}),
		click.WithBackgrounds(imgs),
	)

	return &imageClickService{
		captcha: builder.Make(),
	}
}

func (s *imageClickService) GetVerifyWay() domain.VerifyWay {
	return domain.WayImageClick
}

func (s *imageClickService) Generate() (*domain.Captcha, string, error) {
	captData, err := s.captcha.Generate()
	if err != nil {
		return nil, "", err
	}
	dotData := captData.GetData()
	if dotData == nil {
		return nil, "", errors.New("generate err")
	}

	// 数据保存之后生成图片
	var mBase64, tBase64 string
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		return nil, "", errors.WithStack(err)
	}
	tBase64, err = captData.GetThumbImage().ToBase64()
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	var coords []string
	for _, dot := range dotData {
		coords = append(coords, strconv.Itoa(dot.X)+","+strconv.Itoa(dot.Y))
	}

	cacheData := strings.Join(coords, ",")

	res := new(domain.Captcha)
	res.Image = mBase64
	res.Thumb = tBase64

	return res, cacheData, nil
}
