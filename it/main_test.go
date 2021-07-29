package it

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/handler"
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/VladimirStepanov/todo-app/internal/repository/postgres"
	"github.com/VladimirStepanov/todo-app/internal/repository/redisrepo"
	"github.com/VladimirStepanov/todo-app/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestingSuite struct {
	suite.Suite
	router http.Handler
}

var (
	defaultPassword = "123456789"
	defaultSalt     = "$2a$10$wHVm4AGd.uq.dR7Zk3VjhOJWLEt9WPXEqoCPx5AEzPtH31o7WiY92"
)

var userForCreate = &models.User{
	Email:         "alreadyexists@mail.ru",
	Password:      defaultPassword,
	IsActivated:   false,
	ActivatedLink: "user_for_create",
}

var notConfirmedUser = &models.User{
	Email:         "notconfirm@mail.ru",
	Password:      defaultPassword,
	IsActivated:   false,
	ActivatedLink: "not_confirmed_user",
}

var confirmedUser = &models.User{
	Email:         "confirmed@mail.ru",
	Password:      defaultPassword,
	IsActivated:   true,
	ActivatedLink: "confirmed_user",
}

var authUser = &models.User{
	Email:         "auth@mail.ru",
	Password:      defaultSalt,
	IsActivated:   true,
	ActivatedLink: "auth",
}

var maxLoggedInUser = &models.User{
	Email:         "maxloggedin@mail.ru",
	Password:      defaultSalt,
	IsActivated:   true,
	ActivatedLink: "auth",
}

var dataForInsert = []*models.User{
	userForCreate,
	notConfirmedUser,
	confirmedUser,
	authUser,
	maxLoggedInUser,
}

var (
	accessKey        = "accessKey"
	refreshKey       = "refreshKey"
	maxLoggenInCount = 3
)

func initDb(t *testing.T, db *sqlx.DB) {

	for _, u := range dataForInsert {
		_, err := db.NamedExec("INSERT INTO users(email, password_hash, is_activated, activated_link) values(:email, :password_hash, :is_activated, :activated_link)", u)
		if err != nil {
			t.Fatal("Error while initDb", err)
		}
	}
}

func initRedis(t *testing.T, client *redis.Client) {
	require.NoError(t, client.Do(context.Background(), "FLUSHALL").Err())
	require.NoError(t, client.Do(context.Background(), "FLUSHDB").Err())
}

func (suite *TestingSuite) SetupSuite() {
	db, err := postgres.NewDB(
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), "disable",
	)

	if err != nil {
		suite.T().Fatal("Can't create NewDB", err)
	}

	redisClient, err := redisrepo.NewRedisClient(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	if err != nil {
		suite.T().Fatal("Can't create NewDB", err)
	}

	initDb(suite.T(), db)
	initRedis(suite.T(), redisClient)

	repo := postgres.NewPostgresRepository(db)
	tokenRepo := redisrepo.NewRedisRepository(redisClient)
	userService := service.NewUserService(repo)
	tokenService := service.NewTokenService(accessKey, refreshKey, maxLoggenInCount, tokenRepo)
	msObj := new(mocks.MailService)
	msObj.On("SendConfirmationsEmail", mock.Anything).Return(nil)
	logger := logrus.New()
	logger.Out = ioutil.Discard
	suite.router = handler.New(userService, msObj, tokenService, logger).InitRoutes(gin.TestMode)
}

func TestSuite(t *testing.T) {
	if testing.Short() == false {
		suite.Run(t, new(TestingSuite))
	}
}
