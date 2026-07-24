package embedding

import "testing"

func TestVectorString(t *testing.T) {
	if got := VectorString([]float64{1, -0.25, 3.5}); got != "[1,-0.25,3.5]" {
		t.Fatalf("got %q", got)
	}
}
