package dns

import (
	"context"
	"strings"

	"github.com/yuzuki616/xray-core/common/net"
	"github.com/yuzuki616/xray-core/features/dns"
	"github.com/yuzuki616/xray-core/features/dns/localdns"
)

// LocalNameServer is an wrapper over local DNS feature.
type LocalNameServer struct {
	client *localdns.Client
}

const errEmptyResponse = "No address associated with hostname"

// QueryIP implements Server.
func (s *LocalNameServer) QueryIP(_ context.Context, domain string, _ net.IP, option dns.IPOption, _ bool) (ips []net.IP, err error) {
	ips, err = s.client.LookupIP(domain, option)

	if err != nil && strings.HasSuffix(err.Error(), errEmptyResponse) {
		err = dns.ErrEmptyResponse
	}

	if len(ips) > 0 {
		newError("Localhost got answer: ", domain, " -> ", ips).AtInfo().WriteToLog()
	}

	return
}

// Name implements Server.
func (s *LocalNameServer) Name() string {
	return "localhost"
}

// NewLocalNameServer creates localdns server object for directly lookup in system DNS.
func NewLocalNameServer() *LocalNameServer {
	newError("DNS: created localhost client").AtInfo().WriteToLog()
	return &LocalNameServer{
		client: localdns.New(),
	}
}

// NewLocalDNSClient creates localdns client object for directly lookup in system DNS.
func NewLocalDNSClient() *Client {
	return &Client{server: NewLocalNameServer()}
}
