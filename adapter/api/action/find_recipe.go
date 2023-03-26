package action

import (
	"net/http"
	"strings"

	"github.com/nicholasanthonys/go-recipe/adapter/api/logging"
	"github.com/nicholasanthonys/go-recipe/adapter/api/response"
	"github.com/nicholasanthonys/go-recipe/adapter/logger"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

type FindRecipeAction struct {
	uc  usecase.FindRecipeUseCase
	log logger.Logger
}

func NewFindRecipeAction(uc usecase.FindRecipeUseCase, log logger.Logger) FindRecipeAction {
	return FindRecipeAction{
		uc:  uc,
		log: log,
	}
}

func (a FindRecipeAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_recipe"

	splitted := strings.Split(r.URL.Path, "/")
	id := splitted[len(splitted)-1]
	input := usecase.FindRecipeInput{
		ID: id,
	}
	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when returning recipe")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	logging.NewInfo(a.log, logKey, http.StatusOK).Log("success when returning recipe list")

	response.NewSuccess(output, http.StatusOK).Send(w)

}
