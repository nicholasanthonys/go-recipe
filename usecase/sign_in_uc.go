package usecase

import (
	"context"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nicholasanthonys/go-recipe/domain"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type (

	// SignInUseCase input port
	SignInUseCase interface {
		Execute(context.Context, SignInInput, *gin.Context) (SignInOutput, error)
	}

	// SignInInput input data
	SignInInput struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	// SignInPresenter output port
	SignInPresenter interface {
		Output(domain.User) SignInOutput
	}

	// SignInOutput output data
	SignInOutput struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	signInInteractor struct {
		repo       domain.UserRepository
		presenter  SignInPresenter
		ctxTimeout time.Duration
	}
)

func NewSignInInteractor(
	repo domain.UserRepository,
	presenter SignInPresenter,
	t time.Duration,

) SignInUseCase {
	return signInInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}

}

func (r signInInteractor) Execute(ctx context.Context, input SignInInput, c *gin.Context) (SignInOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	user, err := r.repo.FindByEmail(ctx, domain.Email(input.Email))
	if err != nil {
		return r.presenter.Output(domain.User{}), err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return r.presenter.Output(domain.User{}), errors.Wrap(err, "login failed")
	}

	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("username", user.Email)
	session.Set("token", sessionToken)
	session.Save()

	return r.presenter.Output(user), nil
}
