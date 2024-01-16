package vmess

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter/outbound"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies"
	"io"
	"net/http"
	"testing"
)

func TestVmess(t *testing.T) {
	dialer := NewDialer(proxies.StringResolver, []outbound.VmessOption{{
		Name:   "vmess",
		Server: "sever_host",
		Port:   10002,
		UUID:   "b831381d-6324-4d53-ad4f-8cda48b30811",
		Cipher: "auto",
		UDP:    true,
	},
		{
			Name:   "vmess",
			Server: "sever_host",
			Port:   10002,
			UUID:   "b831381d-6324-4d53-ad4f-8cda48b30811",
			Cipher: "auto",
			UDP:    true,
		},
		{
			Name:   "vmess",
			Server: "sever_host",
			Port:   10002,
			UUID:   "b831381d-6324-4d53-ad4f-8cda48b30811",
			Cipher: "auto",
			UDP:    true,
		},
		{
			Name:   "vmess",
			Server: "sever_host",
			Port:   10002,
			UUID:   "b831381d-6324-4d53-ad4f-8cda48b30811",
			Cipher: "auto",
			UDP:    true,
		},
	}...,
	)

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
