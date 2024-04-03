package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

type IValidator interface {
	BindAndValidate(c *fiber.Ctx, data interface{}) error
}

type Validator struct {
	Validator *validator.Validate
}

func New() *Validator {
	return &Validator{
		Validator: validator.New(),
	}
}

func (v *Validator) ValidatorError(errs validator.ValidationErrors) string {
	var sb strings.Builder
	for _, err := range errs {
		sb.WriteString(fmt.Sprintf("field: %s must be %s and %s", err.Field(), err.Type().String(), err.Tag()))
	}
	return sb.String()
}

func (v *Validator) RegisterValidation(data interface{}) error {
	v.Validator.RegisterTagNameFunc(v.getTagNameFunc("json"))
	v.Validator.RegisterTagNameFunc(v.getTagNameFunc("query"))
	v.Validator.RegisterTagNameFunc(v.getTagNameFunc("params"))
	var customErr validator.ValidationErrors
	if err := v.Validator.Struct(data); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			customErr = append(customErr, err)
		}
		return errors.New(v.ValidatorError(customErr))
	}
	return nil
}

func (v *Validator) getTagNameFunc(tag string) func(fld reflect.StructField) string {
	return func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	}
}

func (v *Validator) parseAndValidate(c *fiber.Ctx, data interface{}) error {
	var (
		err  error
		err2 error
		err3 error
	)
	err = c.BodyParser(data)
	err2 = c.QueryParser(data)
	err3 = c.ParamsParser(data)
	if err != nil && err2 != nil && err3 != nil {
		return err
	}
	if err := v.RegisterValidation(data); err != nil {
		return err
	}
	return nil
}

func (v *Validator) BindAndValidate(c *fiber.Ctx, data interface{}) error {
	if err := v.parseAndValidate(c, data); err != nil {
		return err
	}
	return nil
}
