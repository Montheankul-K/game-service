package custom

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"sync"
)

type (
	EchoRequest interface {
		Bind(obj interface{}) error
	}

	customEchoRequest struct {
		ctx       echo.Context
		validator *validator.Validate
	}
)

var (
	once              sync.Once
	validatorInstance *validator.Validate
)

func NewCustomEchoRequest(ctx echo.Context) EchoRequest {
	once.Do(func() {
		validatorInstance = validator.New()
	})
	return &customEchoRequest{
		ctx:       ctx,
		validator: validatorInstance,
	}
}

func (r *customEchoRequest) Bind(obj interface{}) error {
	if err := r.ctx.Bind(obj); err != nil {
		return err
	}

	if err := r.validator.Struct(obj); err != nil {
		return err
	}
	return nil
}
