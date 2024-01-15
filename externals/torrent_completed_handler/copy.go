package TorrentCompletedHandler

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type TorrentCompletedHandlerCopy struct {
	baseFolder    string
	destinyFolder string
}

func NewTorrentCompletedHandlerCopy(baseFolder, destinyFolder string) (TorrentCompletedHandlerCopy, error) {
	return TorrentCompletedHandlerCopy{
		baseFolder:    baseFolder,
		destinyFolder: destinyFolder,
	}, nil
}

func (t TorrentCompletedHandlerCopy) CreateExec(destinyPath string) (func(fileName []string) (bool, error), error) {
	return func(fileName []string) (bool, error) {
		os.MkdirAll(path.Join(t.destinyFolder, path.Dir(destinyPath)), 0750)
		cmd := exec.Command("cp", "--reflink=auto", path.Join(t.baseFolder, fileName[0]), path.Join(t.destinyFolder, destinyPath))
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			return false, fmt.Errorf("Failed to execute copy with error: %v", err)
		}
		defer stderrPipe.Close()
		if err := cmd.Start(); err != nil {
			return false, fmt.Errorf("Failed to execute copy with error: %v", err)
		}
		reader := bufio.NewReader(stderrPipe)
		line, err := reader.ReadString('\n')
		err = cmd.Wait()
		if err != nil {
			var exitError *exec.ExitError
			if errors.As(err, &exitError) {
				return false, fmt.Errorf("Failed to execute copy with error: %v", line)
			} else {
				return false, fmt.Errorf("Failed to execute copy with error: %v", err)
			}
		}
		return true, nil
	}, nil
}
