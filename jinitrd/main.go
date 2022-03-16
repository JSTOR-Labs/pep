package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/freddierice/go-losetup.v1"
)

/*
jinitrd respects the following kernel command line parameters:
init_root=<path> // The device containing the squashed root filesystem
ro_root=<path> // The path on the host filesystem to the root filesystem
*/

// This is a very static environment, lets just hardcode some values
const (
	RootFSPath     = "/rootfs.squashfs"
	InitialFSMount = "/mnt/rw"
	RootFSFullPath = InitialFSMount + RootFSPath
	RootFSMount    = "/mnt/rootfs"
	OverlayFSPath  = "/mnt/overlay"
	OverlayLower   = "/mnt/lower"
	OverlayUpper   = OverlayFSPath + "/upper"
	OverlayWork    = OverlayFSPath + "/work"
)

var cmdline = make(map[string]string)

type FailureHook struct{}

func (h FailureHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {

}

func Fail() {
	log.Error().Msg("rebooting in 10 seconds")
	time.Sleep(10 * time.Second)
	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

func parseCmdline() error {
	b, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		return err
	}
	parts := strings.Split(strings.TrimSpace(string(b)), " ")
	for _, part := range parts {
		if idx := strings.IndexByte(part, '='); idx > -1 {
			cmdline[part[:idx]] = part[idx+1:]
		} else {
			cmdline[part] = ""
		}
	}
	return nil
}

func init() {
	if err := os.MkdirAll("/dev", 0o777); err != nil {
		log.Error().Err(err).Msg("failed to create /dev")
		Fail()
	}

	// mount devtmpfs
	if err := mount("devtmpfs", "/dev", "devtmpfs", 0); err != nil {
		log.Error().Err(err).Msg("failed to mount devtmpfs")
		Fail()
	}

	// Set up /dev/console
	f, err := os.OpenFile("/dev/console", os.O_RDWR, 0)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open /dev/console")
		Fail()
	}
	// Redirect stdin, stdout, stderr to /dev/console
	os.Stdout = f
	os.Stderr = f
	os.Stdin = f

	output := zerolog.ConsoleWriter{Out: f}

	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}

func main() {
	log.Info().Msg("jinitrd starting")

	log.Info().Msg("mounting filesystems")

	if err := MkDirsAll(InitialFSMount, RootFSMount, OverlayFSPath, OverlayLower); err != nil {
		log.Error().Err(err).Msg("failed to create required directories")
		Fail()
	}

	// mount sysfs
	if err := mount("sysfs", "/sys", "sysfs", 0); err != nil {
		log.Error().Err(err).Msg("failed to mount sysfs")
		Fail()
	}

	// mount proc
	if err := mount("proc", "/proc", "proc", 0); err != nil {
		log.Error().Err(err).Msg("failed to mount proc")
		Fail()
	}

	log.Info().Msg("parsing kernel command line")
	if err := parseCmdline(); err != nil {
		log.Error().Err(err).Msg("failed to parse kernel command line")
		Fail()
	}

	// Mount the main filesystem
	initRoot, ok := cmdline["init_root"]
	mountDev := ""
	if ok {
		parts := strings.Split(initRoot, "=")
		switch len(parts) {
		case 1:
			// This is a device name
			mountDev = parts[0]
		case 2:
			var err error
			// This is a special identifier
			switch parts[0] {
			case "UUID":
				mountDev, err = filepath.EvalSymlinks(filepath.Join("/dev/disk/by-uuid", parts[1]))
			case "PARTUUID":
				mountDev, err = filepath.EvalSymlinks(filepath.Join("/dev/disk/by-partuuid", parts[1]))
			case "LABEL":
				mountDev, err = filepath.EvalSymlinks(filepath.Join("/dev/disk/by-label", parts[1]))
			default:
				log.Error().Msgf("unknown special identifier %q", parts[0])
				Fail()
			}
			if err != nil {
				log.Error().Err(err).Msgf("failed to resolve %q", parts[0])
				Fail()
			}
		}
	} else {
		mountDev = InitialFSPath
	}

	log.Info().Str("dev", mountDev).Msg("mounting initial filesystem")
	if err := mount(mountDev, InitialFSMount, "ext4", 0); err != nil {
		log.Error().Err(err).Msg("failed to mount initial filesystem")
		Fail()
	}

	if err := checkUpdate(); err != nil {
		log.Error().Err(err).Msg("failed to check for updates")
	}

	// Get the root filesystem location
	rootPath, ok := cmdline["ro_root"]
	if !ok {
		rootPath = RootFSPath
	}
	rootFullPath := filepath.Join(InitialFSMount, rootPath)
	log.Info().Str("path", rootFullPath).Msg("looping root filesystem")
	// Create new loop device
	dev, err := losetup.Attach(rootFullPath, 0, true)
	if err != nil {
		log.Error().Err(err).Msg("failed to create loop device")
		Fail()
	}

	log.Info().Str("dev", dev.Path()).Msg("mounting root filesystem")
	// Mount the root filesystem
	if err := rootFS(dev.Path()); err != nil {
		log.Error().Err(err).Msg("failed to mount root filesystem")
		Fail()
	}

	if err := switchRoot(); err != nil {
		log.Error().Err(err).Msg("failed to switch root")
		Fail()
	}
	time.Sleep(time.Second * 2)
}

func rootFS(dev string) error {
	if err := syscall.Mount("tmpfs", OverlayFSPath, "tmpfs", syscall.MS_NOSUID|syscall.MS_NODEV, "size=1G"); err != nil {
		return err
	}

	if err := MkDirsAll(OverlayUpper, OverlayWork); err != nil {
		return err
	}

	if err := syscall.Mount(dev, OverlayLower, "squashfs", syscall.MS_RDONLY, ""); err != nil {
		return err
	}

	if err := syscall.Mount("overlay", RootFSMount, "overlay", 0, "lowerdir="+OverlayLower+",upperdir="+OverlayUpper+",workdir="+OverlayWork); err != nil {
		return err
	}

	return nil
}

func switchRoot() error {
	// mount devtmpfs on new root /dev
	if err := mount("devtmpfs", fmt.Sprintf("%s/dev", RootFSMount), "devtmpfs", 0); err != nil {
		log.Error().Err(err).Msg("failed to mount devtmpfs")
		Fail()
	}

	// I think busybox will mount dev, proc, sys, and tmp filesystems
	if err := os.Chdir(RootFSMount); err != nil {
		log.Error().Err(err).Msgf("failed to chdir to %s", RootFSMount)
		return err
	}

	log.Info().Msg("mounting real root")
	if err := syscall.Mount(".", "/", "", syscall.MS_MOVE, ""); err != nil {
		log.Error().Err(err).Msg("failed to move mount to real root")
		return err
	}

	log.Info().Msg("chroot into actual root")
	if err := syscall.Chroot("."); err != nil {
		log.Error().Err(err).Msg("failed to chroot")
		return err
	}

	log.Info().Msg("chaining to real init")
	if err := syscall.Exec("/sbin/init", []string{"/init"}, os.Environ()); err != nil {
		log.Error().Err(err).Msg("failed to exec init")
		return err
	}

	return nil
}
