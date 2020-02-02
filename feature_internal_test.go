package feature

import (
	"testing"
)

func TestMinSample(t *testing.T) {
	b := []byte{0, 0, 0, 0}
	if !includes(1.0, b) {
		t.Error("sample.include should include min value of 0 at 100% size")
	}

	if includes(0.0, b) {
		t.Error("sample.include should NOT include min value of 0 at 0% size")
	}
}

func TestMaxSample(t *testing.T) {
	b := []byte{255, 255, 255, 255}
	if !includes(1.0, b) {
		t.Error("sample.include should include max value of 255 at 100% size")
	}

	if includes(0.0, b) {
		t.Error("sample.include should NOT include max value of 255 at 0% size")
	}
}
