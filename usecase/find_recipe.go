package usecase

import (
	"context"
	"time"

	"github.com/nicholasanthonys/go-recipe/domain"
)

type (

	// FindRecipeUseCase input port
	FindRecipeUseCase interface {
		Execute(context.Context, FindRecipeInput) (FindRecipeOutput, error)
	}

	// FindRecipeInput input data
	FindRecipeInput struct {
		ID string `json:"id" validate:"required"`
	}

	// FindRecipePresenter output port
	FindRecipePresenter interface {
		Output(domain.Recipe) FindRecipeOutput
	}

	// FindRecipeOutput output data
	FindRecipeOutput struct {
		ID           string   `json:"id"`
		Name         string   `json:"name"`
		Tags         []string `json:"tags"`
		Ingredients  []string `json:"ingredients"`
		Instructions []string `json:"instructions"`
		PublishedAt  string   `json:"published_at"`
	}

	findRecipeInteractor struct {
		repo       domain.RecipeRepository
		presenter  FindRecipePresenter
		ctxTimeout time.Duration
	}
)

func NewFindRecipeInteractor(
	repo domain.RecipeRepository,
	presenter FindRecipePresenter,
	t time.Duration,

) FindRecipeUseCase {
	return findRecipeInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}

}

func (r findRecipeInteractor) Execute(ctx context.Context, input FindRecipeInput) (FindRecipeOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	id := domain.RecipeID(input.ID)

	recipe, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return r.presenter.Output(domain.Recipe{}), err
	}
	return r.presenter.Output(recipe), nil

}
