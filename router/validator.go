package router

import (
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
)

// ValidatorManager 验证器管理
var ValidatorManager *Validator

// Validator 验证器
type Validator struct {
	Validate *validator.Validate
	Trans    ut.Translator
}

// InitializeValidator 初始化验证器
func InitializeValidator() {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	ValidatorManager = &Validator{
		Validate: validate,
		Trans:    trans,
	}
}

// Struct 验证结构体
func (c *Validator) Struct(s any) error {
	errs := c.Validate.Struct(s)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			return errors.New(err.Translate(c.Trans))
		}
	}
	return nil
}

// RegisterTagNameFunc 注册Tag名称方法
func (c *Validator) RegisterTagNameFunc() {
	c.Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		// 	如果有接入翻译, 那么显示翻译结果
		return field.Name
	})
}

// AddTranslation 添加翻译文本
func (c *Validator) AddTranslation(tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag = fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = c.Validate.RegisterTranslation(tag, c.Trans, registerFn, transFn)
}
