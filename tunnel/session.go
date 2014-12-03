package tunnel

import (
	"encoding/base64"
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

func (session *Session) Run() {

}
