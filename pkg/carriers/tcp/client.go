package tcp

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"

	"github.com/Jarusk/dunwich/pkg/corpus/dunwich"
)

type TcpClient struct {
	c net.Conn
}

func NewTcpClient() *TcpClient {
	return &TcpClient{}
}

func (c *TcpClient) StartClient(server string) error {
	var err error
	c.c, err = net.Dial("tcp", server)
	if err != nil {
		slog.Debug("unable to start client", "err", err)
		return fmt.Errorf("failed to start client: %v", err)
	}

	c.client()

	return nil
}

func (c *TcpClient) Shutdown() {
	err := c.c.Close()
	if err != nil {
		slog.Debug("client shutdown")
	}
}

func (c *TcpClient) client() {

	w := bufio.NewWriter(c.c)
	defer w.Flush()

	for {
		for i := 0; i < dunwich.GetNumSegments(); i++ {
			segment, err := dunwich.GetSegment(i)
			if err != nil {
				slog.Error("failed to get segment", "id", i, "err", err)
				break
			}
			count, err := w.WriteString(*segment)
			if err != nil {
				slog.Error("failed to write segment", "id", i, "err", err)
				return
			}

			slog.Debug("sent bytes", slog.Int("count", count))
		}
	}
}
