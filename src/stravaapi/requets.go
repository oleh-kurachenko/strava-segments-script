package stravaapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strava-segments-script/credentials"
	"strconv"
	"time"
)

const RequestTimeout = time.Second * 10
const APIURL = "https://www.strava.com/api/v3/"

type StarredSegmentJson struct {
	Id           int                        `json:"id"`
	Name         string                     `json:"name"`
	ActivityType string                     `json:"activity_type"`
	Distance     float64                    `json:"distance"`
	City         string                     `json:"city"`
	Country      string                     `json:"country"`
	PrEffort     StarredSegmentPrEffortJson `json:"athlete_pr_effort"`
}

type StarredSegmentPrEffortJson struct {
	ElapsedTime int  `json:"elapsed_time"`
	IsKom       bool `json:"is_kom"`
}

type FullSegmentJson struct {
	Id           int                 `json:"id"`
	EffortCount  int                 `json:"effort_count"`
	AthleteCount int                 `json:"athlete_count"`
	StarCount    int                 `json:"star_count"`
	Xom          FullSegmentXomsJson `json:"xoms"`
}

type FullSegmentXomsJson struct {
	Xom string `json:"overall"`
}

func MakeRequest(accessTokenProvider *credentials.AccessTokenProvider,
	urlPath string, urlParams map[string]string) (
	responseBody []byte, err error) {

	parsedUrl, err := url.Parse(APIURL + urlPath)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	for k, v := range urlParams {
		params.Add(k, v)
	}
	parsedUrl.RawQuery = params.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", parsedUrl.String(), nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	accessToken, err := accessTokenProvider.GetAccessToken(RequestTimeout)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	responseBodyRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status: %s", resp.Status)
	}

	return responseBodyRaw, nil
}

func GetStarredSegmentJsons(accessTokenProvider *credentials.
	AccessTokenProvider) ([]StarredSegmentJson, error) {

	segments := make([]StarredSegmentJson, 0)

	for i := 1; ; i++ {
		data, err := MakeRequest(
			accessTokenProvider, "/segments/starred", map[string]string{
				"page":     strconv.Itoa(i),
				"per_page": "200",
			})
		if err != nil {
			return nil, err
		}

		// TODO pass Reader
		var segmentsPage []StarredSegmentJson
		if err := json.Unmarshal(data, &segmentsPage); err != nil {
			return nil, err
		}

		if len(segmentsPage) == 0 {
			break
		}

		segments = append(segments, segmentsPage...)
	}

	return segments, nil
}

func GetSegment(accessTokenProvider *credentials.AccessTokenProvider,
	segmentId int) (segment FullSegmentJson, err error) {
	data, err := MakeRequest(accessTokenProvider,
		"/segments/"+strconv.Itoa(segmentId), nil)
	if err != nil {
		return segment, err
	}

	if err := json.Unmarshal(data, &segment); err != nil {
		return segment, err
	}

	fmt.Println(segment)

	return segment, nil
}
