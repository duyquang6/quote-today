package datequote

import (
	"github.com/duyquang6/quote-today/internal/service"
)

type Controller struct {
	service service.DateQuoteService
}

// NewController creates a new controller.
func NewController(dateQuoteService service.DateQuoteService) *Controller {
	return &Controller{
		service: dateQuoteService,
	}
}
