package request

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	ContextWrapper interface {
		Bind(data any) error
	}

	contextWrapper struct {
		Context   echo.Context
		Validator *validator.Validate
	}
)

func NewContextWrapper(ctx echo.Context) ContextWrapper {
	return &contextWrapper{
		Context:   ctx,
		Validator: validator.New(),
	}
}

func (c *contextWrapper) Bind(data any) error {
	if err := c.Context.Bind(data); err != nil {
		log.Printf("Error: Bind data failed: %s", err.Error())
	}

	if err := c.Validator.Struct(data); err != nil {
		log.Printf("Error: Validate data failed: %s", err.Error())
	}

	return nil
}
