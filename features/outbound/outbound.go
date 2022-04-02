package outbound

import (
	"context"

	"github.com/yuzuki616/xray-core/common"
	"github.com/yuzuki616/xray-core/features"
	"github.com/yuzuki616/xray-core/transport"
)

// Handler is the interface for handlers that process outbound connections.
//
// xray:api:stable
type Handler interface {
	common.Runnable
	Tag() string
	Dispatch(ctx context.Context, link *transport.Link)
}

type HandlerSelector interface {
	Select([]string) []string
}

// Manager is a feature that manages outbound.Handlers.
//
// xray:api:stable
type Manager interface {
	features.Feature
	// GetHandler returns an outbound.Handler for the given tag.
	GetHandler(tag string) Handler
	// GetDefaultHandler returns the default outbound.Handler. It is usually the first outbound.Handler specified in the configuration.
	GetDefaultHandler() Handler
	// AddHandler adds a handler into this outbound.Manager.
	AddHandler(ctx context.Context, handler Handler) error

	// RemoveHandler removes a handler from outbound.Manager.
	RemoveHandler(ctx context.Context, tag string) error
}

// ManagerType returns the type of Manager interface. Can be used to implement common.HasType.
//
// xray:api:stable
func ManagerType() interface{} {
	return (*Manager)(nil)
}
