package handler

import (
	"fmt"
	"net/http"

	"github.com/michelaquino/check-password/context"
	"github.com/michelaquino/check-password/models"
	"github.com/michelaquino/check-password/repository"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
)

func GetCredentials(echoContext echo.Context) error {
	return echoContext.Render(http.StatusOK, "getCredentials", nil)
}

func PostCredentials(echoContext echo.Context) error {
	log := context.GetLogger()

	credentials := new(models.Credentials)
	if err := echoContext.Bind(credentials); err != nil {
		log.Error("Bind form to object", "Error", err.Error())

		viewModel := GetCredentialViewModel{
			HasError:     true,
			ErrorMessage: "Ocorreu um erro inesperado",
		}

		return echoContext.Render(http.StatusInternalServerError, "getCredentials", viewModel)
	}

	if credentials.Email == "" || !govalidator.IsEmail(credentials.Email) {
		viewModel := GetCredentialViewModel{
			HasError:     true,
			ErrorMessage: "O e-mail é inválido",
		}

		log.Info("Validate email", "Error", fmt.Sprintf("Invalid email: %s", credentials.Email))
		return echoContext.Render(http.StatusInternalServerError, "getCredentials", viewModel)
	}

	if credentials.Password == "" {
		viewModel := GetCredentialViewModel{
			HasError:     true,
			ErrorMessage: "A senha não pode ser vazia",
		}

		log.Info("Validate password length", "Error", "Password length invalid")
		return echoContext.Render(http.StatusInternalServerError, "getCredentials", viewModel)
	}

	credentials.SetPasswordHash()
	if err := repository.SaveCredentials(credentials); err != nil {
		viewModel := GetCredentialViewModel{
			HasError:     true,
			ErrorMessage: "Ocorreu um erro ao salvar as credenciais",
		}

		log.Error("Save credentials on database", "Error", err.Error())
		return echoContext.Render(http.StatusInternalServerError, "getCredentials", viewModel)
	}

	log.Info("Save credentials", "Success", "Credentials save with success")
	return echoContext.Render(http.StatusOK, "getCredentials", nil)
}

type GetCredentialViewModel struct {
	HasError     bool
	ErrorMessage string
}
