package credentials

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type AccessTokenProvider struct {
	RefreshTokenJsonPath string
	RefreshToken         *RefreshToken
	AccessToken          *AccessToken
	APICallCounter       *APICallCounter
}

type APILimitExceededError struct {
	DurationUntilReset time.Duration
}

func (e *APILimitExceededError) Error() string {
	return fmt.Sprintf("API limit exceeded, duration until reset: %s",
		e.DurationUntilReset)
}

func NewAccessTokenProvider(refreshTokenJsonPath string,
	redisClient *redis.Client) (
	*AccessTokenProvider, error) {

	refreshToken, err := MakeRefreshToken(refreshTokenJsonPath)
	if err != nil {
		return nil, err
	}

	accessToken, err := GetAccessTokenFromCache()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if accessToken.ExpiresAt.Before(time.Now()) || errors.Is(err,
		os.ErrNotExist) {
		accessToken, err = getAccessTokenFromRefresh(refreshToken)
	}

	return &AccessTokenProvider{
			RefreshTokenJsonPath: refreshTokenJsonPath,
			RefreshToken:         &refreshToken,
			AccessToken:          &accessToken,
			APICallCounter:       NewAPICallCounter(redisClient)},
		nil
}

func (provider *AccessTokenProvider) GetAccessToken(timeout time.Duration) (
	accessToken string, err error) {

	if provider.AccessToken.ExpiresAt.Before(time.Now().Add(timeout)) {
		*provider.AccessToken, err = getAccessTokenFromRefresh(
			*provider.RefreshToken)
		if err != nil {
			return "", err
		}
	}

	aPILimitIsFine, aPILimitTTL, err := provider.APICallCounter.IsFine()
	if err != nil {
		return "", err
	}

	if !aPILimitIsFine {
		return "", &APILimitExceededError{DurationUntilReset: aPILimitTTL}
	}

	return provider.AccessToken.AccessToken, nil
}

func getAccessTokenFromRefresh(refreshToken RefreshToken) (
	accessToken AccessToken, err error) {

	accessToken, err = GetAccessTokenFromRefresh(refreshToken)
	if err != nil {
		return AccessToken{}, err
	}
	if err := SaveAccessTokenToCache(accessToken); err != nil {
		log.Fatal(err)
	}
	// TODO updated refresh token if changed

	return accessToken, nil
}
