package shadowsocks

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter/outbound"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies"
	"io"
	"net/http"
	"testing"
)

func TestSS(t *testing.T) {
	dialer := NewDialer(proxies.StringResolver, []outbound.ShadowSocksOption{{
		Name:     "ss",
		Server:   "sever_host",
		Port:     12001,
		Password: "xxx",
		Cipher:   "aes-128-gcm",
		UDP:      true,
		Plugin:   "obfs",
		PluginOpts: map[string]any{
			"mode": "http",
			"host": "xxx.download.windowsupdate.com",
		}},
		{
			Name:     "ss",
			Server:   "sever_host",
			Port:     12001,
			Password: "xxx",
			Cipher:   "aes-128-gcm",
			UDP:      true,
			Plugin:   "obfs",
			PluginOpts: map[string]any{
				"mode": "http",
				"host": "xxx.download.windowsupdate.com",
			}},
		{
			Name:     "ss",
			Server:   "sever_host",
			Port:     12001,
			Password: "xxx",
			Cipher:   "aes-128-gcm",
			UDP:      true,
			Plugin:   "obfs",
			PluginOpts: map[string]any{
				"mode": "http",
				"host": "xxx.download.windowsupdate.com",
			}},
		{
			Name:     "ss",
			Server:   "sever_host",
			Port:     12001,
			Password: "xxx",
			Cipher:   "aes-128-gcm",
			UDP:      true,
			Plugin:   "obfs",
			PluginOpts: map[string]any{
				"mode": "http",
				"host": "xxx.download.windowsupdate.com",
			}},
	}...)

	httpTransport := &http.Transport{
		DialContext: dialer.NewConn,
	}
	httpC := &http.Client{
		Transport: httpTransport,
	}
	resp, err := httpC.Get("https://www.youtube.com/")
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func TestProxiesYaml_Load(t *testing.T) {
	cfg := new(ProxiesYaml)
	err := cfg.Load("ss.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	outbounds := cfg.CovertOption()
	fmt.Println(outbounds)
}
