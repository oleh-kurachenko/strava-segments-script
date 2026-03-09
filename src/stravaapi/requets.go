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
	"strings"
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

type Segment struct {
	Id           int
	Name         string
	ActivityType string
	Distance     float64
	City         string
	Country      string
	HasXom       bool
	MyTime       time.Duration
	XomTime      time.Duration
	EffortCount  int
	AthleteCount int
	StarCount    int
}

// Augment Segment with detailed data
func (s *Segment) Augment(ac *credentials.AccessTokenProvider) error {
	detailedSegment, err := GetSegment(ac, s.Id)
	if err != nil {
		return err
	}

	timeStrParts := strings.Split(detailedSegment.Xom.Xom, ":")
	leadingZeros := make([]string, 3-len(timeStrParts))
	for i := range leadingZeros {
		leadingZeros[i] = "0"
	}
	timeStrParts = append(leadingZeros, timeStrParts...)
	s.XomTime, err = time.ParseDuration(
		timeStrParts[0] + "h" + timeStrParts[1] + "m" + timeStrParts[2] + "s")
	if err != nil {
		return err
	}
	s.EffortCount = detailedSegment.EffortCount
	s.AthleteCount = detailedSegment.AthleteCount
	s.StarCount = detailedSegment.StarCount

	return nil
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

func GetStarredSegments(accessTokenProvider *credentials.
AccessTokenProvider) (segments []Segment, err error) {

	segmentJsons := make([]StarredSegmentJson, 0)

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

		segmentJsons = append(segmentJsons, segmentsPage...)
	}

	segments = make([]Segment, len(segmentJsons))
	for i, segmentJson := range segmentJsons {
		segments[i].Id = segmentJson.Id
		segments[i].Name = segmentJson.Name
		segments[i].ActivityType = segmentJson.ActivityType
		segments[i].Distance = segmentJson.Distance
		segments[i].City = segmentJson.City
		segments[i].Country = segmentJson.Country
		segments[i].HasXom = segmentJson.PrEffort.IsKom
		segments[i].MyTime = time.Second * time.Duration(segmentJson.PrEffort.ElapsedTime)

		segments[i].XomTime = -1
		segments[i].EffortCount = -1
		segments[i].AthleteCount = -1
		segments[i].StarCount = -1
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
