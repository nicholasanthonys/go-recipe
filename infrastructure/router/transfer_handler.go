package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nicholasanthonys/go-recipe/adapter/api/action"
	"github.com/nicholasanthonys/go-recipe/adapter/presenter"
	"github.com/nicholasanthonys/go-recipe/adapter/repository"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

func (g ginEngine) buildCreateTransferAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewCreateTransferInteractor(
				repository.NewTransferNoSQL(g.db),
				repository.NewAccountNoSQL(g.db),
				presenter.NewCreateTransferPresenter(),
				g.ctxTimeout,
			)

			act = action.NewCreateTransferAction(uc, g.log, g.validator)
		)

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildFindAllTransferAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewFindAllTransferInteractor(
				repository.NewTransferNoSQL(g.db),
				presenter.NewFindAllTransferPresenter(),
				g.ctxTimeout,
			)
			act = action.NewFindAllTransferAction(uc, g.log)
		)

		act.Execute(c.Writer, c.Request)
	}
}
