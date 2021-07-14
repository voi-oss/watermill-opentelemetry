# Watermill OpenTelemetry integration

[![Go Report Card](https://goreportcard.com/badge/github.com/voi-oss/watermill-opentelemetry?style=flat-square)](https://goreportcard.com/report/github.com/voi-oss/watermill-opentelemetry)
[![GolangCI](https://golangci.com/badges/github.com/voi-oss/watermill-opentelemetry.svg)](https://golangci.com/r/github.com/voi-oss/watermill-opentelemetry)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/voi-oss/watermill-opentelemetry)

Bringing distributed tracing support to [Watermill](https://watermill.io/) with [OpenTelemetry](https://opentelemetry.io/). 

## Usage

### For publishers

```go
package example

import (
    "github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
    "github.com/ThreeDotsLabs/watermill/message"
    "github.com/garsue/watermillzap"
    wot "github.com/voi-oss/watermill-opentelemetry"
    "go.uber.org/zap"
)

type PublisherConfig struct {
	Name         string
	GCPProjectID string
}

// NewPublisher instantiates a GCP Pub/Sub Publisher with tracing capabilities.
func NewPublisher(logger *zap.Logger, config PublisherConfig) (message.Publisher, error) {
	publisher, err := googlecloud.NewPublisher(
        googlecloud.PublisherConfig{ProjectID: config.GCPProjectID},
        watermillzap.NewLogger(logger),
    )
	if err != nil {
		return nil, err
	}

	if config.Name == "" {
		return wot.NewPublisherDecorator(publisher), nil
	}

	return wot.NewNamedPublisherDecorator(config.Name, publisher), nil
}
```

### For subscribers

A tracing middleware can be defined at the router level:

```go
package example

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
    wotel "github.com/voi-oss/watermill-opentelemetry"
)

func InitTracedRouter() (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, watermill.NopLogger{})
	if err != nil {
		return nil, err
	}

	router.AddMiddleware(wotel.Trace())

	return router, nil
}
```

Alternatively, individual handlers can be traced: 

```go
package example

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
    wotel "github.com/voi-oss/watermill-opentelemetry"
)

func InitRouter() (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, watermill.NopLogger{})
	if err != nil {
		return nil, err
	}
    
    // subscriber definition omitted for clarity
    subscriber := (message.Subscriber)(nil)

	router.AddNoPublisherHandler(
        "handler_name",
        "subscribeTopic",
        subscriber,
        wotel.TraceNoPublishHandler(func(msg *message.Message) error {
            return nil
        }),
    )

	return router, nil
}
```

## Contributions

We encourage and support an active, healthy community of contributors &mdash;
including you! Details are in the [contribution guide](CONTRIBUTING.md) and
the [code of conduct](CODE_OF_CONDUCT.md). The maintainers keep an eye on
issues and pull requests, but you can also report any negative conduct to
opensource@voiapp.io.

### Contributors

- [@K-Phoen](https://github.com/K-Phoen)
- [@jeespers](https://github.com/jeespers)

#### I am missing?

If you feel you should be on this list, create a PR to add yourself.

## License

Apache 2.0, see [LICENSE.md](LICENSE.md).