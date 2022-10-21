package action

import (
	"encoding/json"
	"net/http"

	"github.com/nicholasanthonys/go-recipe/adapter/api/logging"
	"github.com/nicholasanthonys/go-recipe/adapter/api/response"
	"github.com/nicholasanthonys/go-recipe/adapter/logger"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

type RegisterAction struct {
	uc  usecase.RegisterUseCase
	log logger.Logger
}

func NewRegisterAction(uc usecase.RegisterUseCase, log logger.Logger) RegisterAction {
	return RegisterAction{
		uc:  uc,
		log: log,
	}
}

func (a RegisterAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "register"

	var input usecase.RegisterInput

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

	output, err := a.uc.Execute(r.Context(), input)
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

	logging.NewInfo(a.log, logKey, http.StatusCreated).Log("success register")
	response.NewSuccess(output, http.StatusOK).Send(w)
}
