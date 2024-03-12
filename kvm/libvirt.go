package kvm

import (
	"fmt"
	"github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

/*
虽然这个库现在大部分功能都正常，但尽量不要使用，qemu和virsh命令会随着系统升级随时更新，但这个库很久没有更新了，
除非命令解决不了的，比如修改虚拟机最大内存等，非必要不使用此库
*/

// ConnectLibvirt 连接libvirt库
func ConnectLibvirt() (*libvirt.Libvirt, error) {
	c, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", 10*time.Second)
	if err != nil {
		logrus.Errorf("failed to dial libvirt: %v", err)
		return nil, err
	}
	defer c.Close()

	l := libvirt.New(c)
	if e := l.Connect(); e != nil {
		logrus.Errorf("failed to connect: %v", e)
		return nil, e
	}

	return l, nil
}

// VirtDomainList 虚拟机Domain列表
func VirtDomainList(l *libvirt.Libvirt) ([]libvirt.Domain, error) {
	domains, err := l.Domains()
	if err != nil {
		logrus.Errorf("failed to retrieve domains: %v", err)
		return nil, err
	}
	return domains, nil
}

// FindDomain 查找和返回指定虚拟机的Domain
func FindDomain(l *libvirt.Libvirt, name string) (libvirt.Domain, error) {
	domains, err := VirtDomainList(l)
	if err != nil {
		logrus.Errorf("failed to VirtDomainList: %v", err)
		return libvirt.Domain{}, err
	}
	for _, d := range domains {
		if d.Name == name {
			logrus.Printf("%d\t%s\t%x\n", d.ID, d.Name, d.UUID)
			return d, nil
		}
	}
	return libvirt.Domain{}, fmt.Errorf("cannot be found Domains")
}

// Disconnect 关闭libvirt连接
func Disconnect(l *libvirt.Libvirt) {
	if e := l.Disconnect(); e != nil {
		logrus.Errorf("failed to Disconnect: %v", e)
	}
}

// SetVMMaxMemory 修改虚拟机内存最大值
func SetVMMaxMemory(name string, max uint64) error {
	l, err := ConnectLibvirt()
	if err != nil {
		logrus.Errorf("failed to ConnectLibvirt: %v", err)
		return err
	}

	domain, err := FindDomain(l, name)
	if err != nil {
		logrus.Errorf("failed to FindDomain: %v", err)
		return err
	}
	defer Disconnect(l)

	return l.DomainSetMaxMemory(domain, max)
}

// GetDomainXML 获取指定虚拟机的xml配置（暂时没什么用，可通过修改xml内容来修改虚拟机各种配置）
func GetDomainXML(name string) (string, error) {
	l, err := ConnectLibvirt()
	if err != nil {
		logrus.Errorf("failed to ConnectLibvirt: %v", err)
		return "", err
	}
	domain, err := FindDomain(l, name)
	if err != nil {
		logrus.Errorf("failed to FindDomain: %v", err)
		return "", err
	}
	defer Disconnect(l)

	return l.DomainGetXMLDesc(domain, 0)
}

// SetInterfaceParameters 用来修改xml，非必要不使用
func SetInterfaceParameters(name string) error {
	l, err := ConnectLibvirt()
	if err != nil {
		logrus.Errorf("failed to ConnectLibvirt: %v", err)
		return err
	}

	defer Disconnect(l)
	return nil
}
