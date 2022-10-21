package presenter

import (
	"github.com/nicholasanthonys/go-recipe/domain"
	"github.com/nicholasanthonys/go-recipe/usecase"
)

type findSignInPresenter struct{}

func NewSignInPresenter() usecase.SignInPresenter {
	return findSignInPresenter{}
}

func (a findSignInPresenter) Output(user domain.User) usecase.SignInOutput {
	o := usecase.SignInOutput{
		ID:    user.ID.String(),
		Email: user.Email,
	}

	return o
}
