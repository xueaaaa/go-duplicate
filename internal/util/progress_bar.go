package util

import (
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

func NewBar(silent bool, max int, description string) *progressbar.ProgressBar {
	if silent {
		return nil
	}
	return progressbar.NewOptions(max,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetDescription(description),
		progressbar.OptionThrottle(200*time.Millisecond),
		progressbar.OptionShowCount(),
	)
}
