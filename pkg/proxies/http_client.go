package proxies

import "net/http"

func NewHttpClient(client *http.Client, dialer Dialer) *http.Client {
	client.Transport = &http.Transport{
		DialContext: dialer.NewConn,
	}
	return client
}
