package it

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/handler"
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/VladimirStepanov/todo-app/internal/repository/postgres"
	"github.com/VladimirStepanov/todo-app/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestingSuite struct {
	suite.Suite
	router http.Handler
}

var userForCreate = &models.User{
	Email:         "alreadyexists@mail.ru",
	Password:      "123456789",
	IsActivated:   false,
	ActivatedLink: "user_for_create",
}

var notConfirmedUser = &models.User{
	Email:         "notconfirm@mail.ru",
	Password:      "123456789",
	IsActivated:   false,
	ActivatedLink: "not_confirmed_user",
}

var confirmedUser = &models.User{
	Email:         "confirmed@mail.ru",
	Password:      "123456789",
	IsActivated:   true,
	ActivatedLink: "confirmed_user",
}

var dataForInsert = []*models.User{
	userForCreate,
	notConfirmedUser,
	confirmedUser,
}

func initDb(t *testing.T, db *sqlx.DB) {

	for _, u := range dataForInsert {
		_, err := db.NamedExec("INSERT INTO users(email, password_hash, is_activated, activated_link) values(:email, :password_hash, :is_activated, :activated_link)", u)
		if err != nil {
			t.Fatal("Error while initDb", err)
		}
	}
}

func (suite *TestingSuite) SetupSuite() {
	db, err := postgres.NewDB(
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), "disable",
	)

	if err != nil {
		suite.T().Fatal("Can't create NewDB", err)
	}

	initDb(suite.T(), db)

	repo := postgres.NewPostgresRepository(db)
	userService := service.NewUserService(repo)
	msObj := new(mocks.MailService)
	msObj.On("SendConfirmationsEmail", mock.Anything).Return(nil)
	logger := logrus.New()
	logger.Out = ioutil.Discard
	suite.router = handler.New(userService, msObj, logger).InitRoutes(gin.TestMode)
}

func TestSuite(t *testing.T) {
	if testing.Short() == false {
		suite.Run(t, new(TestingSuite))
	}
}
