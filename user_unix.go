// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package wingedsnake

import (
	"os/user"
	"strconv"
	"syscall"
)

// getConfigUser 修改进程uid gid
func getConfigUser(conf *config) (int, int, error) {
	// 获取 user 的 uid
	ui, err := user.Lookup(conf.Base.User)
	if err != nil {
		logf("user.Lookup(%v) error(%v)", conf.Base.User, err)
		return 0, 0, err
	}
	// 获取 group 的 gid
	gi, err := user.LookupGroup(conf.Base.Group)
	if err != nil {
		logf("user.LookupGroup(%v) error(%v)", conf.Base.Group, err)
		return 0, 0, err
	}
	uid, _ := strconv.Atoi(ui.Uid)
	gid, _ := strconv.Atoi(gi.Gid)
	return uid, gid, nil
}

func exchangeOwner(uid, gid int) error {
	// 修改 进程 uid
	if err := syscall.Setregid(gid, gid); err != nil {
		logf("syscall.Setregid(%v,%v) error(%v)", gid, gid, err)
		return err
	}
	// 修改 进程 gid
	if err := syscall.Setreuid(uid, uid); err != nil {
		logf("syscall.Setreuid(%v,%v) error(%v)", uid, uid, err)
		return err
	}
	return nil
}
