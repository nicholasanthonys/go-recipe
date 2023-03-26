package action

import (
	"net/http"

	"github.com/nicholasanthonys/go-recipe/adapter/api/logging"
	"github.com/nicholasanthonys/go-recipe/adapter/api/response"
	"github.com/nicholasanthonys/go-recipe/adapter/logger"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

type FindAllRecipeAction struct {
	uc  usecase.FindAllRecipeUseCase
	log logger.Logger
}

func NewFindAllRecipeAction(uc usecase.FindAllRecipeUseCase, log logger.Logger) FindAllRecipeAction {
	return FindAllRecipeAction{
		uc:  uc,
		log: log,
	}
}

func (a FindAllRecipeAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_all_recipe"

	output, err := a.uc.Execute(r.Context())
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when returning recipe list")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	logging.NewInfo(a.log, logKey, http.StatusOK).Log("success when returning recipe list")

	response.NewSuccess(output, http.StatusOK).Send(w)
}
