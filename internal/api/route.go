package api

import (
	"context"
	"net/http"
	"os"

	dateQuoteControllerPkg "github.com/duyquang6/quote-today/internal/controller/datequote"
	telegramBotControllerPkg "github.com/duyquang6/quote-today/internal/controller/telegrambot"
	"github.com/duyquang6/quote-today/internal/database"
	"github.com/duyquang6/quote-today/internal/middleware"
	"github.com/duyquang6/quote-today/internal/repository"
	"github.com/duyquang6/quote-today/internal/service"
	"github.com/duyquang6/quote-today/pkg/logging"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *httpServer) setupDependencyAndRouter(ctx context.Context, r *gin.Engine, db *database.DB) {
	logger := logging.FromContext(ctx).Named("setupDependencyAndRouter")
	quoteRepo := repository.NewQuoteRepository()
	dateQuoteRepo := repository.NewDateQuoteRepository()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	logger.Info("Authorized on account %s", bot.Self.UserName)
	if err != nil {
		logger.Error("cannot connect telegram bot, error:", err)
	}
	dateQuoteService := service.NewDateQuoteService(db, quoteRepo, dateQuoteRepo)
	dateQuoteController := dateQuoteControllerPkg.NewController(dateQuoteService)
	teleBotController := telegramBotControllerPkg.NewController(bot)
	initRoute(ctx, r, dateQuoteController, teleBotController)
}

func initRoute(ctx context.Context, r *gin.Engine,
	dateQuoteController *dateQuoteControllerPkg.Controller,
	teleBotController *telegramBotControllerPkg.Controller) {
	r.Use(middleware.PopulateRequestID())
	r.Use(middleware.PopulateLogger(logging.FromContext(ctx)))

	appUrl := os.Getenv("APP_URL")
	if len(appUrl) == 0 {
		appUrl = "http://localhost:8080"
	}
	r.LoadHTMLFiles("web/index.tmpl")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"apiEndpoint": appUrl})
	})
	r.Static("/css", "web/css")
	api := r.Group("/api")
	{
		sub := api.Group("/date-quote")
		sub.POST("/like", dateQuoteController.HandleLike())
		sub.POST("/dislike", dateQuoteController.HandleDislike())
		sub.GET("", dateQuoteController.HandleGetRandomDateQuote())
	}

	webhook := r.Group("/webhook")
	{
		webhook.POST("/telebot", teleBotController.HandleTelegramWebhook())
	}
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
