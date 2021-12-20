//go:build !linux
// +build !linux

package usb

/*
This file provides the functions from usb_linux, they are effectively noop on other platforms
due to no support from the supporting libraries on other platforms
*/

func FindUSBDrives() []USBDrive {
	return make([]USBDrive, 0)
}

func FormatDrive(_ string) error {
	return nil
}

func BuildFlashDrive(_ string, _ string, _ []string) error {
	return nil
}
