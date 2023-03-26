package action

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholasanthonys/go-recipe/adapter/api/logging"
	"github.com/nicholasanthonys/go-recipe/adapter/api/response"
	"github.com/nicholasanthonys/go-recipe/adapter/logger"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

type SignInAction struct {
	uc  usecase.SignInUseCase
	log logger.Logger
	c   *gin.Context
}

func NewSignInAction(uc usecase.SignInUseCase, log logger.Logger, c *gin.Context) SignInAction {
	return SignInAction{
		uc:  uc,
		log: log,
		c:   c,
	}
}

func (a SignInAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "sign_in"

	var input usecase.SignInInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log("error when decoding json")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	output, err := a.uc.Execute(r.Context(), input, a.c)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when sign in")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	logging.NewInfo(a.log, logKey, http.StatusCreated).Log("success sign in")
	response.NewSuccess(output, http.StatusOK).Send(w)
}
