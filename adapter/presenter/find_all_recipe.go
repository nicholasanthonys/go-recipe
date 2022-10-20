package presenter

import (
	"github.com/nicholasanthonys/go-recipe/domain"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

type findAllRecipePresenter struct{}

func NewFindAllRecipePresenter() usecase.FindAllRecipePresenter {
	return findAllRecipePresenter{}
}

func (a findAllRecipePresenter) Output(recipes []domain.Recipe) []usecase.FindAllRecipeOutput {
	var o = make([]usecase.FindAllRecipeOutput, 0)

	for _, recipe := range recipes {
		o = append(o, usecase.FindAllRecipeOutput{
			ID:           recipe.ID.String(),
			Name:         recipe.Name,
			Tags:         recipe.Tags,
			Ingredients:  recipe.Ingredients,
			Instructions: recipe.Instructions,
			PublishedAt:  recipe.PublishedAt.String(),
		})
	}

	return o
}
