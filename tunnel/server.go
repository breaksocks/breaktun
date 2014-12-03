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

func (ser *Server) Close() error {
	return ser.conn.Close()
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
			session.writeToTun <- buf[:n]
		} else {
			session := ser.mgr.NewSession(cliaddr)
			session.writeToTun <- buf[:n]
			go ser.tunWriteBack(addr, session.writeToClient)
			go session.Run()
		}
	}
}

func (ser *Server) tunWriteBack(addr *net.UDPAddr, write chan []byte) {
	for {
		if data, ok := <-write; ok {
			if _, err := ser.conn.WriteToUDP(data, addr); err != nil {
				glog.Errorf("conn write fail: %v", err)
				break
			}
		} else {
			break
		}
	}
}
