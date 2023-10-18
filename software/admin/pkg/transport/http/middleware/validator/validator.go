package validator

import (
	"regexp"
	"strings"

	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

func TransactionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		enTranslator := en.New()
		zhTranslator := zh.New()

		uni := ut.New(enTranslator, zhTranslator, enTranslator)
		v := validator.New()

		locale := c.Request.Header.Get("Accept-Language")

		trans, _ := uni.GetTranslator(locale)

		// 注册翻译器
		switch locale {
		case "zh":
			err := zhTrans.RegisterDefaultTranslations(v, trans)
			if err != nil {
				break
			}
			v.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("zh_comment")
			})
			//break
		default:
			err := enTrans.RegisterDefaultTranslations(v, trans)
			if err != nil {
				break
			}
			v.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})
			//break
		}
		//自定义验证方法
		_ = v.RegisterValidation("valid_password", func(fl validator.FieldLevel) bool {
			matched, _ := regexp.Match(`^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`, []byte(fl.Field().String()))
			return matched
		})
		_ = v.RegisterValidation("valid_iplist", func(fl validator.FieldLevel) bool {
			if fl.Field().String() == "" {
				return true
			}
			for _, item := range strings.Split(fl.Field().String(), ",") {
				matched, _ := regexp.Match(`\S+`, []byte(item)) //ip_addr
				if !matched {
					return false
				}
			}
			return true
		})
		//自定义翻译器
		_ = v.RegisterTranslation("valid_password", trans, func(ut ut.Translator) error {
			return ut.Add("valid_password", "{0} "+GetValidMsg("valid_password", locale), true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("valid_password", fe.Field())
			return t
		})
		_ = v.RegisterTranslation("valid_iplist", trans, func(ut ut.Translator) error {
			return ut.Add("valid_iplist", "{0} "+GetValidMsg("valid_iplist", locale), true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("valid_iplist", fe.Field())
			return t
		})
		c.Set(TranslatorKey, trans)
		c.Set(ValidatorKey, v)
		c.Next()
	}

}
