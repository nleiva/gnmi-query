package main

import (
	"context"
	"crypto/tls"
	"fmt"

	gnmi "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func initClient(ctx context.Context, targetAddr, user, pass string) (*ygnmi.Client, error) {
	skipVerify := true

	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipVerify,
		MinVersion:         tls.VersionTLS12,
	}
	creds := credentials.NewTLS(tlsConfig)

	var opts []grpc.DialOption

	//Add TLS credentials to config options array.
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// Add user/password to config options array.
	opts = append(opts, grpc.WithPerRPCCredentials(
		&loginCreds{
			Username:   user,
			Password:   pass,
			requireTLS: true,
		},
	),
	)

	conn, err := grpc.NewClient(targetAddr, opts...)

	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}
	// Send Capabilities Request to test connection
	gnmiClient := gnmi.NewGNMIClient(conn)
	resp, err := gnmiClient.Capabilities(ctx, &gnmi.CapabilityRequest{})
	if err != nil {
		return nil, fmt.Errorf("capabilities RPC failed: %w", err)
	}

	fmt.Println("Supported Encodings:")
	for _, enc := range resp.SupportedEncodings {
		fmt.Printf("  - %s\n", enc.String())
	}
	fmt.Printf("gNMI Version: %s\n\n", resp.GNMIVersion)

	yclient, err := ygnmi.NewClient(
		gnmiClient,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating ygnmi client: %w", err)
	}

	return yclient, nil
}
