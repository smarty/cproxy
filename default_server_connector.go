package cproxy

type defaultServerConnector struct {
	dialer      dialer
	initializer initializer
}

func newServerConnector(dialer dialer, initializer initializer) *defaultServerConnector {
	return &defaultServerConnector{dialer: dialer, initializer: initializer}
}

func (this *defaultServerConnector) Connect(client socket, serverAddress string) proxy {
	server := this.dialer.Dial(serverAddress)
	if server == nil {
		return nil
	}

	if !this.initializer.Initialize(client, server) {
		_ = server.Close()
		return nil
	}

	return newProxy(client, server)
}
