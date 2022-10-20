package presenter

import (
	"github.com/nicholasanthonys/go-recipe/domain"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

type findRecipePresenter struct{}

func NewFindRecipePresenter() usecase.FindRecipePresenter {
	return findRecipePresenter{}
}

func (a findRecipePresenter) Output(recipe domain.Recipe) usecase.FindRecipeOutput {
	o := usecase.FindRecipeOutput{
		ID:           recipe.ID.String(),
		Name:         recipe.Name,
		Tags:         recipe.Tags,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		PublishedAt:  recipe.PublishedAt.String(),
	}

	return o
}
