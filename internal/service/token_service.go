package service

import (
	"fmt"
	"time"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenService struct {
	AccessKey   string
	RefreshKey  string
	MaxLoggedIn int
	repo        models.TokenRepository
}

func GenerateToken(uuid string, uID int64, iat, exp int64, key string) (string, error) {
	claims := jwt.MapClaims{}
	claims["uuid"] = uuid
	claims["user_id"] = uID
	claims["exp"] = exp
	claims["iat"] = iat

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return rt.SignedString([]byte(key))
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

	res.AccessToken, err = GenerateToken(res.UUID, userID, res.AccessIAT, res.AccessET, ts.AccessKey)
	if err != nil {
		return nil, err
	}

	res.RefreshET = time.Now().Add(time.Hour * 24 * 7).Unix()
	res.RefreshIAT = time.Now().Unix()

	res.RefreshToken, err = GenerateToken(res.UUID, userID, res.RefreshIAT, res.RefreshET, ts.RefreshKey)
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

func (ts *TokenService) getClaims(tokenString, key string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrBadToken
		}
		return []byte(key), nil
	})

	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			if e.Errors == jwt.ValidationErrorExpired {
				return nil, models.ErrTokenExpired
			}
		}
		return nil, models.ErrBadToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("can't convert token claims to standart claim")
	}
}

func (ts *TokenService) verify(token, key, prefix string) (jwt.MapClaims, error) {
	claims, err := ts.getClaims(token, key)

	if err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf(
		"%s:%d:%s", prefix, int64(claims["user_id"].(float64)), claims["uuid"].(string),
	)

	val, err := ts.repo.Get(redisKey)

	if err != nil {
		return nil, err
	} else if !val {
		return nil, models.ErrUserUnauthorized
	}

	return claims, nil
}

func (ts *TokenService) Refresh(refreshToken string) (*models.TokenDetails, error) {
	claims, err := ts.verify(refreshToken, ts.RefreshKey, "r")

	if err != nil {
		return nil, err
	}

	accessRedisKey := fmt.Sprintf(
		"a:%d:%s", int64(claims["user_id"].(float64)), claims["uuid"].(string),
	)

	refreshRedisKey := fmt.Sprintf(
		"r:%d:%s", int64(claims["user_id"].(float64)), claims["uuid"].(string),
	)

	err = ts.repo.Delete(refreshRedisKey, accessRedisKey)
	if err != nil {
		return nil, err
	}
	return ts.NewTokenPair(int64(claims["user_id"].(float64)))
}

func (ts *TokenService) Verify(token string) (int64, string, error) {
	claims, err := ts.verify(token, ts.AccessKey, "a")

	if err != nil {
		return 0, "", err
	}

	return int64(claims["user_id"].(float64)), claims["uuid"].(string), nil
}

func (ts *TokenService) Logout(userID int64, userUUID string) error {
	accessRedisKey := fmt.Sprintf("a:%d:%s", userID, userUUID)
	refreshRedisKey := fmt.Sprintf("r:%d:%s", userID, userUUID)

	return ts.repo.Delete(refreshRedisKey, accessRedisKey)
}
