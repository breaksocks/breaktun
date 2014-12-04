package tunnel

import (
	"encoding/base64"
	"github.com/golang/glog"
	"net"
)

type SessionId string

func SessionIdFromBytes(bs []byte) SessionId {
	return SessionId(base64.StdEncoding.EncodeToString(bs))
}

func (sid SessionId) Bytes() ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(sid))
}

func (sid SessionId) size() int {
	return len(sid)
}

type Session struct {
	Id         SessionId
	RemoteAddr *net.UDPAddr
	Username   string
	cipherExg  *CipherExchange
	cipherCfg  *CipherConfig

	exit chan *Session

	dev           TunDev
	writeToTun    chan []byte
	writeToClient chan []byte
}

func NewSession(sid SessionId, addr *net.UDPAddr) *Session {
	session := new(Session)
	session.Id = sid
	session.RemoteAddr = addr
	session.writeToTun = make([]byte, 1024)
	session.writeToClient = make([]byte, 1024)
	return session
}

func (session *Session) Close() {
	session.exit <- session
}

func (session *Session) Run() {
	defer session.Close()

	if !session.login() || !session.setupDev() {
		return
	}

	normal_close_dev := false

	go func() {
		for {
			buf := make([]byte, 2048)
			if n, err := serssion.dev.Read(buf); err == nil {
				// TODO: encrypt
				session.writeToClient <- buf[:n]
			} else if !normal_close_dev {
				glog.Errorf("read tun fail: %v", err)
				break
			} else {
				break
			}
		}
	}()

	for {
		if data, ok := session.writeToTun; ok {
			// TODO: decrypt
			if _, err := session.dev.Write(data); err != nil {
				glog.Errorf("write to tun fail: %v", err)
				break
			}
		} else {
			normal_close_dev = true
			if err := session.dev.Close(); err != nil {
				glog.Errorf("close tun fail: %v", err)
			}
			break
		}
	}
}

func (session *Session) login() bool {
	reqer := NewRetryRequester(5, 200, session.writeToTun, session.writeToClient)
	reqer.RetryGet(3)
}

func (session *Session) setupDev() bool {
}
