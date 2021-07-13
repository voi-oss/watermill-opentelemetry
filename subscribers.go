package middleware

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const subscriberTracerName = "watermill/subscriber"

// HandlerMiddleware decorates a watermill HandlerFunc to add tracing when a
// message is received.
func HandlerMiddleware(h message.HandlerFunc, options ...Option) message.HandlerFunc {
	tracer := otel.Tracer(subscriberTracerName)
	config := &config{}

	for _, opt := range options {
		opt(config)
	}

	spanOptions := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(config.spanAttributes...),
	}

	return func(msg *message.Message) ([]*message.Message, error) {
		spanName := message.HandlerNameFromCtx(msg.Context())
		ctx, span := tracer.Start(msg.Context(), spanName, spanOptions...)
		span.SetAttributes(
			semconv.MessagingDestinationKindTopic,
			semconv.MessagingDestinationKey.String(message.SubscribeTopicFromCtx(ctx)),
			semconv.MessagingOperationReceive,
		)
		msg.SetContext(ctx)

		events, err := h(msg)

		if err != nil {
			span.RecordError(err)
		}
		span.End()

		return events, err
	}
}

// NoPublishHandlerFuncMiddleware decorates a watermill NoPublishHandlerFunc to
// add tracing when a message is received.
func NoPublishHandlerFuncMiddleware(h message.NoPublishHandlerFunc, options ...Option) message.NoPublishHandlerFunc {
	decoratedHandler := HandlerMiddleware(func(msg *message.Message) ([]*message.Message, error) {
		return nil, h(msg)
	}, options...)

	return func(msg *message.Message) error {
		_, err := decoratedHandler(msg)

		return err
	}
}
