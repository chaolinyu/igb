package liveclock

import (
	"errors"
	"reflect"
	"time"

	"github.com/uxtoolkit/cog"
)

var cogType reflect.Type

type LiveClock struct {
	cog.UXCog
	ticker *time.Ticker
}

func NewLiveClock() *LiveClock {
	liveClock := &LiveClock{}
	liveClock.SetCogType(cogType)
	liveClock.SetCleanupFunc(liveClock.Cleanup)
	return liveClock
}

func (lc *LiveClock) Cleanup() {
	lc.ticker.Stop()
}

func (lc *LiveClock) Start() error {

	const layout = time.RFC1123Z
	var location *time.Location

	if lc.Props["timezoneName"] != nil && lc.Props["timezoneOffset"] != nil {
		location = time.FixedZone(lc.Props["timezoneName"].(string), lc.Props["timezoneOffset"].(int))
	} else {
		return errors.New("The timezoneName and timezoneOffset props need to be set!")
	}

	lc.ticker = time.NewTicker(time.Millisecond * 500)

	go func() {
		for t := range lc.ticker.C {

			if location != nil {
				lc.SetProp("currentTime", t.In(location).Format(layout))
			} else {
				lc.SetProp("currentTime", t.Format(layout))
			}

		}
	}()

	return nil
}

func init() {
	cogType = reflect.TypeOf(LiveClock{})
}
