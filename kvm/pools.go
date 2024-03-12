package kvm

import (
	"fmt"
	"libvirt-manager/utils"
)

/*
存储池
virsh pool-list --all
存储卷
virsh vol-list default
*/

// CloneVM 从存储池克隆已存在的虚拟机
func CloneVM(src, dst string) (string, error) {
	return PublicVMControl("vol-clone", utils.JoinVMPath(src), fmt.Sprintf("%s%s", dst, ".qcow2"))
}
