package battlegrounds

import (
	"EventSender/config"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"net/smtp"
)

type Response struct {
	status int
	err    error
}

func SendInvite(logger *slog.Logger, config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawByte, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to ")
			response(w, r, http.StatusInternalServerError, err)
			return
		}
		auth := smtp.PlainAuth("", config.MailSender.Mail, config.MailSender.Password, config.MailSender.Host)

		to := []string{"prettypepe@mail.ru"}

		msg := []byte("To: prettypepe@mail.ru\r\n" +
			"From:" + config.MailSender.Mail + "\r\n" +
			"Subject: Want to play battlegrounds\r\n" +
			"\r\n" +
			"I want to play bg, lets go to play at 17:30 PM \r\n")

		if err = smtp.SendMail(config.MailSender.Host+":"+config.MailSender.Port, auth, config.MailSender.Mail, to, msg); err != nil {
			logger.Error("failed to send email", slog.String("error", err.Error()))
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		logger.Info("successfully send email")
		logger.Info("decoded body", slog.String("body", string(rawByte)))
		response(w, r, http.StatusOK, nil)
	}
}

func response(w http.ResponseWriter, r *http.Request, status int, err error) {
	render.JSON(w, r, Response{
		status: status,
		err:    err,
	})
}
