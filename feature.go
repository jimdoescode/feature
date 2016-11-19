package feature

import (
	"crypto/rand"
	"crypto/sha256"
)

// Group is the interface used by EnabledFor to consistently
// ramp a feature to a certain group like users.
type Group interface {
	// GetGroupIdentifier provides a unique identifier that can be
	// hashed to consistently maintain if a flag is enabled or not.
	GetGroupIdentifier() []byte
	// AlwaysEnabled forces any flag to be enabled for this group.
	AlwaysEnabled() bool
}

// Map a byte slice to a single floating point
// number in [0, 1)
func mapBytes(h []byte) float64 {
	l := uint(len(h))

	vmax := 1 << l
	v := 0
	for _, b := range h {
		v = v << 1
		if b >= 128 {
			v += 1
		}
	}

	return float64(v) / float64(vmax)
}

// Flag represents a feature flag
type Flag struct {
	name    string
	percent float64
	offset  []byte
}

// Create a new feature flag.
func NewFlag(name string, percent float64) *Flag {
	return &Flag{name, percent, []byte(name)}
}

// EnabledFor applies a feature flag consistently
// across a group, based on a flag's percent.
func (f *Flag) EnabledFor(g Group) bool {
	if g.AlwaysEnabled() {
		return true
	}

	h := sha256.New()
	h.Write(g.GetGroupIdentifier())
	h.Write(f.offset)

	return f.percent > mapBytes(h.Sum(nil))
}

// Enabled randomly applies a feature flag based
// on the flag's percent.
func (f *Flag) Enabled() bool {
	b := make([]byte, 32)
	rand.Read(b)

	return f.percent > mapBytes(b)
}
