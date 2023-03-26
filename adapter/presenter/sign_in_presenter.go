package presenter

import (
	"github.com/nicholasanthonys/go-recipe/domain"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

type signInPresenter struct{}

func NewSignInPresenter() usecase.SignInPresenter {
	return signInPresenter{}
}

func (a signInPresenter) Output(user domain.User) usecase.SignInOutput {
	o := usecase.SignInOutput{
		ID:    user.ID.String(),
		Email: user.Email,
	}

	return o
}
