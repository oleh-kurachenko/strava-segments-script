package stravaapi

import (
	"log"
	"strava-segments-script/credentials"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Segment struct {
	Id           int           `json:"id"`
	Name         string        `json:"name"`
	ActivityType string        `json:"activity_type"`
	Distance     float64       `json:"distance"`
	City         string        `json:"city"`
	Country      string        `json:"country"`
	HasXom       bool          `json:"has_xom"`
	MyTime       time.Duration `json:"my_time"`
	XomTime      time.Duration `json:"xom_time"`
	EffortCount  int           `json:"effort_count"`
	AthleteCount int           `json:"athlete_count"`
	StarCount    int           `json:"star_count"`
}

// Augment Segment with detailed data
func (s *Segment) Augment(rc *redis.Client, ac *credentials.
	AccessTokenProvider) error {
	cachedSegment, err := GetSegmentFromCache(rc, s.Id)
	if err == nil {
		if cachedSegment.Name == s.Name &&
			cachedSegment.Distance == s.Distance &&
			cachedSegment.City == s.City &&
			cachedSegment.Country == s.Country &&
			cachedSegment.MyTime == s.MyTime &&
			cachedSegment.HasXom == s.HasXom {

			*s = cachedSegment
			return nil
		}

		if err := DeleteSegmentFromCache(rc, s.Id); err != nil {
			log.Fatal(err)
		}
	}

	detailedSegment, err := GetSegment(ac, s.Id)
	if err != nil {
		return err
	}

	timeStrParts := strings.Split(detailedSegment.Xom.Xom, ":")
	timeStr := timeStrParts[len(timeStrParts)-1]
	if len(timeStrParts) != 1 {
		leadingZeros := make([]string, 3-len(timeStrParts))
		for i := range leadingZeros {
			leadingZeros[i] = "0"
		}
		timeStrParts = append(leadingZeros, timeStrParts...)
		timeStr = timeStrParts[0] + "h" + timeStrParts[1] + "m" + timeStrParts[2] + "s"
	}
	s.XomTime, err = time.ParseDuration(timeStr)
	if err != nil {
		return err
	}
	s.EffortCount = detailedSegment.EffortCount
	s.AthleteCount = detailedSegment.AthleteCount
	s.StarCount = detailedSegment.StarCount

	if err := PutSegmentInCache(rc, *s); err != nil {
		log.Fatal(err)
	}

	return nil
}

func GetStarredSegments(ac *credentials.AccessTokenProvider) (
	segments []Segment, err error) {

	segmentJsons, err := GetStarredSegmentJsons(ac)
	if err != nil {
		return nil, err
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
