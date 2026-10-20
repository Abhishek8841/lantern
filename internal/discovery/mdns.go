package discovery

import (
	"fmt"

	"github.com/hashicorp/mdns"
)

const ServiceType = "_lantern._tcp"

func Advertise(instanceName string, port int, peerID string) (*mdns.Server, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instance name cannot be empty")
	}

	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid TCP port: %d", port)
	}

	if peerID == "" {
		return nil, fmt.Errorf("peer ID cannot be empty")
	}

	service, err := mdns.NewMDNSService(
		instanceName,
		ServiceType,
		"",
		"",
		port,
		nil,
		[]string{
			"id=" + peerID,
			"version=1",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create mDNS service: %w", err)
	}

	server, err := mdns.NewServer(&mdns.Config{
		Zone: service,
	})
	if err != nil {
		return nil, fmt.Errorf("start mDNS server: %w", err)
	}

	return server, nil
}
