package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/rs/xid"
)

// TODO : this is a dummy data of recipes
var recipes []domain.Recipe

func init() {
	recipes = make([]domain.Recipe, 0)
}

func (g ginEngine) NewRecipeHandler(c *gin.Context) {
	var recipe domain.Recipe

	// marshal json into struct
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

func (g ginEngine) ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
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
		if recipes[i].ID == id {
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

	recipe.ID = id
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
		if recipes[i].ID == id {
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
