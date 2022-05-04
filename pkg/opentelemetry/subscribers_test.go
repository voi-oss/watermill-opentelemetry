package opentelemetry

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel/attribute"
)

func TestTrace(t *testing.T) {
	middleware := Trace(WithSpanAttributes(
		attribute.Bool("test", true),
	))

	var (
		uuid    = "52219531-0cd8-4b64-be31-ba6b4ef01472"
		payload = message.Payload("test payload for Trace")
		msg     = message.NewMessage(uuid, payload)
	)

	h := func(m *message.Message) ([]*message.Message, error) {
		if got, want := m.UUID, uuid; got != want {
			t.Fatalf("m.UUID = %q, want %q", got, want)
		}

		return nil, nil
	}

	if _, err := middleware(h)(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTraceNoPublishHandler(t *testing.T) {
	var (
		uuid    = "88b433e5-12fa-4eb7-9229-6bfd67de5c4f"
		payload = message.Payload("test payload for TraceNoPublishHandler")
		msg     = message.NewMessage(uuid, payload)
	)

	h := func(m *message.Message) error {
		if got, want := m.UUID, uuid; got != want {
			t.Fatalf("m.UUID = %q, want %q", got, want)
		}

		return nil
	}

	handlerFunc := TraceNoPublishHandler(h, WithSpanAttributes(
		attribute.Bool("test", true),
	))

	if err := handlerFunc(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
