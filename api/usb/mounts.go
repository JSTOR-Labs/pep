//go:build linux
// +build linux

package usb

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/rs/zerolog/log"
)

var (
	ErrAlreadyMounted  = errors.New("drive already mounted")
	ErrNotMounted      = errors.New("drive isn't mounted")
	ErrMountingFailure = errors.New("unable to change drive mounted state")
)

type (
	Drive struct {
		mountPoint string
		drivePath  string
		mounted    bool
	}
)

func NewDrive(path string, mountPoint string) *Drive {
	return &Drive{
		mountPoint: mountPoint,
		drivePath:  path,
	}
}

func (d *Drive) IsMounted() bool {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), d.drivePath) {
			return true
		}
	}
	return false
}

func (d *Drive) Mount(rdonly bool) error {
	if d.mounted {
		return ErrAlreadyMounted
	}

	flags := syscall.MS_NOSUID | syscall.MS_NODEV
	if rdonly {
		flags |= syscall.MS_RDONLY
	}

	cmd := exec.Command("/bin/mount", d.drivePath, d.mountPoint, "-o", "uid=elasticsearch,gid=elasticsearch,umask=000")
	output, err := cmd.CombinedOutput()
	log.Debug().Err(err).Str("output", string(output)).Msg("mounting drive")
	if err != nil {
		return err
	}

	d.mounted = true
	return nil
}

func (d *Drive) Unmount() error {
	if !d.mounted {
		return ErrNotMounted
	}
	cmd := exec.Command("/bin/umount", d.drivePath)
	err := cmd.Run()
	if err != nil {
		return ErrMountingFailure
	}

	d.mounted = false
	return nil
}

func (d *Drive) MountPath() string {
	return d.mountPoint
}
