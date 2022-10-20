package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	"github.com/nicholasanthonys/go-recipe/domain"
)

type recipeBSON struct {
	ID           string    `bson:"id"`
	Name         string    `bson:"name"`
	Tags         []string  `bson:"tags"`
	Ingredients  []string  `bson:"ingredients"`
	Instructions []string  `bson:"instructions"`
	PublishedAt  time.Time `bson:"publishedAt"`
}

type RecipeNoSQL struct {
	collectionName string
	db             NoSQL
}

func NewRecipeNoSQL(db NoSQL) RecipeNoSQL {
	return RecipeNoSQL{
		db:             db,
		collectionName: "recipes",
	}
}

func (r RecipeNoSQL) Create(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	var recipeBSON = recipeBSON{
		ID:           recipe.ID.String(),
		Name:         recipe.Name,
		Tags:         recipe.Tags,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		PublishedAt:  recipe.PublishedAt,
	}

	if err := r.db.Store(ctx, r.collectionName, recipeBSON); err != nil {
		return domain.Recipe{}, errors.Wrap(err, "error creating recipe")

	}

	return recipe, nil

}

func (r RecipeNoSQL) Update(ctx context.Context, ID domain.RecipeID, recipe domain.Recipe) (domain.Recipe, error) {
	query := bson.M{"id": recipe.ID}
	update := bson.M{"$set": bson.M{
		"name":         recipe.Name,
		"tags":         recipe.Tags,
		"ingredients":  recipe.Ingredients,
		"instructions": recipe.Instructions,
		"publishedAt":  recipe.PublishedAt,
	}}

	if err := r.db.Update(ctx, r.collectionName, query, update); err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return domain.Recipe{}, errors.Wrap(domain.ErrRecipeNotFound, "error update recipe")
		default:
			return domain.Recipe{}, errors.Wrap(err, "error update recipe")
		}
	}
	return recipe, nil

}

func (r RecipeNoSQL) FindAll(ctx context.Context) ([]domain.Recipe, error) {
	var recipesBSON = make([]recipeBSON, 0)
	if err := r.db.FindAll(ctx, r.collectionName, bson.M{}, &recipesBSON); err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return []domain.Recipe{}, errors.Wrap(domain.ErrRecipeNotFound, "error listing recipes")
		default:
			return []domain.Recipe{}, errors.Wrap(err, "error listing recipes")
		}
	}

	var recipes = make([]domain.Recipe, 0)
	for _, recipeBson := range recipesBSON {
		var recipe = domain.NewRecipe(
			domain.RecipeID(recipeBson.ID),
			recipeBson.Name,
			recipeBson.Tags,
			recipeBson.Ingredients,
			recipeBson.Instructions,
			recipeBson.PublishedAt,
		)

		recipes = append(recipes, recipe)

	}

	return recipes, nil

}

func (r RecipeNoSQL) FindByID(ctx context.Context, ID domain.RecipeID) (domain.Recipe, error) {

	var (
		recipeBSON = &recipeBSON{}
		query      = bson.M{"id": ID}
	)

	if err := r.db.FindOne(ctx, r.collectionName, query, nil, recipeBSON); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return domain.Recipe{}, domain.ErrRecipeNotFound
		default:
			return domain.Recipe{}, errors.Wrap(err, "error fetching recipe")
		}

	}
	return domain.NewRecipe(
		domain.RecipeID(recipeBSON.ID),
		recipeBSON.Name,
		recipeBSON.Tags,
		recipeBSON.Ingredients,
		recipeBSON.Instructions,
		recipeBSON.PublishedAt,
	), nil

}

func (r RecipeNoSQL) Delete(ctx context.Context, ID domain.RecipeID) error {
	query := bson.M{"id": ID}
	if err := r.db.Delete(ctx, r.collectionName, query); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return domain.ErrRecipeNotFound
		default:
			return errors.Wrap(err, "error deleting recipe")
		}
	}
	return nil

}
