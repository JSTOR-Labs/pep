package main

import (
	"os"
	"syscall"

	"github.com/rs/zerolog/log"
)

const (
	PersistentVolumeDevPath    = "/dev/sdb1"
	BootVolumeDevPath          = "/dev/sda1"
	PersistentVolumeMountPoint = "/mnt/persist"
	BootVolumeMountPoint       = "/mnt/boot"
	UpdateStagingPath          = PersistentVolumeMountPoint + "/update"
	RootFSUpdatePath           = UpdateStagingPath + "/rootfs.squashfs"
	KernelUpdatePath           = UpdateStagingPath + "/bzImage"
	KernelFullPath             = BootVolumeMountPoint + "/bzImage"
)

func checkUpdate() error {
	if err := MkDirsAll(PersistentVolumeMountPoint); err != nil {
		return err
	}
	if err := syscall.Mount(PersistentVolumeDevPath, PersistentVolumeMountPoint, "ext4", 0, ""); err != nil {
		return err
	}
	defer func() {
		if err := syscall.Unmount(PersistentVolumeMountPoint, 0); err != nil {
			log.Error().Err(err).Msg("unmount persist")
		}
	}()

	if _, err := os.Stat(RootFSUpdatePath); err == nil {
		log.Info().Msg("root filesystem update found, installing")
		if err := os.Rename(RootFSUpdatePath, RootFSFullPath); err != nil {
			return err
		}
	}

	if _, err := os.Stat(KernelUpdatePath); err == nil {
		log.Info().Msg("kernel update found, installing")
		if err := syscall.Mount(BootVolumeDevPath, BootVolumeMountPoint, "vfat", syscall.MS_NOATIME, "iocharset=utf8"); err != nil {
			return err
		}

		defer func() {
			if err := syscall.Unmount(BootVolumeMountPoint, 0); err != nil {
				log.Error().Err(err).Msg("unmount boot")
			}
		}()

		if err := os.Rename(KernelUpdatePath, KernelFullPath); err != nil {
			return err
		}
	}

	return nil
}
