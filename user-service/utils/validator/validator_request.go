package validator

import (
	"errors"
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	translator, found := uni.GetTranslator("en")
	if !found {
		log.Fatalf("[NewValidator] Translator not found: %v", found)
	}

	validator := validator.New()

	return &Validator{
		Validator:  validator,
		Translator: translator,
	}
}

func (v *Validator) Validate(s interface{}) error {
	err := v.Validator.Struct(s)
	if err != nil {
		object, _ := err.(validator.ValidationErrors)
		for _, e := range object {
			log.Printf("[Validator-1] Validation error: Field %s, Tag %s, Param %s, Translate: %s", e.Field(), e.Tag(), e.Param(), e.Translate(v.Translator))

			return errors.New(e.Translate(v.Translator))
		}
	}
	return nil
}
