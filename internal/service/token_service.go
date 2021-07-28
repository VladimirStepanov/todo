package service

import (
	"fmt"
	"time"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type TokenService struct {
	AccessKey   string
	RefreshKey  string
	MaxLoggedIn int
	repo        models.TokenRepository
}

func NewTokenService(accessKey, refreshKey string, maxLoggedIn int, repo models.TokenRepository) models.TokenService {
	return &TokenService{
		AccessKey:   accessKey,
		RefreshKey:  refreshKey,
		MaxLoggedIn: maxLoggedIn,
		repo:        repo,
	}
}

func (ts *TokenService) generatePair(userID int64) (*models.TokenDetails, error) {
	var err error

	res := &models.TokenDetails{}

	res.UUID = uuid.NewString()
	res.AccessET = time.Now().Add(time.Minute * 15).Unix()
	res.AccessIAT = time.Now().Unix()

	atClaims := jwt.MapClaims{}
	atClaims["uuid"] = res.UUID
	atClaims["user_id"] = userID
	atClaims["exp"] = res.AccessET
	atClaims["iat"] = res.AccessIAT

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	res.AccessToken, err = at.SignedString([]byte(ts.AccessKey))
	if err != nil {
		return nil, err
	}

	res.RefreshET = time.Now().Add(time.Hour * 24 * 7).Unix()
	res.RefreshIAT = time.Now().Unix()

	rtClaims := jwt.MapClaims{}
	rtClaims["uuid"] = res.UUID
	rtClaims["user_id"] = userID
	rtClaims["exp"] = res.RefreshET
	rtClaims["iat"] = res.RefreshIAT
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	res.RefreshToken, err = rt.SignedString([]byte(ts.RefreshKey))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ts *TokenService) saveTokenPair(userID int64, td *models.TokenDetails) error {
	pattern := fmt.Sprintf("r:%d:*", userID)

	count, err := ts.repo.Count(pattern)

	if err != nil {
		return err
	}

	if count >= ts.MaxLoggedIn {
		return models.ErrMaxLoggedIn
	}

	accessRedisKey := fmt.Sprintf("a:%d:%s", userID, td.UUID)
	refreshRedisKey := fmt.Sprintf("r:%d:%s", userID, td.UUID)

	accessDur := time.Duration((td.AccessET - td.AccessIAT)) * time.Second
	refreshDur := time.Duration(td.RefreshET-td.RefreshIAT) * time.Second

	return ts.repo.SetTokens(accessRedisKey, accessDur, refreshRedisKey, refreshDur)

}

func (ts *TokenService) NewTokenPair(userID int64) (*models.TokenDetails, error) {
	td, err := ts.generatePair(userID)

	if err != nil {
		return nil, err
	}

	err = ts.saveTokenPair(userID, td)
	if err != nil {
		return nil, err
	}
	return td, err
}
