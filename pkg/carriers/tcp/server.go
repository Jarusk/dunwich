package tcp

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
)

type TcpServer struct {
	l net.Listener
}

func NewTcpServer() *TcpServer {
	return &TcpServer{}
}

func (s *TcpServer) StartServer(listen string) error {
	var err error
	s.l, err = net.Listen("tcp", listen)
	if err != nil {
		slog.Debug("was unable to create listener", "err", err)
		return fmt.Errorf("failed to create listener: %v", err)
	}

	slog.Debug("starting to listen")
	go s.serve()

	return nil
}

func (s *TcpServer) Shutdown() {
	err := s.l.Close()
	if err != nil {
		slog.Debug("closed listener")
	}
}

func (s *TcpServer) serve() {
	for {
		// Accept incoming connections
		conn, err := s.l.Accept()
		if err != nil {
			slog.Debug(
				"failed to accept incoming connection",
				slog.String("err", err.Error()),
			)
			continue
		}

		slog.Info("received new connection", slog.String("addr", conn.RemoteAddr().String()))

		// Handle client connection in a goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	slog.Debug(
		"handling connection",
		slog.String("client", conn.RemoteAddr().String()),
	)

	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			slog.Error("client hung up", "err", err)
			break
		}

		slog.Debug("received message",
			slog.String("content", msg),
			slog.String("client", conn.RemoteAddr().String()),
		)
	}
}
