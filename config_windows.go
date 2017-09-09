package wingedsnake

import (
	"github.com/go-ini/ini"
)

// config 配置文件
type config struct {
	Base struct {
		PidPath  string   // 进程pid存储位置
		Affinity []string // CPU亲和性
	}
	Client struct {
		Config string   // 配置文件路径
		Pprof  []string // pprof 监听
		Socket []string // socket 监听
		Ping   []string // ping 监听
	}
}

// newConfig 初始化 加载并解析配置文件到 Conf 对象
func newConfig(confPath string) (*config, error) {
	conf := &config{}
	conf.Base.PidPath = "c:\\"
	conf.Base.Affinity = []string{"0001"}
	if err := ini.MapTo(conf, confPath); err != nil {
		return nil, err
	}
	return conf, nil
}

// Reload 重新加载 重载配置文件到 Conf 对象
func (cfg *config) Reload(confPath string) error {
	tmp := &config{}
	if err := ini.MapTo(tmp, confPath); err != nil {
		return err
	}
	cfg = tmp
	return nil
}
