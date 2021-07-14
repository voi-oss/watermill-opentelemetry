package opentelemetry

import (
	"fmt"
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const publisherTracerName = "watermill/publisher"

// PublisherDecorator decorates a standard watermill publisher to add tracing capabilities.
type PublisherDecorator struct {
	pub           message.Publisher
	publisherName string
	config        config
	tracer        trace.Tracer
}

// NewPublisherDecorator instantiates a PublisherDecorator with a default name.
func NewPublisherDecorator(pub message.Publisher, options ...Option) message.Publisher {
	return NewNamedPublisherDecorator(structName(pub), pub, options...)
}

// NewNamedPublisherDecorator instantiates a PublisherDecorator with a provided name.
func NewNamedPublisherDecorator(name string, pub message.Publisher, options ...Option) message.Publisher {
	config := config{}

	for _, opt := range options {
		opt(&config)
	}

	return &PublisherDecorator{
		pub:           pub,
		publisherName: name,
		config:        config,
		tracer:        otel.Tracer(publisherTracerName),
	}
}

// Publish implements the watermill Publisher interface and creates traces.
// Publishing of messages are delegated to the decorated Publisher.
func (p *PublisherDecorator) Publish(topic string, messages ...*message.Message) error {
	if len(messages) == 0 {
		return nil
	}

	ctx := messages[0].Context()
	spanName := message.PublisherNameFromCtx(ctx)
	if spanName == "" {
		spanName = p.publisherName
	}

	ctx, span := p.tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindProducer))
	messages[0].SetContext(ctx)

	spanAttributes := []attribute.KeyValue{
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationKey.String(topic),
		semconv.MessagingOperationProcess,
	}
	spanAttributes = append(spanAttributes, spanAttributes...)
	span.SetAttributes(spanAttributes...)

	err := p.pub.Publish(topic, messages...)
	if err != nil {
		span.RecordError(err)
	}

	span.End()

	return err
}

// Close implements the watermill Publisher interface.
func (p *PublisherDecorator) Close() error {
	return p.pub.Close()
}

func structName(v interface{}) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}

	s := fmt.Sprintf("%T", v)
	// trim the pointer marker, if any
	return strings.TrimLeft(s, "*")
}
