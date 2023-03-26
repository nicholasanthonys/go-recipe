package usecase

import (
	"context"
	"time"

	"github.com/nicholasanthonys/go-recipe/domain"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type (

	// RegisterUseCase input port
	RegisterUseCase interface {
		Execute(context.Context, RegisterInput) (RegisterOutput, error)
	}

	// RegisterInput input data
	RegisterInput struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	// RegisterPresenter output port
	RegisterPresenter interface {
		Output(domain.User) RegisterOutput
	}

	// RegisterOutput output data
	RegisterOutput struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	RegisterInteractor struct {
		repo       domain.UserRepository
		presenter  RegisterPresenter
		ctxTimeout time.Duration
	}
)

func NewRegisterInteractor(
	repo domain.UserRepository,
	presenter RegisterPresenter,
	t time.Duration,

) RegisterUseCase {
	return RegisterInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}

}

func (r RegisterInteractor) Execute(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	_, err := r.repo.FindByEmail(ctx, domain.Email(input.Email))
	if err == nil {
		return r.presenter.Output(domain.User{}), errors.New("register failed")
	}

	if err != domain.ErrUserNotFound {
		return r.presenter.Output(domain.User{}), errors.Wrap(err, "register failed")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return r.presenter.Output(domain.User{}), err
	}
	user := domain.User{
		ID:        domain.UserID(domain.NewUUID()),
		Email:     input.Email,
		Password:  string(hashedPass),
		CreatedAt: time.Now(),
	}
	input.Password = string(hashedPass)
	user, err = r.repo.Create(ctx, user)
	if err != nil {
		return r.presenter.Output(domain.User{}), err
	}

	return r.presenter.Output(user), nil
}
