package credentials

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const AccessTokenCacheFilename = "access_token_cache.json"

type AccessToken struct {
	AccessToken  string
	ExpiresAt    time.Time
	RefreshToken string
}

type AccessTokenJson struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenCacheJson struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int    `json:"expires_at"`
}

func GetAccessTokenFromRefresh(refreshToken RefreshToken) (AccessToken,
	error) {

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", strconv.Itoa(refreshToken.ClientID))
	data.Set("client_secret", refreshToken.ClientSecret)
	data.Set("refresh_token", refreshToken.RefreshToken)

	resp, err := http.PostForm("https://www.strava.com/api/v3/oauth/token",
		data)
	if err != nil {
		return AccessToken{}, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return AccessToken{}, fmt.Errorf("API returned status %s", resp.Status)
	}

	var accessToken AccessTokenJson
	if err := json.NewDecoder(resp.Body).Decode(&accessToken); err != nil {
		return AccessToken{}, err
	}

	token := AccessToken{AccessToken: accessToken.AccessToken,
		ExpiresAt: time.Unix(int64(accessToken.ExpiresAt), 0)}
	if accessToken.RefreshToken != refreshToken.RefreshToken {
		token.RefreshToken = accessToken.RefreshToken
	}
	return token, nil
}

func GetAccessTokenFromCache() (AccessToken, error) {
	file, err := os.ReadFile(AccessTokenCacheFilename)
	if err != nil {
		return AccessToken{}, err
	}

	var accessToken AccessTokenCacheJson
	err = json.Unmarshal(file, &accessToken)
	if err != nil {
		return AccessToken{}, err
	}

	if accessToken.AccessToken == "" {
		return AccessToken{}, fmt.Errorf("invalid %s: no access_token",
			AccessTokenCacheFilename)
	}
	if accessToken.ExpiresAt == 0 {
		return AccessToken{}, fmt.Errorf("invalid %s: no expires_at",
			AccessTokenCacheFilename)
	}

	return AccessToken{AccessToken: accessToken.AccessToken,
			ExpiresAt: time.Unix(int64(accessToken.ExpiresAt), 0)},
		nil
}

func SaveAccessTokenToCache(accessToken AccessToken) error {
	accessTokenJson := AccessTokenCacheJson{
		AccessToken: accessToken.AccessToken,
		ExpiresAt:   int(accessToken.ExpiresAt.Unix())}
	file, err := json.MarshalIndent(&accessTokenJson, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(AccessTokenCacheFilename, file, 0644)
}
