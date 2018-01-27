// encvol.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package encvol

import (
	"fmt"
	"os"
	"time"
)

const (
	EncvolOffset	= 1048576
	VolKeySize	= 256/8
)

type EncVol struct {
	fh		*os.File
	sz		uint64
	cm		*ContainerMeta
	key		[]byte;
}


type ProgressCallback func(uint64,uint64)


func New(ff *os.File, size uint64, tag string) *EncVol {
	cm := new(ContainerMeta)
	cm.TagName = tag
	return &EncVol{ff,size,cm,nil}
}

func NewX(ff *os.File, size uint64, tag string, key []byte) *EncVol {
	cm := new(ContainerMeta)
	cm.TagName = tag
	return &EncVol{ff,size,cm,key}
}


func (ev *EncVol) Release() {
	ev.cm = nil
	ev.key = nil
}

func (ev *EncVol) Close() {

	if nil != ev.fh {
		ev.fh.Close()
	}
	ev.key = nil
	ev.cm = nil
}


func (ev *EncVol) Touch(ts time.Time) error {
	ev.cm.LastUsed = ts
	return writeContainerMeta(ev.fh, ev.key, ev.cm)
}

func (ev *EncVol) Commit() {
	// Just write metaInfo to the file
	writeContainerMeta(ev.fh, ev.key, ev.cm)
}


func (ev *EncVol) Validate() bool {


	return true
}


// Implement "Stringer"
func (ev *EncVol) String() string {

	if nil== ev.cm {
		return "ENCVOL not ready"
	}

	cm := ev.cm
	return fmt.Sprintf(`{ "tag": "%s", "size": "%d", "created": "%v", "lastUsed": "%v", "physSize": "%d"}`,
	cm.TagName, cm.Size, cm.Created, cm.LastUsed, ev.sz)
}

func (ev *EncVol) Metadata() *ContainerMeta {
	return ev.cm
}

