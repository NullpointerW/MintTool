package shadowsocks

import (
	"context"
	"fmt"
	"github.com/Dreamacro/clash/adapter/outbound"
	C "github.com/Dreamacro/clash/constant"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies"
	"golang.org/x/exp/rand"
	"net"
	"time"
)

type Dialer struct {
	nodes        []outbound.ShadowSocksOption
	nl           int
	addrResolver proxies.AddrResolver
}

func NewDialer(addrResolver proxies.AddrResolver, nodes ...outbound.ShadowSocksOption) *Dialer {
	dl := new(Dialer)
	dl.nodes = nodes
	dl.nl = len(nodes)
	dl.addrResolver = addrResolver
	return dl
}

func (dialer *Dialer) NewConn(ctx context.Context, network, addr string) (net.Conn, error) {
	rand.Seed(uint64(time.Now().UnixNano()))
	r := rand.Int() % len(dialer.nodes)
	fmt.Println("use node::", r)
	proxy, err := outbound.NewShadowSocks(dialer.nodes[r])
	if err != nil {
		return nil, err
	}
	host, port, err := dialer.addrResolver(network, addr)
	if err != nil {
		return nil, err
	}
	conn, err := proxy.DialContext(context.Background(), &C.Metadata{
		Host:    host,
		DstPort: C.Port(port),
	})
	return conn, err
}
