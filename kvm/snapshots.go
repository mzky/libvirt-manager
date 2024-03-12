package kvm

import (
	"libvirt-manager/utils"
	"path/filepath"
	"strings"
)

type SnapList struct {
	ID   string
	Name string
	DATE string
	ICT  string
}

// CreateSnapshot 创建虚拟机快照
func CreateSnapshot(name string, snapshotName string) (string, error) {
	return utils.Exec("qemu-img", "snapshot", "-c", snapshotName, filepath.Join(utils.VMPath, name+".qcow2"))
}

// RevertSnapshot 回滚快照
func RevertSnapshot(name string, snapshotName string) (string, error) {
	return utils.Exec("qemu-img", "snapshot", "-a", snapshotName, filepath.Join(utils.VMPath, name+".qcow2"))
}

// DeleteSnapshot 删除虚拟机快照
func DeleteSnapshot(name string, snapshotName string) (string, error) {
	return utils.Exec("qemu-img", "snapshot", "-d", snapshotName, filepath.Join(utils.VMPath, name+".qcow2"))
}

// SnapshotList 列出虚拟机快照
func SnapshotList(name string) ([]SnapList, error) {
	var snapList []SnapList
	s, err := utils.Exec("qemu-img", "snapshot", "-l", filepath.Join(utils.VMPath, name+".qcow2"))
	if err != nil {
		return snapList, err
	}
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for _, line := range lines[2:] {
		var snList SnapList
		lin := strings.Fields(line)
		snList.ID = lin[0]
		snList.Name = lin[1]
		snList.DATE = lin[4] + " " + lin[5]
		snList.ICT = lin[7]
		snapList = append(snapList, snList)
	}
	return snapList, err
}
