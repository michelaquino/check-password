package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.globoi.com/michel.aquino/check-password/repository"
)

func GetListCredentials(echoContext echo.Context) error {
	credetialsList := repository.ListCredentials()

	return echoContext.Render(http.StatusOK, "listCredentials", credetialsList)
}
