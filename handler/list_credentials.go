package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"gitlab.globoi.com/michel.aquino/check-password/context"
	"gitlab.globoi.com/michel.aquino/check-password/repository"
)

func GetListCredentials(echoContext echo.Context) error {
	log := context.GetLogger()

	onlyHackedCredentials, err := strconv.ParseBool(echoContext.QueryParam("onlyHacked"))
	if err != nil {
		onlyHackedCredentials = false
	}

	credetialsList, err := repository.ListCredentials(onlyHackedCredentials)
	if err != nil {
		log.Error("Get credentials list", "Error", err.Error())
		return echoContext.Render(http.StatusOK, "listCredentials", nil)
	}

	log.Info("Get credentials list", "Success", "")
	return echoContext.Render(http.StatusOK, "listCredentials", credetialsList)
}
