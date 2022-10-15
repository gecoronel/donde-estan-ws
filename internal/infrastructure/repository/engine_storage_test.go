package repository

import (
	"testing"

	"github.com/gcoron/donde-estan-ws/internal/bussiness/model"
	mock "github.com/gcoron/donde-estan-ws/internal/infrastructure/repository/mocks"
	"github.com/golang/mock/gomock"
)

var user = model.User{
	ID:       1,
	Name:     "user",
	LastName: "user",
	NumberID: "11100011",
	Username: "user",
	Password: "user",
	Email:    "user@user.com",
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDB(ctrl)
	mockDB.
		EXPECT().
		Get(user.ID).
		Return(&user, nil)

	engineStorage := NewEngineStorage(mockDB)

	userFound, err := engineStorage.GetUser(user.ID)

	if err != nil {
		t.Error(err.Error())
		t.Log("error should be nil")
		t.Fail()
	}

	if userFound == nil {
		t.Error("user should not be nil")
		t.Fail()
	}

	if user.ID != userFound.ID {
		t.Fail()
	}
}
