package domain

type CaptchaCache interface {
	Save(way VerifyWay, value string) (int64, error)
	Verify(way VerifyWay, id int64, value string) error
	Delete(key string) error
}
