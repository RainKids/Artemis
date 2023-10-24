package captcha

import (
	"admin/pkg/tools/uuid"

	"github.com/mojocn/base64Captcha"
)

type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

var store = NewDefaultRedisStore()

//var store = base64Captcha.DefaultMemStore

func DriverDigitFunc() (id, b64s string, err error) {
	e := configJsonBody{}
	e.Id = uuid.NewUUID()
	e.DriverDigit = base64Captcha.NewDriverDigit(CaptchaImgHeight, CaptchaImgWidth, CaptchaKeyLong, 0.7, 80)
	driver := e.DriverDigit
	cap := base64Captcha.NewCaptcha(driver, store)
	return cap.Generate()
}

func Verify(id, code string, clear bool) bool {
	return store.Verify(id, code, clear)
}
