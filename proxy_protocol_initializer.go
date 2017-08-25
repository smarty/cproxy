package cproxy

import (
	"fmt"
	"io"
	"net"
	"strings"
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
	if strings.Contains(clientAddress, ":") {
		return fmt.Sprintf(proxyProtocolIPv6Preamble, clientAddress, serverAddress, clientPort, serverPort)
	} else {
		return fmt.Sprintf(proxyProtocolIPv4Preamble, clientAddress, serverAddress, clientPort, serverPort)
	}
}
func parseAddress(address string) (string, string) {
	address, port, _ := net.SplitHostPort(address)
	return address, port
}

const proxyProtocolIPv4Preamble = "PROXY TCP4 %s %s %s %s\r\n"
const proxyProtocolIPv6Preamble = "PROXY TCP6 %s %s %s %s\r\n"
