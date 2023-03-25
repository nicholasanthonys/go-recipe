package action

import (
	"encoding/json"
	"net/http"

	"github.com/nicholasanthonys/go-recipe/adapter/api/logging"
	"github.com/nicholasanthonys/go-recipe/adapter/api/response"
	"github.com/nicholasanthonys/go-recipe/adapter/logger"
	"github.com/nicholasanthonys/go-recipe/adapter/validator"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

type CreateRecipeAction struct {
	uc        usecase.CreateRecipeUseCase
	log       logger.Logger
	validator validator.Validator
}

func NewCreateRecipeAction(uc usecase.CreateRecipeUseCase, log logger.Logger, v validator.Validator) CreateRecipeAction {
	return CreateRecipeAction{
		uc:        uc,
		log:       log,
		validator: v,
	}
}

func (a CreateRecipeAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_recipe"

	var input usecase.CreateRecipeInput
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

	if errs := a.validateInput(input); len(errs) > 0 {
		logging.NewError(
			a.log,
			response.ErrInvalidInput,
			logKey,
			http.StatusBadRequest,
		).Log("invalid input")

		response.NewErrorMessage(errs, http.StatusBadRequest).Send(w)
		return
	}

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when creating a new recipe")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	logging.NewInfo(a.log, logKey, http.StatusCreated).Log("success creating recipe")

	response.NewSuccess(output, http.StatusCreated).Send(w)
}

func (a CreateRecipeAction) validateInput(input usecase.CreateRecipeInput) []string {
	var msgs []string

	err := a.validator.Validate(input)
	if err != nil {
		for _, msg := range a.validator.Messages() {
			msgs = append(msgs, msg)
		}
	}

	return msgs
}
