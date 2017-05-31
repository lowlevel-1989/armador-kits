package fs


const (

		MS_ACTIVE                        = 0x40000000
        MS_ASYNC                         = 0x1
        MS_BIND                          = 0x1000
        MS_DIRSYNC                       = 0x80
        MS_INVALIDATE                    = 0x2
        MS_I_VERSION                     = 0x800000
        MS_KERNMOUNT                     = 0x400000
        MS_MANDLOCK                      = 0x40
        MS_MGC_MSK                       = 0xffff0000
        MS_MGC_VAL                       = 0xc0ed0000
        MS_MOVE                          = 0x2000
        MS_NOATIME                       = 0x400
        MS_NODEV                         = 0x4
        MS_NODIRATIME                    = 0x800
        MS_NOEXEC                        = 0x8
        MS_NOSUID                        = 0x2
        MS_NOUSER                        = -0x80000000
        MS_POSIXACL                      = 0x10000
        MS_PRIVATE                       = 0x40000
        MS_RDONLY                        = 0x1
        MS_REC                           = 0x4000
        MS_RELATIME                      = 0x200000
        MS_REMOUNT                       = 0x20
        MS_RMT_MASK                      = 0x800051
        MS_SHARED                        = 0x100000
        MS_SILENT                        = 0x8000
        MS_SLAVE                         = 0x80000
        MS_STRICTATIME                   = 0x1000000
        MS_SYNC                          = 0x4
        MS_SYNCHRONOUS                   = 0x10
        MS_UNBINDABLE                    = 0x20000
)

const (
	fsType	= "xfs"
)