package proxies

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Dialer interface {
	NewConn(ctx context.Context, network, addr string) (net.Conn, error)
}

type AddrResolver func(network, addr string) (string, uint16, error)

func StringResolver(_, addr string) (string, uint16, error) {
	resolved := strings.Split(addr, ":")
	if len(resolved) != 2 {
		return "", 0, errors.New(fmt.Sprintf("StringResolver: invalid address: %s", addr))
	}
	port, err := strconv.Atoi(resolved[1])
	if err != nil {
		return "", 0, errors.New(fmt.Sprintf("StringResolver: invalid address: %s", addr))
	}

	return resolved[0], uint16(port), nil
}

func TcpResolver(network, addr string) (string, uint16, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return "", 0, err
	}
	return tcpAddr.IP.String(), uint16(tcpAddr.Port), nil
}
