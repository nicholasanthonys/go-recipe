package router

import (
	"net/http"
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
