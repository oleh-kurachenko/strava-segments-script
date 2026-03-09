package credentials

import (
	"errors"
	"log"
	"os"
	"time"
)

type AccessTokenProvider struct {
	RefreshTokenJsonPath string
	RefreshToken         *RefreshToken
	AccessToken          *AccessToken
}

func NewAccessTokenProvider(refreshTokenJsonPath string) (
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
			AccessToken:          &accessToken},
		nil
}

func (provider *AccessTokenProvider) GetAccessToken(timeout time.Duration) (
	accessToken string, err error) {

	if provider.AccessToken.ExpiresAt.Before(time.Now().Add(timeout)) {
		*provider.AccessToken, err = getAccessTokenFromRefresh(
			*provider.RefreshToken)
		if err != nil {
			return accessToken, err
		}
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
