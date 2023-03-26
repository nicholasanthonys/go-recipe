package usecase

import (
	"context"
	"time"

	"github.com/nicholasanthonys/go-recipe/domain"
)

type (
	// CreateRecipe UseCase input port
	CreateRecipeUseCase interface {
		Execute(context.Context, CreateRecipeInput) (CreateRecipeOutput, error)
	}

	// CreateRecipe Input input data
	CreateRecipeInput struct {
		Name         string   `json:"name" validate:"required"`
		Tags         []string `json:"tags" validate:"required"`
		Ingredients  []string `json:"ingredients" validate:"required"`
		Instructions []string `json:"instructions" validate:"required"`
	}

	// CreateRecipePresenter output port
	CreateRecipePresenter interface {
		Output(domain.Recipe) CreateRecipeOutput
	}

	// CreateRecipeOutput output data
	CreateRecipeOutput struct {
		ID           string   `json:"id"`
		Name         string   `json:"name"`
		Tags         []string `json:"tags"`
		Ingredients  []string `json:"ingredients"`
		Instructions []string `json:"instructions"`
		PublishedAt  string   `json:"published_at"`
	}

	createRecipeInteractor struct {
		repo       domain.RecipeRepository
		presenter  CreateRecipePresenter
		ctxTimeout time.Duration
	}
)

// NewCreateRecipeInteractor creates new createRecipeInteractor with its dependencies
func NewCreateRecipeInteractor(
	repo domain.RecipeRepository,
	presenter CreateRecipePresenter,
	t time.Duration,
) CreateRecipeUseCase {
	return createRecipeInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (a createRecipeInteractor) Execute(ctx context.Context, input CreateRecipeInput) (CreateRecipeOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	var Recipe = domain.NewRecipe(
		domain.RecipeID(domain.NewUUID()),
		input.Name,
		input.Tags,
		input.Ingredients,
		input.Instructions,
		time.Now(),
	)

	Recipe, err := a.repo.Create(ctx, Recipe)
	if err != nil {
		return a.presenter.Output(domain.Recipe{}), err
	}

	return a.presenter.Output(Recipe), nil
}
