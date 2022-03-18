package video

import (
	"path/filepath"
	"strings"
)

// GetFunscriptPath returns the path of a file
// with the extension changed to .funscript
func GetFunscriptPath(path string, funscriptDir string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(base, ext)
	dir := filepath.Dir(funscriptDir)
	return dir + "/" + fn + ".funscript"
}
