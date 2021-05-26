package utils

// #include <unistd.h>
// #include <sys/reboot.h>
import "C"

// Reboot reboots the system
func Reboot() {
	C.sync()
	C.setuid(0)
	C.reboot(C.RB_AUTOBOOT)
}
