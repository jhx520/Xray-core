package conf

import (
	"github.com/golang/protobuf/proto"

	"github.com/yuzuki616/xray-core/transport/internet/grpc"
)

type GRPCConfig struct {
	ServiceName         string `json:"serviceName" `
	MultiMode           bool   `json:"multiMode"`
	IdleTimeout         int32  `json:"idle_timeout"`
	HealthCheckTimeout  int32  `json:"health_check_timeout"`
	PermitWithoutStream bool   `json:"permit_without_stream"`
	InitialWindowsSize  int32  `json:"initial_windows_size"`
}

func (g *GRPCConfig) Build() (proto.Message, error) {
	if g.IdleTimeout <= 0 {
		g.IdleTimeout = 0
	}
	if g.HealthCheckTimeout <= 0 {
		g.HealthCheckTimeout = 0
	}
	if g.InitialWindowsSize < 0 {
		// default window size of gRPC-go
		g.InitialWindowsSize = 0
	}

	return &grpc.Config{
		ServiceName:         g.ServiceName,
		MultiMode:           g.MultiMode,
		IdleTimeout:         g.IdleTimeout,
		HealthCheckTimeout:  g.HealthCheckTimeout,
		PermitWithoutStream: g.PermitWithoutStream,
		InitialWindowsSize:  g.InitialWindowsSize,
	}, nil
}
