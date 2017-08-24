package cproxy

type DefaultServerConnector struct {
	dialer      Dialer
	initializer Initializer
}

func NewServerConnector(dialer Dialer, initializer Initializer) *DefaultServerConnector {
	return &DefaultServerConnector{dialer: dialer, initializer: initializer}
}

func (this *DefaultServerConnector) Connect(client Socket, serverAddress string) Proxy {
	server := this.dialer.Dial(serverAddress)
	if server == nil {
		return nil
	}

	if !this.initializer.Initialize(client, server) {
		server.Close()
		return nil
	}

	return NewProxy(client, server)
}
