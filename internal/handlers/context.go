package handlers

import (
	"errors"
	"net/http"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/utils"
)

var (
	ErrWrongValueType = errors.New("wrong value type")
)

func getUser(r *http.Request) (user storage.User, err error) {
	user, ok := r.Context().Value(utils.ContextKey("user")).(storage.User)
	if !ok {
		return storage.User{}, ErrWrongValueType
	}
	return user, nil
}
