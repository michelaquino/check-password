package handler

import (
	"net/http"

	"gitlab.globoi.com/michel.aquino/check-password/models"
	"gitlab.globoi.com/michel.aquino/check-password/repository"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
)

func GetCredentials(echoContext echo.Context) error {
	return echoContext.Render(http.StatusOK, "getCredentials", nil)
}

func PostCredentials(echoContext echo.Context) error {
	credentials := new(models.Credentials)

	if err := echoContext.Bind(credentials); err != nil {
		return echoContext.Render(http.StatusInternalServerError, "getCredentials", nil)
	}

	if credentials.Email == "" || !govalidator.IsEmail(credentials.Email) {
		viewModel := GetCredentialViewModel{
			HasError:     true,
			ErrorMessage: "O e-mail é inválido",
		}

		return echoContext.Render(http.StatusInternalServerError, "getCredentials", viewModel)
	}

	if len(credentials.Password) < 8 {
		viewModel := GetCredentialViewModel{
			HasError:     true,
			ErrorMessage: "A senha deve ter mais que 8 caracteres",
		}

		return echoContext.Render(http.StatusInternalServerError, "getCredentials", viewModel)
	}

	credentials.SetPasswordHash()
	if err := repository.SaveCredentials(credentials); err != nil {
		viewModel := GetCredentialViewModel{
			HasError:     true,
			ErrorMessage: "Ocorreu um erro ao salvar as credenciais",
		}

		return echoContext.Render(http.StatusInternalServerError, "getCredentials", viewModel)
	}

	return echoContext.Render(http.StatusOK, "getCredentials", nil)
}

type GetCredentialViewModel struct {
	HasError     bool
	ErrorMessage string
}
