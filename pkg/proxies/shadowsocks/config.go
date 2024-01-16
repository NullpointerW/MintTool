package shadowsocks

import (
	"github.com/Dreamacro/clash/adapter/outbound"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

type ProxiesYaml struct {
	Proxies []Config `yaml:"proxies"`
}

type Config struct {
	Server     string         `yaml:"server"`
	Port       any            `yaml:"port"`
	Password   string         `yaml:"password"`
	Cipher     string         `yaml:"cipher"`
	UDP        bool           `yaml:"udp"`
	Plugin     string         `yaml:"plugin"`
	PluginOpts map[string]any `yaml:"plugin-opts"`
}

func (yamlCfg *ProxiesYaml) Load(fp string) error {
	b, err := os.ReadFile(fp)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, yamlCfg)
	return err
}
func (yamlCfg *ProxiesYaml) CovertOption() []outbound.ShadowSocksOption {
	var ssOps []outbound.ShadowSocksOption
	for _, p := range yamlCfg.Proxies {
		var ssOp = outbound.ShadowSocksOption{
			Name:       "ss",
			Server:     p.Server,
			Password:   p.Password,
			Cipher:     p.Cipher,
			UDP:        p.UDP,
			Plugin:     p.Plugin,
			PluginOpts: p.PluginOpts,
		}
		if portStr, ok := p.Port.(string); ok {
			ssOp.Port, _ = strconv.Atoi(portStr)
		} else {
			ssOp.Port = p.Port.(int)
		}
		ssOps = append(ssOps, ssOp)
	}
	return ssOps
}
