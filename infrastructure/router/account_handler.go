package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gsabadini/go-bank-transfer/adapter/api/action"
	"github.com/gsabadini/go-bank-transfer/adapter/presenter"
	"github.com/gsabadini/go-bank-transfer/adapter/repository"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

func (g ginEngine) buildCreateAccountAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewCreateAccountInteractor(
				repository.NewAccountNoSQL(g.db),
				presenter.NewCreateAccountPresenter(),
				g.ctxTimeout,
			)
			act = action.NewCreateAccountAction(uc, g.log, g.validator)
		)

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildFindAllAccountAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewFindAllAccountInteractor(
				repository.NewAccountNoSQL(g.db),
				presenter.NewFindAllAccountPresenter(),
				g.ctxTimeout,
			)
			act = action.NewFindAllAccountAction(uc, g.log)
		)

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildFindBalanceAccountAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewFindBalanceAccountInteractor(
				repository.NewAccountNoSQL(g.db),
				presenter.NewFindAccountBalancePresenter(),
				g.ctxTimeout,
			)
			act = action.NewFindAccountBalanceAction(uc, g.log)
		)

		q := c.Request.URL.Query()
		q.Add("account_id", c.Param("account_id"))
		c.Request.URL.RawQuery = q.Encode()

		act.Execute(c.Writer, c.Request)
	}
}
