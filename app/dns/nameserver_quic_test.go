package dns_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	. "github.com/yuzuki616/xray-core/app/dns"
	"github.com/yuzuki616/xray-core/common"
	"github.com/yuzuki616/xray-core/common/net"
	"github.com/yuzuki616/xray-core/features/dns"
)

func TestQUICNameServer(t *testing.T) {
	url, err := url.Parse("quic://dns.adguard.com")
	common.Must(err)
	s, err := NewQUICNameServer(url)
	common.Must(err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	ips, err := s.QueryIP(ctx, "google.com", net.IP(nil), dns.IPOption{
		IPv4Enable: true,
		IPv6Enable: true,
	}, false)
	cancel()
	common.Must(err)
	if len(ips) == 0 {
		t.Error("expect some ips, but got 0")
	}

	ctx2, cancel := context.WithTimeout(context.Background(), time.Second*5)
	ips2, err := s.QueryIP(ctx2, "google.com", net.IP(nil), dns.IPOption{
		IPv4Enable: true,
		IPv6Enable: true,
	}, true)
	cancel()
	common.Must(err)
	if r := cmp.Diff(ips2, ips); r != "" {
		t.Fatal(r)
	}
}
