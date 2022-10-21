package usecase

import (
	"context"
	"crypto/sha256"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nicholasanthonys/go-recipe/domain"
	"github.com/rs/xid"
)

type (

	// SignInUseCase input port
	SignInUseCase interface {
		Execute(context.Context, SignInInput, *gin.Context) (SignInOutput, error)
	}

	// SignInInput input data
	SignInInput struct {
		Email    string `json:"string" validate:"required"`
		Password string `json:"string" validate:"required"`
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

	h := sha256.New()

	user, err := r.repo.FindByEmailAndPass(ctx, input.Email, string(h.Sum([]byte(input.Password))))
	if err != nil {
		return r.presenter.Output(domain.User{}), err
	}

	sessionToken := xid.New()
	session := sessions.Default(c)
	session.Set("username", user.Email)
	session.Set("token", sessionToken)

	return r.presenter.Output(user), nil
}
