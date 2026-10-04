package objects

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewNotFound(t *testing.T) {
	err := NewNotFound("Pod", "web")
	if err == nil {
		t.Fatal("NewNotFound returned nil")
	}
	if got, want := err.Error(), `Pod "web" not found`; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	if err.Code != 404 || err.Reason != ReasonNotFound {
		t.Errorf("got Code=%d Reason=%q, want 404 %q", err.Code, err.Reason, ReasonNotFound)
	}
}

func TestNewAlreadyExists(t *testing.T) {
	err := NewAlreadyExists("ConfigMap", "cfg")
	if err == nil {
		t.Fatal("NewAlreadyExists returned nil")
	}
	if got, want := err.Error(), `ConfigMap "cfg" already exists`; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	if err.Code != 409 || err.Reason != ReasonAlreadyExists {
		t.Errorf("got Code=%d Reason=%q, want 409 %q", err.Code, err.Reason, ReasonAlreadyExists)
	}
}

func TestErrorHelpers(t *testing.T) {
	notFound := NewNotFound("Pod", "web")
	tests := []struct {
		name          string
		err           error
		reason        StatusReason
		notFound      bool
		alreadyExists bool
	}{
		{"not found", notFound, ReasonNotFound, true, false},
		{"already exists", NewAlreadyExists("Pod", "web"), ReasonAlreadyExists, false, true},
		{"wrapped once", fmt.Errorf("getting pod: %w", notFound), ReasonNotFound, true, false},
		{"wrapped twice", fmt.Errorf("reconcile: %w", fmt.Errorf("getting pod: %w", notFound)), ReasonNotFound, true, false},
		{"wrapped with %v loses the chain", fmt.Errorf("getting pod: %v", notFound), "", false, false},
		{"plain error with same text", errors.New(`Pod "web" not found`), "", false, false},
		{"nil error", nil, "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReasonForError(tt.err); got != tt.reason {
				t.Errorf("ReasonForError() = %q, want %q", got, tt.reason)
			}
			if got := IsNotFound(tt.err); got != tt.notFound {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.notFound)
			}
			if got := IsAlreadyExists(tt.err); got != tt.alreadyExists {
				t.Errorf("IsAlreadyExists() = %v, want %v", got, tt.alreadyExists)
			}
		})
	}
}
