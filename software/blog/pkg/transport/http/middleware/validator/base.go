package validator

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"strings"
)

const (
	ValidatorKey  = "ValidatorKey"
	TranslatorKey = "TranslatorKey"
)

var ValidMap = map[string]map[string]string{
	"valid_password": {"en": "Must be at least 8 in length and contain at least one letter and one number",
		"zh": "长度至少为8，至少含有一个字母和一个数字"},
	"valid_iplist": {"en": "invalid ip format", "zh": "不符合ip格式"},
}

func GetValidMsg(valildKey, language string) string {
	val, ok := ValidMap[valildKey]
	if !ok {
		return ""
	}
	msg, ok := val[language]
	if !ok {
		return val["en"]
	}
	return msg
}

func GetValidator(c *gin.Context) (*validator.Validate, error) {
	v, ok := c.Get(ValidatorKey)
	if !ok {
		return nil, errors.New("未设置验证器")
	}
	validator, ok := v.(*validator.Validate)
	if !ok {
		return nil, errors.New("获取验证器失败")
	}
	return validator, nil
}

func GetTranslation(c *gin.Context) (ut.Translator, error) {
	trans, ok := c.Get(TranslatorKey)
	if !ok {
		return nil, errors.New("未设置翻译器")
	}
	translator, ok := trans.(ut.Translator)
	if !ok {
		return nil, errors.New("获取翻译器失败")
	}
	return translator, nil
}

func DefaultGetValidParams(c *gin.Context, params interface{}) error {
	if err := c.ShouldBindJSON(params); err != nil {
		return err
	}

	valid, err := GetValidator(c)
	if err != nil {
		return err
	}

	trans, err := GetTranslation(c)
	if err != nil {
		return err
	}

	err = valid.Struct(params)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, e.Translate(trans))
		}
		return errors.New(strings.Join(sliceErrs, ","))
	}
	return nil
}
