package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrRecipeNotFound = errors.New("recipe not found")
)

type RecipeRepository interface {
	Create(context.Context, Recipe) (Recipe, error)
	Update(context.Context, RecipeID, Recipe) (Recipe, error)
	FindAll(context.Context) ([]Recipe, error)
	FindByID(context.Context, RecipeID) (Recipe, error)
	Delete(context.Context, RecipeID) error
}

type Recipe struct {
	ID           RecipeID  `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

type RecipeID string

func (r RecipeID) String() string {
	return string(r)
}

func NewRecipe(ID RecipeID, name string, tags []string, ingredients []string, instructions []string, publishedAt time.Time) Recipe {
	return Recipe{
		ID:           ID,
		Name:         name,
		Tags:         tags,
		Ingredients:  ingredients,
		Instructions: instructions,
		PublishedAt:  publishedAt,
	}

}
