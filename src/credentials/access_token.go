package credentials

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

func GetAccessToken(credentials RefreshToken) (AccessToken, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", strconv.Itoa(credentials.ClientID))
	data.Set("client_secret", credentials.ClientSecret)
	data.Set("refresh_token", credentials.RefreshToken)

	resp, err := http.PostForm("https://www.strava.com/api/v3/oauth/token", data)
	if err != nil {
		return AccessToken{}, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return AccessToken{}, fmt.Errorf("API returned status %s", resp.Status)
	}

	var accessToken AccessToken
	if err := json.NewDecoder(resp.Body).Decode(&accessToken); err != nil {
		return AccessToken{}, err
	}

	return accessToken, nil
}
