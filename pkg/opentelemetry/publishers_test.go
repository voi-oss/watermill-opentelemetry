package opentelemetry

import (
	"bytes"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel/attribute"
)

func TestNewPublisherDecorator(t *testing.T) {
	var (
		uuid    = "0d5427ea-7ab4-4ef1-b80d-0a22bd54a98f"
		payload = message.Payload("test payload")

		pub = &mockMessagePublisher{
			PublishFunc: func(topic string, messages ...*message.Message) error {
				if got, want := topic, "test.topic"; got != want {
					t.Fatalf("topic = %q, want %q", got, want)
				}

				if got, want := len(messages), 1; got != want {
					t.Fatalf("len(messages) = %d, want %d", got, want)
				}

				message := messages[0]

				if got, want := message.UUID, uuid; got != want {
					t.Fatalf("message.UUID = %q, want %q", got, want)
				}

				if !bytes.Equal(payload, message.Payload) {
					t.Fatalf("unexpected payload")
				}

				return nil
			},
		}

		dec = NewPublisherDecorator(pub, WithSpanAttributes(
			attribute.Bool("test", true),
		))
	)

	pd, ok := dec.(*PublisherDecorator)
	if !ok {
		t.Fatalf("expected message.Publisher to be *PublisherDecorator")
	}

	if got, want := len(pd.config.spanAttributes), 1; got != want {
		t.Fatalf("len(pd.config.spanAttributes) = %d, want %d", got, want)
	}

	msg := message.NewMessage(uuid, payload)

	if err := dec.Publish("test.topic", msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

var _ message.Publisher = &mockMessagePublisher{}

type mockMessagePublisher struct {
	CloseFunc   func() error
	PublishFunc func(topic string, messages ...*message.Message) error
}

func (mock *mockMessagePublisher) Close() error {
	if mock.CloseFunc == nil {
		panic("MessagePublisher.CloseFunc: method is nil but Publisher.Close was just called")
	}

	return mock.CloseFunc()
}

func (mock *mockMessagePublisher) Publish(topic string, messages ...*message.Message) error {
	if mock.PublishFunc == nil {
		panic("MessagePublisher.PublishFunc: method is nil but Publisher.Publish was just called")
	}

	return mock.PublishFunc(topic, messages...)
}
