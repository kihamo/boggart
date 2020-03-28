package vacuum

import (
	"context"
	"time"

	"github.com/kihamo/boggart/providers/xiaomi/miio"
)

type CleanSummary struct {
	TotalTime     time.Duration
	TotalArea     uint64 // mm2
	TotalCleanups uint64
	CleanupIDs    []uint64
}

type CleanDetail struct {
	StartTime        time.Time
	EndTime          time.Time
	CleaningDuration time.Duration
	Area             uint64 // mm2
	Completed        bool
}

func (d *Device) CleanSummary(ctx context.Context) (CleanSummary, error) {
	type response struct {
		miio.Response

		Result []interface{} `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_clean_summary", nil, &reply)
	if err != nil {
		return CleanSummary{}, err
	}

	result := CleanSummary{}

	for i, v := range reply.Result {
		switch i {
		case 0:
			result.TotalTime = time.Duration(v.(float64)) * time.Second

		case 1:
			result.TotalArea = uint64(v.(float64))

		case 2:
			result.TotalCleanups = uint64(v.(float64))

		case 3:
			values := v.([]interface{})
			result.CleanupIDs = make([]uint64, len(values))

			for i2, v2 := range values {
				result.CleanupIDs[i2] = uint64(v2.(float64))
			}
		}
	}

	return result, nil
}

func (d *Device) CleanDetails(ctx context.Context, id uint64) (CleanDetail, error) {
	type response struct {
		miio.Response

		Result [][]int64 `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_clean_record", []uint64{id}, &reply)
	if err != nil {
		return CleanDetail{}, err
	}

	result := CleanDetail{}

	for i, v := range reply.Result[0] {
		switch i {
		case 0:
			result.StartTime = time.Unix(v, 0)

		case 1:
			result.EndTime = time.Unix(v, 0)

		case 2:
			result.CleaningDuration = time.Duration(v) * time.Second

		case 3:
			result.Area = uint64(v)

		case 5:
			result.Completed = v == 1
		}
	}

	return result, nil
}
