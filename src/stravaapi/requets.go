package stravaapi

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strava-segments-script/credentials"
	"time"
)

const RequestTimeout = time.Second * 10
const APIURL = "https://www.strava.com/api/v3/"

func MakeSampleRequest(accessTokenProvider *credentials.
AccessTokenProvider) error {

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", APIURL+"athlete", nil)
	if err != nil {
		log.Println(err)
		return err
	}

	accessToken, err := accessTokenProvider.GetAccessToken(RequestTimeout)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", string(responseBody))

	return nil
}
