package vmess

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
	Server string `yaml:"server"`
	Port   any    `yaml:"port"`
	UUID   string `yaml:"uuid"`
	Cipher string `yaml:"cipher"`
	UDP    bool   `yaml:"udp"`
}

func (yamlCfg *ProxiesYaml) Load(fp string) error {
	b, err := os.ReadFile(fp)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, yamlCfg)
	return err
}
func (yamlCfg *ProxiesYaml) CovertOption() []outbound.VmessOption {
	var vmOps []outbound.VmessOption
	for _, p := range yamlCfg.Proxies {
		var vmOp = outbound.VmessOption{
			Name:   "vmess",
			Server: p.Server,
			UUID:   p.UUID,
			Cipher: p.Cipher,
			UDP:    p.UDP,
		}
		if portStr, ok := p.Port.(string); ok {
			vmOp.Port, _ = strconv.Atoi(portStr)
		} else {
			vmOp.Port = p.Port.(int)
		}
		vmOps = append(vmOps, vmOp)
	}
	return vmOps
}
