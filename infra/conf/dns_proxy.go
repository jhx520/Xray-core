package conf

import (
	"github.com/golang/protobuf/proto"

	"github.com/yuzuki616/xray-core/common/net"
	"github.com/yuzuki616/xray-core/proxy/dns"
)

type DNSOutboundConfig struct {
	Network   Network  `json:"network"`
	Address   *Address `json:"address"`
	Port      uint16   `json:"port"`
	UserLevel uint32   `json:"userLevel"`
}

func (c *DNSOutboundConfig) Build() (proto.Message, error) {
	config := &dns.Config{
		Server: &net.Endpoint{
			Network: c.Network.Build(),
			Port:    uint32(c.Port),
		},
		UserLevel: c.UserLevel,
	}
	if c.Address != nil {
		config.Server.Address = c.Address.Build()
	}
	return config, nil
}
