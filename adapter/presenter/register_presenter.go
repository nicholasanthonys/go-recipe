package presenter

import (
	"github.com/nicholasanthonys/go-recipe/domain"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

type registerPresenter struct{}

func NewRegisterPresenter() usecase.RegisterPresenter {
	return registerPresenter{}
}

func (a registerPresenter) Output(user domain.User) usecase.RegisterOutput {
	o := usecase.RegisterOutput{
		ID:    user.ID.String(),
		Email: user.Email,
	}

	return o
}
