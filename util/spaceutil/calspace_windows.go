// +build windows

package spaceutil

import "github.com/StackExchange/wmi"

func GetDiskInfo() (diskSize, diskFree uint64) {
	var storageinfo []storageInfo
	err := wmi.Query("Select * from Win32_LogicalDisk", &storageinfo)
	if err != nil {
		return
	}
	for _, storage := range storageinfo {
		diskSize += storage.Size
		diskFree += storage.FreeSpace
	}
	return
}
