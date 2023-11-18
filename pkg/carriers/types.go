package carriers

type Server interface {
	StartServer(listen string)
	Shutdown()
}

type Client interface {
	StartClient(server string)
	Shutdown()
}
