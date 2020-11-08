package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/renosyah/go-mongo/api"
	"github.com/renosyah/go-mongo/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandlerAddUser(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()

	var param model.User

	err := ParseBodyData(ctx, r, &param)
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "User/create/param"),
			http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	return userModule.Add(ctx, param)
}

func HandlerAllUser(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()
	var param model.AllUser

	err := ParseBodyData(ctx, r, &param)
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "User/all/param"),
			http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	return userModule.All(ctx, param)
}

func HandlerOneUser(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "User/detail"),
			http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	return userModule.One(ctx, model.User{ID: id})
}

func HandlerUpdateUser(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	var param model.User

	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "User/update"),
			http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	err = ParseBodyData(ctx, r, &param)
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "User/update/param"),
			http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	param.ID = id

	return userModule.Update(ctx, param)
}

func HandlerDeleteUser(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "User/delete"),
			http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	return userModule.Delete(ctx, model.User{ID: id})
}
