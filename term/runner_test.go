package term

import "testing"

func TestRun(t *testing.T) {
	r := New()
	err := r.Run("ls", []string{})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestInteractiveRun(t *testing.T) {
	r := New()
	err := r.InteractiveRun("ls", []string{})
	if err != nil {
		t.Error(err.Error())
	}
}
