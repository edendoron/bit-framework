package models

import (
	"github.com/go-playground/validator"
	"log"
)

func ValidateType(response interface{}) error {
	v := validator.New()
	err := v.Struct(response)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e.Error())
		}
		return err
	}
	return nil
}
