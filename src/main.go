package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strava-segments-script/credentials"
	"strava-segments-script/stravaapi"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type SegmentsByTerrain struct {
	DownhillSegments      []stravaapi.Segment
	DownhillSegmentsWOXor int
	UphillSegments        []stravaapi.Segment
	UphillSegmentsWOXor   int
	FlatSegments          []stravaapi.Segment
	FlatSegmentsWOXor     int
}

func digData(rc *redis.Client, atp *credentials.AccessTokenProvider) (
	segmentsByTerrain SegmentsByTerrain, err error) {

	segments, err := stravaapi.GetStarredSegments(atp)
	if err != nil {
		return SegmentsByTerrain{}, err
	}

	for i := range segments {
		if !segments[i].HasXom {
			if err := segments[i].Augment(rc,
				atp); err != nil {

				return SegmentsByTerrain{}, err
			}
		}

		if strings.HasSuffix(segments[i].Name, "| Downhill") {
			segmentsByTerrain.DownhillSegments = append(segmentsByTerrain.
				DownhillSegments, segments[i])
			if !segments[i].HasXom {
				segmentsByTerrain.DownhillSegmentsWOXor++
			}
		} else if strings.HasSuffix(segments[i].Name, "| Uphill") {
			segmentsByTerrain.UphillSegments = append(segmentsByTerrain.
				UphillSegments, segments[i])
			if !segments[i].HasXom {
				segmentsByTerrain.UphillSegmentsWOXor++
			}
		} else {
			segmentsByTerrain.FlatSegments = append(
				segmentsByTerrain.FlatSegments, segments[i])
			if !segments[i].HasXom {
				segmentsByTerrain.FlatSegmentsWOXor++
			}
		}
	}

	return
}

func presentSegment(segment stravaapi.Segment) {
	fmt.Printf("- \"%-60.60s\" : %v -> %v on %vm distance\n",
		segment.Name,
		segment.MyTime,
		segment.XomTime,
		segment.Distance)
}

func presentData(segs SegmentsByTerrain) {
	fmt.Printf("do not have XOM on %d segments of %d starred\n\n",
		segs.DownhillSegmentsWOXor+
			segs.UphillSegmentsWOXor+
			segs.FlatSegmentsWOXor,
		len(segs.DownhillSegments)+
			len(segs.UphillSegments)+
			len(segs.FlatSegments))

	fmt.Printf("DH: do not have XOM on %d segments of %d starred\n",
		segs.DownhillSegmentsWOXor, len(segs.DownhillSegments))
	sort.Slice(segs.DownhillSegments, func(i, j int) bool {
		return segs.DownhillSegments[i].EffortCount <
			segs.DownhillSegments[j].EffortCount
	})
	for _, seg := range segs.DownhillSegments {
		if !seg.HasXom {
			presentSegment(seg)
		}
	}
	fmt.Println()

	fmt.Printf("UH: do not have XOM on %d segments of %d starred\n",
		segs.UphillSegmentsWOXor, len(segs.UphillSegments))
	sort.Slice(segs.UphillSegments, func(i, j int) bool {
		return segs.UphillSegments[i].EffortCount <
			segs.UphillSegments[j].EffortCount
	})
	for _, seg := range segs.UphillSegments {
		if !seg.HasXom {
			presentSegment(seg)
		}
	}
	fmt.Println()

	fmt.Printf("FL: do not have XOM on %d segments of %d starred\n",
		segs.FlatSegmentsWOXor, len(segs.FlatSegments))
	sort.Slice(segs.FlatSegments, func(i, j int) bool {
		return segs.FlatSegments[i].EffortCount <
			segs.FlatSegments[j].EffortCount
	})
	for _, seg := range segs.FlatSegments {
		if !seg.HasXom {
			presentSegment(seg)
		}
	}
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	accessTokenProvider, err :=
		credentials.NewAccessTokenProvider("local/credentials.json",
			redisClient)
	if err != nil {
		log.Fatal(err)
	}

	data, err := digData(redisClient, accessTokenProvider)
	if err != nil {
		log.Fatal(err)
	}
	presentData(data)
}
