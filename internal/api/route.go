package api

import (
	"context"
	"net/http"

	dateQuoteControllerPkg "github.com/duyquang6/quote-today/internal/controller/datequote"
	"github.com/duyquang6/quote-today/internal/database"
	"github.com/duyquang6/quote-today/internal/middleware"
	"github.com/duyquang6/quote-today/internal/repository"
	"github.com/duyquang6/quote-today/internal/service"
	"github.com/duyquang6/quote-today/pkg/logging"
	"github.com/gin-gonic/gin"
)

func (s *httpServer) setupDependencyAndRouter(ctx context.Context, r *gin.Engine, db *database.DB) {
	quoteRepo := repository.NewQuoteRepository()
	dateQuoteRepo := repository.NewDateQuoteRepository()

	wagerService := service.NewDateQuoteService(db, quoteRepo, dateQuoteRepo)

	wagerController := dateQuoteControllerPkg.NewController(wagerService)

	initRoute(ctx, r, wagerController)
}

func initRoute(ctx context.Context, r *gin.Engine,
	dateQuoteController *dateQuoteControllerPkg.Controller) {
	r.Use(middleware.PopulateRequestID())
	r.Use(middleware.PopulateLogger(logging.FromContext(ctx)))
	//r.LoadHTMLGlob("web/**/*.html")
	r.LoadHTMLFiles("web/index.tmpl")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"apiEndpoint": "http://localhost:8080"})
	})
	api := r.Group("/api")
	{
		sub := api.Group("/date-quote")
		sub.POST("/like", dateQuoteController.HandleLike())
		sub.POST("/dislike", dateQuoteController.HandleDislike())
		sub.GET("", dateQuoteController.HandleGetRandomDateQuote())
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
