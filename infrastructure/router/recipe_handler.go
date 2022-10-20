package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasanthonys/go-recipe/adapter/api/action"
	"github.com/nicholasanthonys/go-recipe/adapter/presenter"
	"github.com/nicholasanthonys/go-recipe/adapter/repository"
	"github.com/nicholasanthonys/go-recipe/domain"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

// TODO : this is a dummy data of recipes
var recipes []domain.Recipe

func init() {
	recipes = make([]domain.Recipe, 0)
}

func (g ginEngine) FindRecipeHandler(c *gin.Context) {
	uc := usecase.NewFindRecipeInteractor(
		repository.NewRecipeNoSQL(g.db),
		presenter.NewFindRecipePresenter(),
		g.ctxTimeout,
	)

	ac := action.NewFindRecipeAction(uc, g.log)
	ac.Execute(c.Writer, c.Request)

}

func (g ginEngine) NewRecipeHandler(c *gin.Context) {
	var uc = usecase.NewCreateRecipeInteractor(
		repository.NewRecipeNoSQL(g.db),
		presenter.NewCreateRecipePResenter(),
		g.ctxTimeout,
	)
	act := action.NewCreateRecipeAction(uc, g.log, g.validator)
	act.Execute(c.Writer, c.Request)
}

func (g ginEngine) ListRecipesHandler(c *gin.Context) {
	var (
		uc = usecase.NewFindAllRecipeInteractor(
			repository.NewRecipeNoSQL(g.db),
			presenter.NewFindAllRecipePresenter(),
			g.ctxTimeout,
		)
		act = action.NewFindAllRecipeAction(uc, g.log)
	)
	act.Execute(c.Writer, c.Request)
}

func (g ginEngine) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe domain.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID.String() == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe Not Found",
		})
		return
	}

	recipe.ID = domain.RecipeID(id)
	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
}

func (g ginEngine) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe domain.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID.String() == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe Not Found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)

	c.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
}

func (g ginEngine) SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]domain.Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				listOfRecipes = append(listOfRecipes, recipes[i])
				c.JSON(http.StatusOK, listOfRecipes)
				return
			}

		}

	}

	c.JSON(http.StatusOK, listOfRecipes)

}
