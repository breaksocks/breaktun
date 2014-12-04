package tunnel

import (
	"sync"
)

type SessionManager struct {
	lock     sync.RWMutex
	sessions map[SessionId]*Session
	reftable map[string]*Session
	exit     chan *Session
}

func NewSessionManager() *SessionManager {
	mgr := &SessionManager{}
	mgr.sessions = make(map[SessionId]*Session)
	mgr.exit = make(chan *Session, 64)
	return mgr
}

func (mgr *SessionManager) NewSession(addr *net.UDPAddr) *Session {
	session := NewSession("", addr)
	session.exit = mgr.exit

	mgr.lock.Lock()
	mgr.reftable[addr] = session
	mgr.lock.Unlock()
	return session, nil
}

func (mgr *SessionManager) SessionFeedId(session *Session, sid SessionId) {
	mgr.lock.Lock()
	mgr.sessions[sid] = session
	mgr.lock.Unlock()
}

func (mgr *SessionManager) GetSessionById(sid SessionId) *Session {
	mgr.lock.RLock()
	session := mgr.sessions[sid]
	mgr.lock.RUnlock()

	return session
}

func (mgr *SessionManager) GetSessionByAddr(addr string) *Session {
	mgr.lock.RLock()
	session := mgr.reftable[addr]
	mgr.lock.RUnlock()

	return session
}

func (mgr *SessionManager) DelSessionById(sid SessionId) {
	mgr.lock.Lock()
	if session, ok := mgr.sessions[sid]; ok {
		delete(mgr.sessions, sid)
		delete(mgr.reftable, session.RemoteAddr)
	}
	mgr.lock.Unlock()
}

func (mgr *SessionManager) DelSession(session *Session) {
	mgr.lock.Lock()
	delete(mgr.sessions, session.Id)
	delete(mgr.reftable, session.RemoteAddr)
	mgr.lock.Unlock()
}

func (mgr *SessionManager) AutoDelSession() {
	for {
		if session, ok := <-mgr.exit; ok {
			mgr.DelSession(session)
		} else {
			break
		}
	}
}
