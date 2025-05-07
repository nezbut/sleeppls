package sleeppls

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ShutDowner interface {
	ShutDown() error
}

type ComputerShutDowner struct {
	os string
}

func NewComputerShutDowner(Os string) ShutDowner {
	return &ComputerShutDowner{
		os: strings.ToLower(Os),
	}
}

func (s *ComputerShutDowner) ShutDown() error {
	var err error
	switch s.os {
	case "windows":
		err = exec.Command("shutdown", "/s", "/t", "30", "/c", "SleepPls: Computer will shutdown in 30 seconds").Run()
	case "linux":
		err = exec.Command("shutdown", "-h", "now").Run()
	case "darwin":
		if err = exec.Command("osascript", "-e", `tell app "System Events" to shut down`).Run(); err != nil {
			err = exec.Command("shutdown", "-h", "now").Run()
		}
	default:
		err = fmt.Errorf("unsupported OS")
	}
	if err != nil {
		if os.IsPermission(err) {
			err = errors.New("the Permission denied error, please run program with \"sudo\"")
		} else {
			err = fmt.Errorf("error during computer shutdown: %w", err)
		}
	}
	return err
}
