package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	ErrWrongOffsetFormat = errors.New("wrong offset format")
)

type Offset struct {
	time.Duration
}

func (i Offset) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *Offset) UnmarshalJSON(data []byte) (err error) {
	interval := string(data)

	trimmed := strings.Trim(interval, `"-`)
	parts := strings.Split(trimmed, ":")
	if len(parts) != 3 {
		return ErrWrongOffsetFormat
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}

	duration := time.Duration(hours) * time.Hour
	duration += time.Duration(minutes) * time.Minute
	duration += time.Duration(seconds) * time.Second

	if []rune(interval)[0] == '-' {
		duration *= -1
	}

	i.Duration = duration

	return nil
}

func (i *Offset) String() (interval string) {
	interval = fmt.Sprintf(
		"%02d:%02d:%02d",
		int(math.Abs(i.Hours())),
		int(math.Abs(i.Minutes()))%60,
		int(math.Abs(i.Seconds()))%60)
	if i.Seconds() < 0 {
		interval = fmt.Sprintf("-%s", interval)
	}
	return
}
