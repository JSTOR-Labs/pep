//go:build (linux && ignore) || cgo
// +build linux,ignore cgo

package usb

/*
 #cgo LDFLAGS: -lblkid
#include "blkid/blkid.h"
#include "string.h"
#include "stdlib.h"
*/
import "C"
import "unsafe"

const (
	UUIDIdent     = "UUID"
	TypeIdent     = "TYPE"
	LabelIdent    = "LABEL"
	PartUUIDIdent = "PARTUUID"
)

type BlkIDData struct {
	UUID     string
	Type     string
	Label    string
	PartUUID string
}

func BlkID(path string) (*BlkIDData, error) {
	var (
		blkidUUID     *C.char
		blkidType     *C.char
		blkidLabel    *C.char
		blkidPartUUID *C.char
		device        *C.char
	)

	blkidUUID = C.CString(UUIDIdent)
	blkidType = C.CString(TypeIdent)
	blkidLabel = C.CString(LabelIdent)
	blkidPartUUID = C.CString(PartUUIDIdent)

	defer C.free(unsafe.Pointer(blkidUUID))
	defer C.free(unsafe.Pointer(blkidType))
	defer C.free(unsafe.Pointer(blkidLabel))
	defer C.free(unsafe.Pointer(blkidPartUUID))

	device = C.CString(path)
	defer C.free(unsafe.Pointer(device))

	var (
		uuid, fsType, label, partUUID *C.char
	)
	uuid = C.blkid_get_tag_value(nil, blkidUUID, device)
	fsType = C.blkid_get_tag_value(nil, blkidType, device)
	label = C.blkid_get_tag_value(nil, blkidLabel, device)
	partUUID = C.blkid_get_tag_value(nil, blkidPartUUID, device)

	defer C.free(unsafe.Pointer(uuid))
	defer C.free(unsafe.Pointer(fsType))
	defer C.free(unsafe.Pointer(label))
	defer C.free(unsafe.Pointer(partUUID))

	return &BlkIDData{
		UUID:     C.GoString(uuid),
		Type:     C.GoString(fsType),
		Label:    C.GoString(label),
		PartUUID: C.GoString(partUUID),
	}, nil
}

// func BlkID(path string) (*BlkIDData, error) {
// 	cmd := exec.Command(which.LookupExecutable("blkid"), path)
// 	out, err := cmd.Output()
// 	if err != nil {
// 		log.Warn().Err(err).Str("out", string(out)).Msg("blkid failed")
// 		return nil, err
// 	}
// 	var b BlkIDData
// 	rawData := make(map[string]string)
// 	fields := strings.Split(string(out), " ")
// 	for _, field := range fields[0:] {
// 		values := strings.Split(field, "=")
// 		if len(values) < 2 {
// 			continue
// 		}
// 		rawData[values[0]] = strings.Trim(values[1], "\"")
// 	}
// 	setValue(rawData, "UUID", &b.UUID)
// 	setValue(rawData, "TYPE", &b.Type)
// 	setValue(rawData, "LABEL", &b.Label)
// 	setValue(rawData, "PARTUUID", &b.PartUUID)
// 	return &b, nil
// }
