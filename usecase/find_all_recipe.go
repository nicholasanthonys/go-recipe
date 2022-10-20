package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nicholasanthonys/go-recipe/adapter/repository"
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
		repoKV     repository.KeyValStoreIn
		presenter  FindAllRecipePresenter
		ctxTimeout time.Duration
	}
)

// NewFindAllRecipeInteractor creates new findAllRecipeInteractor with its dependencies
func NewFindAllRecipeInteractor(
	repo domain.RecipeRepository,
	repoKV repository.KeyValStoreIn,
	presenter FindAllRecipePresenter,
	t time.Duration,
) FindAllRecipeUseCase {
	return findAllRecipeInteractor{
		repo:       repo,
		repoKV:     repoKV,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (a findAllRecipeInteractor) Execute(ctx context.Context) ([]FindAllRecipeOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	// try to find from cache first
	val, err := a.repoKV.Get(ctx, "recipes")
	if err == repository.KVNotFound {
		fmt.Println("Cache miss. request to mongodb")
		recipes, err := a.repo.FindAll(ctx)
		if err != nil {
			return a.presenter.Output([]domain.Recipe{}), err
		}
		// set to kv store
		kvdata, _ := json.Marshal(recipes)
		a.repoKV.Set(ctx, "recipes", string(kvdata), 0)
		return a.presenter.Output(recipes), nil
	} else if err != nil {
		return a.presenter.Output([]domain.Recipe{}), err
	} else {
		fmt.Println("Cache hit.")
		recipes := make([]domain.Recipe, 0)
		err := json.Unmarshal([]byte(val), &recipes)
		if err != nil {
			return a.presenter.Output([]domain.Recipe{}), err
		}
		return a.presenter.Output(recipes), nil

	}

}
