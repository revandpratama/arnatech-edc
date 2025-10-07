package middleware

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/edc-service/util"
)

type ValidateError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func Validate(dto any) fiber.Handler {
	return func(c *fiber.Ctx) error {

		body := reflect.New(reflect.TypeOf(dto)).Interface()
		if err := c.BodyParser(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
		}

		if err := util.Validate.Struct(body); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				var errors []ValidateError
				for _, e := range validationErrors {
					errors = append(errors, ValidateError{
						Field: e.Field(),
						Tag:   e.Tag(),
						Value: e.Param(),
					})
				}
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":  "validation failed",
					"fields": errors,
				})
			}
		}

		c.Locals("body", body)

		return c.Next()
	}
}
