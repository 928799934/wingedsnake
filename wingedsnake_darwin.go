package wingedsnake

import (
	"errors"
	"net"
)

var (
	errWindows = errors.New("not support darwin")
)

// Main Main
func Main(initFunc func(config string, socket []net.Listener), quitFunc func()) error {
	return errWindows
}
