package scenarios

import (
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/yuzuki616/xray-core/app/log"
	"github.com/yuzuki616/xray-core/app/proxyman"
	"github.com/yuzuki616/xray-core/common"
	clog "github.com/yuzuki616/xray-core/common/log"
	"github.com/yuzuki616/xray-core/common/net"
	"github.com/yuzuki616/xray-core/common/protocol"
	"github.com/yuzuki616/xray-core/common/serial"
	"github.com/yuzuki616/xray-core/common/uuid"
	"github.com/yuzuki616/xray-core/core"
	"github.com/yuzuki616/xray-core/proxy/dokodemo"
	"github.com/yuzuki616/xray-core/proxy/freedom"
	"github.com/yuzuki616/xray-core/proxy/vmess"
	"github.com/yuzuki616/xray-core/proxy/vmess/inbound"
	"github.com/yuzuki616/xray-core/proxy/vmess/outbound"
	"github.com/yuzuki616/xray-core/testing/servers/tcp"
	"github.com/yuzuki616/xray-core/testing/servers/udp"
)

func TestDokodemoTCP(t *testing.T) {
	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	common.Must(err)
	defer tcpServer.Close()

	userID := protocol.NewID(uuid.New())
	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortList: &net.PortList{Range: []*net.PortRange{net.SinglePortRange(serverPort)}},
					Listen:   net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&inbound.Config{
					User: []*protocol.User{
						{
							Account: serial.ToTypedMessage(&vmess.Account{
								Id: userID.String(),
							}),
						},
					},
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}
	server, err := InitializeServerConfig(serverConfig)
	common.Must(err)
	defer CloseServer(server)

	clientPortRange := uint32(5)
	retry := 1
	clientPort := uint32(tcp.PickPort())
	for {
		clientConfig := &core.Config{
			App: []*serial.TypedMessage{
				serial.ToTypedMessage(&log.Config{
					ErrorLogLevel: clog.Severity_Debug,
					ErrorLogType:  log.LogType_Console,
				}),
			},
			Inbound: []*core.InboundHandlerConfig{
				{
					ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
						PortList: &net.PortList{Range: []*net.PortRange{{From: clientPort, To: clientPort + clientPortRange}}},
						Listen:   net.NewIPOrDomain(net.LocalHostIP),
					}),
					ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
						Address: net.NewIPOrDomain(dest.Address),
						Port:    uint32(dest.Port),
						NetworkList: &net.NetworkList{
							Network: []net.Network{net.Network_TCP},
						},
					}),
				},
			},
			Outbound: []*core.OutboundHandlerConfig{
				{
					ProxySettings: serial.ToTypedMessage(&outbound.Config{
						Receiver: []*protocol.ServerEndpoint{
							{
								Address: net.NewIPOrDomain(net.LocalHostIP),
								Port:    uint32(serverPort),
								User: []*protocol.User{
									{
										Account: serial.ToTypedMessage(&vmess.Account{
											Id: userID.String(),
										}),
									},
								},
							},
						},
					}),
				},
			},
		}

		server, _ := InitializeServerConfig(clientConfig)
		if server != nil && WaitConnAvailableWithTest(t, testTCPConn(net.Port(clientPort), 1024, time.Second*2)) {
			defer CloseServer(server)
			break
		}
		retry++
		if retry > 5 {
			t.Fatal("All attempts failed to start client")
		}
		clientPort = uint32(tcp.PickPort())
	}

	for port := clientPort; port <= clientPort+clientPortRange; port++ {
		if err := testTCPConn(net.Port(port), 1024, time.Second*2)(); err != nil {
			t.Error(err)
		}
	}
}

func TestDokodemoUDP(t *testing.T) {
	udpServer := udp.Server{
		MsgProcessor: xor,
	}
	dest, err := udpServer.Start()
	common.Must(err)
	defer udpServer.Close()

	userID := protocol.NewID(uuid.New())
	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortList: &net.PortList{Range: []*net.PortRange{net.SinglePortRange(serverPort)}},
					Listen:   net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&inbound.Config{
					User: []*protocol.User{
						{
							Account: serial.ToTypedMessage(&vmess.Account{
								Id: userID.String(),
							}),
						},
					},
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}
	server, err := InitializeServerConfig(serverConfig)
	common.Must(err)
	defer CloseServer(server)

	clientPortRange := uint32(5)
	retry := 1
	clientPort := uint32(udp.PickPort())
	for {
		clientConfig := &core.Config{
			Inbound: []*core.InboundHandlerConfig{
				{
					ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
						PortList: &net.PortList{Range: []*net.PortRange{{From: clientPort, To: clientPort + clientPortRange}}},
						Listen:   net.NewIPOrDomain(net.LocalHostIP),
					}),
					ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
						Address: net.NewIPOrDomain(dest.Address),
						Port:    uint32(dest.Port),
						NetworkList: &net.NetworkList{
							Network: []net.Network{net.Network_UDP},
						},
					}),
				},
			},
			Outbound: []*core.OutboundHandlerConfig{
				{
					ProxySettings: serial.ToTypedMessage(&outbound.Config{
						Receiver: []*protocol.ServerEndpoint{
							{
								Address: net.NewIPOrDomain(net.LocalHostIP),
								Port:    uint32(serverPort),
								User: []*protocol.User{
									{
										Account: serial.ToTypedMessage(&vmess.Account{
											Id: userID.String(),
										}),
									},
								},
							},
						},
					}),
				},
			},
		}

		server, _ := InitializeServerConfig(clientConfig)
		if server != nil && WaitConnAvailableWithTest(t, testUDPConn(net.Port(clientPort), 1024, time.Second*2)) {
			defer CloseServer(server)
			break
		}
		retry++
		if retry > 5 {
			t.Fatal("All attempts failed to start client")
		}
		clientPort = uint32(udp.PickPort())
	}

	var errg errgroup.Group
	for port := clientPort; port <= clientPort+clientPortRange; port++ {
		errg.Go(testUDPConn(net.Port(port), 1024, time.Second*5))
	}
	if err := errg.Wait(); err != nil {
		t.Error(err)
	}
}
