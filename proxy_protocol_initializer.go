package cproxy

import (
	"fmt"
	"io"
	"net"
)

type ProxyProtocolInitializer struct {
}

func NewProxyProtocolInitializer() *ProxyProtocolInitializer {
	return &ProxyProtocolInitializer{}
}

func (this *ProxyProtocolInitializer) Initialize(client, server Socket) bool {
	header := formatHeader(client.RemoteAddr(), server.RemoteAddr())
	_, err := io.WriteString(server, header)
	return err == nil
}
func formatHeader(client, server net.Addr) string {
	clientAddress, clientPort := parseAddress(client.String())
	serverAddress, serverPort := parseAddress(server.String())
	return fmt.Sprintf(preambleFormat, clientAddress, serverAddress, clientPort, serverPort)
}
func parseAddress(address string) (string, string) {
	address, port, _ := net.SplitHostPort(address)
	return address, port
}

const preambleFormat = "PROXY TCP4 %s %s %s %s\r\n" // TODO (TCP4 vs TCP6)
