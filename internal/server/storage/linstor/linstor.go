package linstor

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	linstorClient "github.com/LINBIT/golinstor/client"

	"github.com/lxc/incus/v6/shared/logger"
)

// Client represents an HTTP Linstor client.
type Client struct {
	Client *linstorClient.Client
}

// NewClient initializes a new Linstor client.
func NewClient(controllerConnection, sslCACert, sslClientCert, sslClientKey string) (*Client, error) {
	logger.Info("Creating new Linstor client", logger.Ctx{"controllerConnection": controllerConnection})

	// Configure the client HTTP transport.
	httpTransport := &http.Transport{}

	// If a CA cert is provided, use it to validate the server certificates.
	if sslCACert != "" {
		rootCAs := x509.NewCertPool()
		certBlock, _ := pem.Decode([]byte(sslCACert))
		caCert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Failed to create Linstor client: %w", err)
		}

		rootCAs.AddCert(caCert)
		httpTransport.TLSClientConfig = &tls.Config{RootCAs: rootCAs}
	}

	// If a client certificate and key pair is provided, submit it to the server.
	if sslClientCert != "" && sslClientKey != "" {
		clientCert, err := tls.X509KeyPair([]byte(sslClientCert), []byte(sslClientKey))
		if err != nil {
			return nil, fmt.Errorf("Failed to create Linstor client: %w", err)
		}

		httpTransport.TLSClientConfig.Certificates = []tls.Certificate{clientCert}
	}

	// Setup the Linstor client.
	httpClient := &http.Client{Transport: httpTransport}
	controllerUrls := strings.Split(controllerConnection, ",")
	c, err := linstorClient.NewClient(linstorClient.Controllers(controllerUrls), linstorClient.HTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("Failed to create Linstor client: %w", err)
	}

	// Get the controller version to check connection.
	ctx := context.TODO()
	version, err := c.Controller.GetVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Linstor client: %w", err)
	}

	logger.Info("Connected to Linstor Controller", logger.Ctx{"version": version})
	return &Client{Client: c}, nil
}
