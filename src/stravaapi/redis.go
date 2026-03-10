package stravaapi

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const SegmentCacheTTL = time.Hour * 24 * 7
const SegmentCacheNamespace = "stravaapi:segment:"

func PutSegmentInCache(client *redis.Client, segment Segment) error {
	segmentJson, err := json.Marshal(segment)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := client.Set(ctx, SegmentCacheNamespace+strconv.Itoa(segment.Id),
		segmentJson, SegmentCacheTTL).Err(); err != nil {

		return err
	}

	return nil
}

func GetSegmentFromCache(client *redis.Client, id int) (segment Segment,
	err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	segmentJson, err := client.Get(ctx, SegmentCacheNamespace+strconv.Itoa(
		id)).Bytes()
	if err != nil {
		return Segment{}, err
	}

	if err := json.Unmarshal(segmentJson, &segment); err != nil {
		return Segment{}, err
	}

	return segment, nil
}
