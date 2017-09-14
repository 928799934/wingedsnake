package wingedsnake

import (
	"os"
	"path/filepath"
)

func startProcess(env []string, files []*os.File) (*os.Process, error) {
	// Fork exec child process
	name := filepath.Base(os.Args[0]) + " worker process"
	process, err := os.StartProcess(os.Args[0], []string{name}, &os.ProcAttr{Env: env, Files: files})
	if err != nil {
		logf("Fail to fork exec %v", err)
		return nil, err
	}
	return process, nil
}
