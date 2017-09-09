package wingedsnake

func getConfigUser(conf *config) (int, int, error) {
	return 0, 0, errWindows
}

// exchangeOwner windows 不支持修改进程 uid 与 gid
func exchangeOwner(uid, gid int) error {
	return errWindows
}
