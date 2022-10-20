package usecase

import (
	"context"
	"time"

	"github.com/nicholasanthonys/go-recipe/domain"
)

type (
	// FindAllRecipeUseCase input port
	FindAllRecipeUseCase interface {
		Execute(context.Context) ([]FindAllRecipeOutput, error)
	}

	// FindAllRecipePresenter output port
	FindAllRecipePresenter interface {
		Output([]domain.Recipe) []FindAllRecipeOutput
	}

	// FindAllRecipeOutput outputData
	FindAllRecipeOutput struct {
		ID           string   `json:"id"`
		Name         string   `json:"name"`
		Tags         []string `json:"tags"`
		Ingredients  []string `json:"ingredients"`
		Instructions []string `json:"instructions"`
		PublishedAt  string   `json:"published_at"`
	}

	findAllRecipeInteractor struct {
		repo       domain.RecipeRepository
		presenter  FindAllRecipePresenter
		ctxTimeout time.Duration
	}
)

// NewFindAllRecipeInteractor creates new findAllRecipeInteractor with its dependencies
func NewFindAllRecipeInteractor(
	repo domain.RecipeRepository,
	presenter FindAllRecipePresenter,
	t time.Duration,
) FindAllRecipeUseCase {
	return findAllRecipeInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (a findAllRecipeInteractor) Execute(ctx context.Context) ([]FindAllRecipeOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	recipes, err := a.repo.FindAll(ctx)
	if err != nil {
		return a.presenter.Output([]domain.Recipe{}), err
	}

	return a.presenter.Output(recipes), nil
}
