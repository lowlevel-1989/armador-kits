// encvol_file.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>

package encvol


import (
    "bufio"
    "errors"
    "io"
    "crypto/rand"
	"os"
)

const (
    bufSize	= 1048576
)


func (ev *EncVol) Fill() error {

    in := rand.Reader
    out := bufio.NewWriterSize(ev.fh, bufSize)
	
    // XXX: TODO: include a progress bar
	
    nc,e := io.CopyN(out,in, int64(ev.sz))

    if nil != e || uint64(nc) < ev.sz {	
		return errors.New("Could not initialize container")
    }
    return nil
}

func (ev *EncVol) Stamp(cm *ContainerMeta) error {
	
	ev.cm=cm
	return writeContainerMeta(ev.fh, ev.key, cm)
}

func OpenRaw(ff *os.File) (*EncVol,error) {
	
	fi,err := ff.Stat()
	if nil!=err {
		return nil,err
	}
	
	pcm := new(ContainerMeta)
	fsz := uint64(fi.Size())
	
	err = readContainerMeta(ff,nil,pcm)
	if nil!= err {
		return nil,err
	}
	
	return &EncVol{ff,fsz,pcm,nil},nil
}

func Open(ff *os.File, key []byte) (*EncVol,error) {
	
	if nil==ff || VolKeySize!=len(key) {
		return nil,errors.New("Invalid key!")
	}
	
	fi,err := ff.Stat()
	if nil!=err {
		return nil,err
	}
	
	pcm := new(ContainerMeta)
	fsz := uint64(fi.Size())
	
	err = readContainerMeta(ff,key,pcm)
	if nil!= err {
		return nil,err
	}
	
	return &EncVol{ff,fsz,pcm,key},nil
}

func (ev *EncVol) File() *os.File {
	return ev.fh
}
