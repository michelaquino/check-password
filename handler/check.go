package handler

import (
	"net/http"

	"gitlab.globoi.com/michel.aquino/check-password/models"
	"gitlab.globoi.com/michel.aquino/check-password/repository"

	"github.com/labstack/echo"
)

func GetCheck(echoContext echo.Context) error {
	return echoContext.Render(http.StatusOK, "login", "")
}

func PostCheck(echoContext echo.Context) error {
	credentials := new(models.Credentials)

	if err := echoContext.Bind(credentials); err != nil {
		return echoContext.Render(http.StatusInternalServerError, "login", "")
	}

	repository.SaveCredentials(credentials)

	return echoContext.Render(http.StatusOK, "login", "")
}
