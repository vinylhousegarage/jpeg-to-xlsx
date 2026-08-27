package storage

import "testing"

func TestBuildObjectKey(t *testing.T) {
	t.Parallel()

	got := BuildObjectKey("001")
	want := "001.jpg"

	if got != want {
		t.Errorf("BuildObjectKey() = %q, want %q", got, want)
	}
}
