package models

import (
	"github.com/go-playground/validator"
	"log"
)

// ValidateType make sure that the type given has all required fields.
func ValidateType(response interface{}) error {
	v := validator.New()
	err := v.Struct(response)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e)
		}
		return err
	}
	return nil
}
