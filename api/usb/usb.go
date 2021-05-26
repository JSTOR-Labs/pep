// +build !linux

package usb

/*
This file provides the functions from usb_linux, they are effectively noop on other platforms
due to no support from the supporting libraries on other platforms
*/

func findUSBDrives() []USBDrive {
	return make([]USBDrive, 0)
}

func formatDrive(_ string) error {
	return nil
}

func initializeDrive(_ string, _ string, _ []string) error {
	return nil
}
