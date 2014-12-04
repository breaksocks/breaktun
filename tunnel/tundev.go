package tunnel

import (
	"net"
)

type TunDev interface {
	SetLocalAddr(ip net.IP) error
	SetRemoteAddr(ip net.IP) error
	SetMTU(mtu int) error

	Read(bs []byte) (int, error)
	Write(bs []byte) (int, error)
	Close() error
}
