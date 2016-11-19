package feature_test

import (
	"fmt"
	"github.com/jimdoescode/feature"
	"math"
	"testing"
)

type mock struct {
	uid    []byte
	always bool
}

func (u mock) GetGroupIdentifier() []byte {
	return u.uid
}

func (u mock) AlwaysEnabled() bool {
	return u.always
}

// This mock has a group identifier that should enable the "feature"
// flag at 50% and disable the "another_feature" flag at 50%
var enabled mock = mock{[]byte{48, 208, 243, 152, 32}, false}

// This mock has a group identifier that should disable the "feature" flag at 50%.
var disabled mock = mock{[]byte{12, 254, 105, 216, 171}, false}

func TestEnabledFor(t *testing.T) {
	flag := feature.NewFlag("feature", 0.5)

	if !flag.EnabledFor(enabled) {
		t.Error("Mock not reporting as enabled.")
	}

	if flag.EnabledFor(disabled) {
		t.Error("Mock reporting as enabled.")
	}

	// Verify that mocks stay in the same buckets on subsequent calls.
	if !flag.EnabledFor(enabled) {
		t.Error("Mock not reporting as enabled on second attempt.")
	}

	if flag.EnabledFor(disabled) {
		t.Error("Mock reporting as enabled on second attempt.")
	}
}

func TestDifferentFlagEnabledFor(t *testing.T) {
	flag := feature.NewFlag("another_feature", 0.5)

	// Verify that different flags enable different things.
	if flag.EnabledFor(enabled) {
		t.Error("Mock is reporting as enabled when it shouldn't.")
	}
}

func TestAlwaysEnabledFor(t *testing.T) {
	flag := feature.NewFlag("feature", 0.5)

	if flag.EnabledFor(disabled) {
		t.Error("Mock reporting as enabled.")
	}

	// This mock has a group identifier that should disable the "test" flag
	// at 50% but with the always enabled flag it will still report as enabled
	always := mock{disabled.GetGroupIdentifier(), true}
	if !flag.EnabledFor(always) {
		t.Error("Always Enabled mock not reporting as enabled.")
	}
}

// Note: This test might be flakey since it's using completely random
// values. They should converge on the flag percentage but might be
// over the desired tolerance sometimes.
func TestEnabled(t *testing.T) {
	percent := 0.75
	tolerance := 0.01
	flag := feature.NewFlag("random_feature", percent)
	max := 100000
	count := 0

	for i := 0; i < max; i++ {
		if flag.Enabled() {
			count++
		}
	}

	hits := float64(count) / float64(max)

	if diff := math.Abs(hits - percent); diff > tolerance {
		t.Error(fmt.Sprintf("flag.Enabled exceeds tolerance of %f. Difference is %f after %d executions.", tolerance, diff, max))
	}
}
