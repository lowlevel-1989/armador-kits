package loop

/*
#include "loop.h"
*/
import "C"


// Setup loopdev /dev/loop0 -- temporarily
tracer.Trace(2,"Attaching loopdev")

if _,err := os.Stat(loopNode); de.CheckErrno(err,syscall.ENOENT) {

	tracer.Warn("node '"+loopNode+"' not found. About to create it")
	if ce := loop.CreateNode(0); nil!=ce {
		tracer.FatalErr(err)
		os.Exit(4)
	}

}

err = loop.SetupX(0,fout,1*units.Mibi, licData.VolKey)
if /*lfd < 0 ||*/ nil!=err  {
	tracer.Fatal("Could not attach to loopdev")
	tracer.FatalErr(err)
	os.Exit(253)
}

//fmt.Scanln()

tracer.Trace(3,"MKFS ...")

var mkfsopts []string
if forceCreate {
	tracer.Info("Forcing overwrite by MKFS")
	mkfsopts = []string{"-f"}
}
tracer.TraceV(5,"MKFS args:",mkfsopts)

