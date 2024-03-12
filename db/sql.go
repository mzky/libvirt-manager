package db

import (
	"libvirt-manager/kvm"
)

type VmInfo struct {
	DeviceIp  string
	VmName    string
	VmIp      string
	VmCpu     string
	VmMem     string
	VmDisk    string
	ImageName string
}

func InsertVmInfo(handle *VirtSql, vm kvm.VMConfig) error {
	/*sql := "INSERT INTO vm_info (id，device_ip, vm_name, vm_cpu, vm_mem, vm_disk, image_type, image_name, vm_remarks)" +
	  	" values (NULL，ip, );"
	  //vm.VMName, vm.
	  _, err := handle.Db.Exec(sql)
	  if err != nil {
	  	logrus.Errorf("exec table error, err = %s", err.Error())
	  	return err
	  }*/
	return nil
}
