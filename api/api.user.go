package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/pkg/errors"
	"github.com/renosyah/go-mongo/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	UserModule struct {
		db   *mongo.Database
		Name string
	}
)

func NewUserModule(db *mongo.Database) *UserModule {
	return &UserModule{
		db:   db,
		Name: "module/User",
	}
}

func (m UserModule) All(ctx context.Context, param model.AllUser) ([]model.UserResponse, *Error) {
	var all []model.UserResponse

	data, err := (&model.User{}).All(ctx, m.db, param)
	if err != nil {
		status := http.StatusInternalServerError
		message := "error on query all User"

		if errors.Cause(err) == sql.ErrNoRows {
			status = http.StatusOK
			message = "no User found"
		}
		return []model.UserResponse{}, NewErrorWrap(err, m.Name, "all/User",
			message, status)
	}
	for _, each := range data {
		all = append(all, each.Response())
	}
	return all, nil

}
func (m UserModule) Add(ctx context.Context, param model.User) (model.UserResponse, *Error) {

	i, err := param.Add(ctx, m.db)
	if err != nil {
		status := http.StatusInternalServerError
		message := "error on add User"

		return model.UserResponse{}, NewErrorWrap(err, m.Name, "add/User",
			message, status)
	}

	param.ID = i

	return param.Response(), nil
}

func (m UserModule) One(ctx context.Context, param model.User) (model.UserResponse, *Error) {
	data, err := param.One(ctx, m.db)
	if err != nil {
		status := http.StatusInternalServerError
		message := "error on get one User"

		return model.UserResponse{}, NewErrorWrap(err, m.Name, "one/User",
			message, status)
	}

	return data.Response(), nil
}

func (m UserModule) Update(ctx context.Context, param model.User) (model.UserResponse, *Error) {
	var emptyID primitive.ObjectID

	i, err := param.Update(ctx, m.db)
	if err != nil || i == emptyID {
		status := http.StatusInternalServerError
		message := "error on update User"

		return model.UserResponse{}, NewErrorWrap(err, m.Name, "update/User",
			message, status)
	}
	return param.Response(), nil
}

func (m UserModule) Delete(ctx context.Context, param model.User) (model.UserResponse, *Error) {
	var emptyID primitive.ObjectID
	i, err := param.Delete(ctx, m.db)
	if err != nil || i == emptyID {
		status := http.StatusInternalServerError
		message := "error on delete User"

		return model.UserResponse{}, NewErrorWrap(err, m.Name, "delete/User",
			message, status)
	}
	return param.Response(), nil
}
