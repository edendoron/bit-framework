package models

import (
	"fmt"
	"github.com/go-playground/validator"
)

func ValidateType(response interface{}) error {
	v := validator.New()
	err := v.Struct(response)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println(e.Error())
		}
		return err
	}
	return nil
}
