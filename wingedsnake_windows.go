package wingedsnake

import (
	"errors"
)

var (
	errWindows = errors.New("not support windows")
)

func Main(initFunc func(config string, socket []net.Listener), quitFunc func()) error {
	return errWindows
}
