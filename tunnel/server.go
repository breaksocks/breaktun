package tunnel

import (
	"github.com/golang/glog"
	"net"
)

type Server struct {
	cfg  *ServerConfig
	conn *net.UDPConn
	mgr  *SessionManager
}

func NewServer(cfg *ServerConfig) (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", cfg.ListenAddr)
	if err != nil {
		return nil, err
	}

	if conn, err := net.ListenUDP("udp", addr); err != nil {
		server := new(Server)
		server.cfg = cfg
		server.conn = conn
		server.mgr = NewSessionManager()
		return server, nil
	} else {
		return nil, err
	}
}

func (ser *Server) Run() {
	for {
		buf := make([]byte, 2048)
		n, cliaddr, err := ser.conn.ReadFromUDP(buf)
		if err != nil {
			glog.Fatalf("server read fail: %v", err)
		}

		addr := cliaddr.String()
		if session := ser.mgr.GetSessionByAddr(addr); session != nil {
			session.write <- buf[:n]
		} else {
			session := ser.mgr.NewSession(addr)
			session.write <- buf[:n]
		}
	}
}
