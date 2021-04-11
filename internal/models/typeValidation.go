package models

import (
	"fmt"
	"github.com/go-playground/validator"
)

func ValidateType(response interface{}) bool {
	v := validator.New()
	err := v.Struct(response)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println(e.Error())
		}
		return false
	}
	return true
}
