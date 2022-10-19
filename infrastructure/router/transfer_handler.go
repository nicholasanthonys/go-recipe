package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gsabadini/go-bank-transfer/adapter/api/action"
	"github.com/gsabadini/go-bank-transfer/adapter/presenter"
	"github.com/gsabadini/go-bank-transfer/adapter/repository"
	"github.com/gsabadini/go-bank-transfer/usecase"
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
