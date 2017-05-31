// mkfs_xfs.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>

package fs

import (
	"os/exec"
)


const (
	mkfs_bin = "/sbin/mkfs.xfs"
)


func MakeFS(devname string, opts []string) error {
	
	var args []string = []string{"mkfs"}
	if len(opts) > 0 {
		args = append(args,opts...)
	}
	args = append(args,devname)
	
	cmd := exec.Cmd{Path: mkfs_bin, Args: args, Env: []string{}}
	
	// XXX: TODO: ExecErrorConv
	return cmd.Run()
}
