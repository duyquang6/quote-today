package wager

import (
	"github.com/duyquang6/quote-today/internal/service"
)

type Controller struct {
	service service.DateQuoteService
}

// NewController creates a new controller.
func NewController(wagerService service.DateQuoteService) *Controller {
	return &Controller{
		service: wagerService,
	}
}
