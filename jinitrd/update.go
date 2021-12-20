package main

import (
	"io"
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

func moveFile(src, dst string) error {
	inputFile, err := os.Open(src)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(dst)
	if err != nil {
		inputFile.Close()
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return err
	}

	return os.Remove(src)
}

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
		if err := moveFile(RootFSUpdatePath, RootFSFullPath); err != nil {
			return err
		}
	}

	if _, err := os.Stat(KernelUpdatePath); err == nil {
		log.Info().Msg("kernel update found, installing")

		if err := MkDirsAll(BootVolumeMountPoint); err != nil {
			return err
		}

		if err := syscall.Mount(BootVolumeDevPath, BootVolumeMountPoint, "vfat", syscall.MS_NOATIME, "iocharset=utf8"); err != nil {
			return err
		}

		defer func() {
			if err := syscall.Unmount(BootVolumeMountPoint, 0); err != nil {
				log.Error().Err(err).Msg("unmount boot")
			}
		}()

		if err := moveFile(KernelUpdatePath, KernelFullPath); err != nil {
			return err
		}
	}

	return nil
}
