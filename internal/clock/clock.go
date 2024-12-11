package clock

import (
	"time"
)

// Clock is an interface for getting the current time in order to ease testing.
type Clock interface {
	Now() time.Time
	NowInZone(zoneName *string) (time.Time, error)
}

// RealClock provides the current real time.
type RealClock struct{}

// Now returns the current time.
func (rc RealClock) Now() time.Time {
	return time.Now()
}

func (rc RealClock) NowInZone(zoneName *string) (time.Time, error) {
	if zoneName == nil {
		return rc.Now(), nil
	}
	loc, err := time.LoadLocation(*zoneName)
	if err != nil {
		return time.Time{}, err
	}
	return rc.Now().In(loc), nil
}

func NewReal() Clock {
	return RealClock{}
}

// FakeClock provides a controllable time for testing.
type FakeClock struct {
	currentTime time.Time
}

// Now returns the current fake time.
func (fc FakeClock) Now() time.Time {
	return fc.currentTime
}

func (fc FakeClock) NowInZone(zoneName *string) (time.Time, error) {
	if zoneName == nil {
		return fc.Now(), nil
	}
	loc, err := time.LoadLocation(*zoneName)
	if err != nil {
		return time.Time{}, err
	}
	return fc.Now().In(loc), nil
}

// Set sets the fake clock's current time.
func (fc FakeClock) Set(t time.Time) {
	fc.currentTime = t
}
