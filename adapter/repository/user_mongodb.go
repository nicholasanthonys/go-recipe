package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	"github.com/nicholasanthonys/go-recipe/domain"
)

type userBSON struct {
	ID        string    `bson:"id"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"createdAt"`
}

type UserNoSQL struct {
	collectionName string
	db             NoSQL
}

func NewUserNoSQL(db NoSQL) UserNoSQL {
	return UserNoSQL{
		db:             db,
		collectionName: "users",
	}
}

func (r UserNoSQL) Create(ctx context.Context, user domain.User) (domain.User, error) {
	var userBSON = userBSON{
		ID:       user.ID.String(),
		Email:    user.Email,
		Password: user.Password,
	}

	if err := r.db.Store(ctx, r.collectionName, userBSON); err != nil {
		return domain.User{}, errors.Wrap(err, "error creating user")

	}

	return user, nil

}

func (a UserNoSQL) FindByEmailAndPass(ctx context.Context, email string, password string) (domain.User, error) {
	var (
		userBSON = &userBSON{}
		query    = bson.M{
			"email":  email,
			password: password,
		}
	)

	if err := a.db.FindOne(ctx, a.collectionName, query, nil, userBSON); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return domain.User{}, domain.ErrAccountNotFound
		default:
			return domain.User{}, errors.Wrap(err, "error fetching account")
		}
	}

	return domain.NewUser(
		domain.UserID(userBSON.ID),
		userBSON.Email,
		userBSON.Password,
		userBSON.CreatedAt,
	), nil
}

func (r UserNoSQL) Update(ctx context.Context, ID domain.UserID, user domain.User) (domain.User, error) {
	query := bson.M{"id": user.ID}
	update := bson.M{"$set": bson.M{
		"email":    user.Email,
		"password": user.Password,
	}}

	if err := r.db.Update(ctx, r.collectionName, query, update); err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return domain.User{}, errors.Wrap(domain.ErrUserNotFound, "error update user")
		default:
			return domain.User{}, errors.Wrap(err, "error update user")
		}
	}
	return user, nil

}

func (r UserNoSQL) FindAll(ctx context.Context) ([]domain.User, error) {
	var usersBSON = make([]userBSON, 0)
	if err := r.db.FindAll(ctx, r.collectionName, bson.M{}, &usersBSON); err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return []domain.User{}, errors.Wrap(domain.ErrUserNotFound, "error listing users")
		default:
			return []domain.User{}, errors.Wrap(err, "error listing users")
		}
	}

	var users = make([]domain.User, 0)
	for _, userBson := range usersBSON {
		var user = domain.NewUser(
			domain.UserID(userBson.ID),
			userBson.Email,
			"",
			userBson.CreatedAt,
		)

		users = append(users, user)

	}

	return users, nil

}

func (r UserNoSQL) FindByID(ctx context.Context, ID domain.UserID) (domain.User, error) {

	var (
		userBSON = &userBSON{}
		query    = bson.M{"id": ID}
	)

	if err := r.db.FindOne(ctx, r.collectionName, query, nil, userBSON); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return domain.User{}, domain.ErrUserNotFound
		default:
			return domain.User{}, errors.Wrap(err, "error fetching user")
		}

	}
	return domain.NewUser(
		domain.UserID(userBSON.ID),
		userBSON.Email,
		"",
		userBSON.CreatedAt,
	), nil

}

func (r UserNoSQL) Delete(ctx context.Context, ID domain.UserID) error {
	query := bson.M{"id": ID}
	if err := r.db.Delete(ctx, r.collectionName, query); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return domain.ErrUserNotFound
		default:
			return errors.Wrap(err, "error deleting user")
		}
	}
	return nil

}
