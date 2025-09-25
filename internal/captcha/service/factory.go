package service

import (
	"scaffold/internal/captcha/domain"
	"github.com/pkg/errors"
)

type CaptchaServiceFactor struct {
	generators	map[domain.VerifyWay]domain.CaptchaService
	cache		domain.CaptchaCache
}

func NewCaptchaServiceFactor(cache domain.CaptchaCache) *CaptchaServiceFactor {
	service := &CaptchaServiceFactor{
		cache:		cache,
		generators:	make(map[domain.VerifyWay]domain.CaptchaService),
	}

	// 注册不同类型的验证码生成器
	service.RegisterGenerator(NewImageClickCaptchaService())

	return service
}

func (s *CaptchaServiceFactor) RegisterGenerator(service domain.CaptchaService) {
	s.generators[service.GetVerifyWay()] = service
}

func (s *CaptchaServiceFactor) Generate(way domain.VerifyWay) (*domain.Captcha, error) {
	generator, exists := s.generators[way]
	if !exists {
		return nil, errors.New("不支持的验证码类型")
	}

	response, cacheData, err := generator.Generate()
	if err != nil {
		return nil, err
	}

	id, err := s.cache.Save(way, cacheData)
	if err != nil {
		return nil, err
	}

	response.ID = id
	return response, nil
}

func (s *CaptchaServiceFactor) GenWithAnswer(way domain.VerifyWay) (*domain.CaptchaAnswer, error) {
	generator, exists := s.generators[way]
	if !exists {
		return nil, errors.New("不支持的验证码类型")
	}

	captcha, cacheData, err := generator.Generate()
	if err != nil {
		return nil, err
	}

	id, err := s.cache.Save(way, cacheData)
	if err != nil {
		return nil, err
	}

	res := &domain.CaptchaAnswer{
		ID:	id,
		Image:	captcha.Image,
		Thumb:	captcha.Thumb,
		Audio:	captcha.Audio,
		Way:	way,
		Value:	cacheData,
	}

	return res, nil
}

func (s *CaptchaServiceFactor) Verify(way domain.VerifyWay, id int64, value string) error {
	return s.cache.Verify(way, id, value)
}
