// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package wingedsnake

import (
	"github.com/go-ini/ini"
)

// config 配置文件
type config struct {
	Base struct {
		User     string   // 用户名
		Group    string   // 用户组
		PidPath  string   // 进程pid存储位置
		Affinity []string // CPU亲和性
		Process  int      // 进程数
	}
	Client struct {
		Config  string   // 配置文件路径
		Sockets []string // socket 监听
	}
}

// newConfig 初始化 加载并解析配置文件到 Conf 对象
func newConfig(confPath string) (*config, error) {
	conf := &config{}
	conf.Base.User = "nobody"
	conf.Base.Group = "nobody"
	conf.Base.PidPath = "/tmp"
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
