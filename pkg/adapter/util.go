package adapter

import "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

// BuildAddress returns a SocketAddress with the given ip and port.
func BuildAddress(host string, port uint32) core.Address {
	return core.Address{
		Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				Address: host,
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: port,
				},
			},
		},
	}
}