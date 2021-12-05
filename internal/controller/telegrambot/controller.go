package telegrambot

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/duyquang6/quote-today/pkg/logging"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Controller struct {
	botApi      *tgbotapi.BotAPI
	chatGroupID int64
}

// NewController creates a new controller.
func NewController(botApi *tgbotapi.BotAPI) *Controller {
	chatID, _ := strconv.ParseInt(os.Getenv("TELEGRAM_BOT_CHAT_ID"), 10, 64)
	return &Controller{botApi: botApi, chatGroupID: chatID}
}

func (s *Controller) HandleTelegramWebhook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			function = "HandleTelegramWebhook"
			ctx      = c.Request.Context()
			logger   = logging.FromContext(ctx).Named(function)
		)
		data, _ := ioutil.ReadAll(c.Request.Body)

		msg := tgbotapi.NewMessage(s.chatGroupID, "App notification release change")
		if _, err := s.botApi.Send(msg); err != nil {
			logger.Error("cannot send tele message, error:", err)
		}
		msg.Text = string(data)
		if _, err := s.botApi.Send(msg); err != nil {
			logger.Error("cannot send tele message, error:", err)
		}
		c.Status(http.StatusOK)
	}
}
