package tcp

import (
	"context"

	"github.com/yuzuki616/xray-core/common"
	"github.com/yuzuki616/xray-core/common/net"
	"github.com/yuzuki616/xray-core/common/session"
	"github.com/yuzuki616/xray-core/transport/internet"
	"github.com/yuzuki616/xray-core/transport/internet/stat"
	"github.com/yuzuki616/xray-core/transport/internet/tls"
	"github.com/yuzuki616/xray-core/transport/internet/xtls"
)

// Dial dials a new TCP connection to the given destination.
func Dial(ctx context.Context, dest net.Destination, streamSettings *internet.MemoryStreamConfig) (stat.Connection, error) {
	newError("dialing TCP to ", dest).WriteToLog(session.ExportIDToError(ctx))
	conn, err := internet.DialSystem(ctx, dest, streamSettings.SocketSettings)
	if err != nil {
		return nil, err
	}

	if config := tls.ConfigFromStreamSettings(streamSettings); config != nil {
		tlsConfig := config.GetTLSConfig(tls.WithDestination(dest))
		if fingerprint, ok := tls.Fingerprints[config.Fingerprint]; ok {
			conn = tls.UClient(conn, tlsConfig, fingerprint)
			if err := conn.(*tls.UConn).Handshake(); err != nil {
				return nil, err
			}
		} else {
			conn = tls.Client(conn, tlsConfig)
		}
	} else if config := xtls.ConfigFromStreamSettings(streamSettings); config != nil {
		xtlsConfig := config.GetXTLSConfig(xtls.WithDestination(dest))
		conn = xtls.Client(conn, xtlsConfig)
	}

	tcpSettings := streamSettings.ProtocolSettings.(*Config)
	if tcpSettings.HeaderSettings != nil {
		headerConfig, err := tcpSettings.HeaderSettings.GetInstance()
		if err != nil {
			return nil, newError("failed to get header settings").Base(err).AtError()
		}
		auth, err := internet.CreateConnectionAuthenticator(headerConfig)
		if err != nil {
			return nil, newError("failed to create header authenticator").Base(err).AtError()
		}
		conn = auth.Client(conn)
	}
	return stat.Connection(conn), nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
