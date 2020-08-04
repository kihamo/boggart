package vacuum

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/kihamo/boggart/providers/xiaomi/miio"
)

var (
	humanTime   = []string{"s", "m", "h"}
	humanArea   = []string{"mm", "cm2", "dm2", "m2"}
	humanPrefix = []string{"K", "M", "G"}
)

type CleanSummary struct {
	TotalTime     time.Duration
	TotalArea     uint64 // mm2
	TotalCleanups uint64
	CleanupIDs    []uint64
}

func (c CleanSummary) TotalTimeHuman() string {
	val := c.TotalTime.Seconds()

	for i := 0; ; i++ {
		v := float64(val / 60)

		if v < 1 || i == len(humanTime)-1 {
			return strconv.Itoa(int(math.Round(val))) + " " + humanTime[i]
		}

		val = v
	}

	return ""
}

func (c CleanSummary) TotalAreaHuman() string {
	val := float64(c.TotalArea)
	step := 100.0

	for i := 0; ; i++ {
		v := val / step

		if v < 1 || i == len(humanArea)-1 {
			prefix := ""
			prefixIndex := 0
			prefixStep := 1000.0

			for {
				val /= 1000

				if val < prefixStep || prefixIndex == len(humanPrefix)-1 {
					prefix = humanPrefix[prefixIndex]
					break
				}

				prefixIndex++
				prefixStep += 1000
			}

			return strconv.Itoa(int(math.Round(val))) + prefix + " " + humanArea[i]
		}

		val = v
	}

	return ""
}

type CleanDetail struct {
	StartTime        time.Time
	EndTime          time.Time
	CleaningDuration time.Duration
	Area             uint64 // mm2
	Completed        bool
}

func (d *Device) CleanSummary(ctx context.Context) (CleanSummary, error) {
	var response struct {
		miio.Response

		Result []interface{} `json:"result"`
	}

	err := d.Client().CallRPC(ctx, "get_clean_summary", nil, &response)
	if err != nil {
		return CleanSummary{}, err
	}

	result := CleanSummary{}

	for i, v := range response.Result {
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
	var response struct {
		miio.Response

		Result [][]int64 `json:"result"`
	}

	err := d.Client().CallRPC(ctx, "get_clean_record", []uint64{id}, &response)
	if err != nil {
		return CleanDetail{}, err
	}

	result := CleanDetail{}

	for i, v := range response.Result[0] {
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
