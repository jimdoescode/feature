package feature

import (
	"crypto/rand"
	"crypto/sha256"
)

// Group is the interface used by EnabledFor to consistently
// enable or disable a feature to a certain group, like users.
type Group interface {
	// GetGroupIdentifier provides a unique identifier that can be
	// hashed to consistently maintain if a flag is enabled or not.
	GetGroupIdentifier() []byte
	// AlwaysEnabled forces any flag to be enabled for this group.
	AlwaysEnabled() bool
}

// A sample represents a grouping of things. The size is a percentage
// of those things. 1.0f means all the things, 0.0f means none of the
// things.
type sample struct {
	size float64
}

// Determines if the byte slice is within the sample
func (s sample) Includes(h []byte) bool {
	l := len(h)
	// 40 bytes is sufficient for our calculation
	if l > 40 {
		h = h[:40]
		l = 40
	}

	vmax := 1 << uint(l)
	v := 0
	for _, b := range h {
		v = v << 1
		if b >= 128 {
			v += 1
		}
	}

	return s.size > (float64(v) / float64(vmax))
}

// Flag represents a feature flag
type Flag struct {
	sample
	name   string
	offset []byte
}

// Create a new feature flag.
func NewFlag(name string, threshold float64) *Flag {
	return &Flag{sample{threshold}, name, []byte(name)}
}

// EnabledFor applies a feature flag consistently
// across a sample group, based on a flag's sample size.
func (f *Flag) EnabledFor(g Group) bool {
	if g.AlwaysEnabled() {
		return true
	}

	h := sha256.New()
	h.Write(g.GetGroupIdentifier())
	h.Write(f.offset)

	return f.Includes(h.Sum(nil))
}

// Enabled randomly applies a feature flag based
// on the flag's percent.
func (f *Flag) Enabled() bool {
	b := make([]byte, 40)
	rand.Read(b)

	return f.Includes(b)
}
