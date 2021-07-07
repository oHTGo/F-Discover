package helpers

import (
	"context"
	"time"

	ffprobe "github.com/vansante/go-ffprobe"
)

func GetDurationVideo(path string) float64 {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.GetProbeDataContext(ctx, path)
	if err != nil {
		return 0
	}

	return data.Format.DurationSeconds
}
