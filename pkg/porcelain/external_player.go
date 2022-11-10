package porcelain

import (
	"os/exec"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

type ExternalPlayer string

func (f *ExternalPlayer) CheckValid() bool {
	valid, _ := fsutil.FileExists(string(*f))
	return valid
}

func (f *ExternalPlayer) RunPathInExternalPlayer(videoPath string) {
	if valid := f.CheckValid(); !valid {
		logger.Warn("Attempted to run External Player when External Player file does not exist.")
		return
	}

	logger.Debugf("External player executed on: %s", videoPath)
	args := []string{videoPath}
	cmd := exec.Command(string(*f), args...)
	_, err := cmd.Output()

	if err != nil {
		logger.Errorf("Eternal Player: \"%s\": %s", videoPath, err.Error())
	}
}
