package discovery

import (
	"context"
	"fmt"
	"time"

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


func Discover(ctx context.Context) ([]*mdns.ServiceEntry, error) {
	entries := make(chan *mdns.ServiceEntry, 16)

	params := mdns.DefaultParams(ServiceType)
	params.Timeout = 3 * time.Second
	params.Entries = entries

	errCh := make(chan error, 1)

	go func() {
		err := mdns.QueryContext(ctx, params)
		close(entries)
		errCh <- err
	}()

	var discovered []*mdns.ServiceEntry

	for entry := range entries {
		discovered = append(discovered, entry)
	}

	if err := <-errCh; err != nil {
		return discovered, fmt.Errorf("discover LANtern services: %w", err)
	}

	return discovered, nil
}

