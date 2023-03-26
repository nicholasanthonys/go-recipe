package presenter

import (
	"github.com/nicholasanthonys/go-recipe/domain"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

type createRecipePresenter struct{}

func NewCreateRecipePresenter() usecase.CreateRecipePresenter {
	return createRecipePresenter{}
}

func (a createRecipePresenter) Output(recipe domain.Recipe) usecase.CreateRecipeOutput {
	return usecase.CreateRecipeOutput{
		ID:           recipe.ID.String(),
		Name:         recipe.Name,
		Tags:         recipe.Tags,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		PublishedAt:  recipe.PublishedAt.String(),
	}
}
