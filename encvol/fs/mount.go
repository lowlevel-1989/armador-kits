// mount.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
// REF: http://man7.org/linux/man-pages/man2/mount.2.html

package fs

import (
	"syscall"
)


func MountVol(devname string, mntpoint string) error {

	var flags uintptr = /*syscall.MS_LAZYTIME|*/syscall.MS_NOATIME|syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV
	var err error

	filesystem := "xfs"

	//	fmt.Println("*** MOUNT '",devname,"' at '",mntpoint,"'; flags=",flags);
	err = syscall.Mount(devname, mntpoint, filesystem, flags, "")
	if nil!=err {
		filesystem = "ext4"
		err = syscall.Mount(devname, mntpoint, filesystem, flags, "")
	}
	if nil!=err {
		return err
	}

	//	fmt.Println("*** mounted. About to remount")

	// Make this mount unbindable.  This is like a private mount,
        // and in addition this mount can't be bind mounted.
	var xfl uintptr = syscall.MS_REMOUNT|syscall.MS_UNBINDABLE /*|syscall.MS_SILENT*/
	err = syscall.Mount(devname,mntpoint, filesystem, xfl, "")

	//  	fmt.Println("*** Ok.")

	return err
}

func UmountVol(mntpoint string) error {

	var err error = nil
	var flags int = 0		// syscall.UMOUNT_NOFOLLOW	// requires Golang 1.10

	flags = syscall.MNT_DETACH
	err = syscall.Unmount(mntpoint, flags)
	if nil!=err {
		return err
	}

	return nil
}
