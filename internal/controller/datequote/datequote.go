package wager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/duyquang6/quote-today/pkg/dto"
	"github.com/duyquang6/quote-today/pkg/exception"
	"github.com/duyquang6/quote-today/pkg/validator"
	"github.com/gin-gonic/gin"
)

func (s *Controller) HandleLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err error
		)
		ctx := c.Request.Context()

		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "read data failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		req := dto.LikeRequest{}
		err = json.Unmarshal(data, &req)

		if err := validator.GetValidate().Struct(req); err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "validation error").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		resp, err := s.service.Like(ctx, req)
		if err != nil {
			if appErr, ok := err.(exception.AppError); ok {
				c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func (s *Controller) HandleDislike() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err error
		)
		ctx := c.Request.Context()

		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "read data failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		req := dto.LikeRequest{}
		err = json.Unmarshal(data, &req)

		if err := validator.GetValidate().Struct(req); err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "validation error").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		resp, err := s.service.Dislike(ctx, req)
		if err != nil {
			if appErr, ok := err.(exception.AppError); ok {
				c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func (s *Controller) HandleGetRandomDateQuote() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		resp, err := s.service.GetRandomQuoteInsertIfNotExist(ctx)
		if err != nil {
			if appErr, ok := err.(exception.AppError); ok {
				c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
