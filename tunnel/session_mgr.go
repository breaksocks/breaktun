package tunnel

import (
	"sync"
)

type SessionManager struct {
	lock     sync.RWMutex
	sessions map[SessionId]*Session
	reftable map[string]*Session
}

func NewSessionManager() *SessionManager {
	mgr := &SessionManager{}
	mgr.sessions = make(map[SessionId]*Session)
	return mgr
}

func (mgr *SessionManager) NewSession(addr string, ctx *CipherContext) (*Session, error) {
	session_id, err := ctx.MakeSessionId()
	if err != nil {
		return nil, err
	}

	session := &Session{}
	session.Id = session_id
	session.RemoteAddr = addr

	mgr.lock.Lock()
	mgr.sessions[session_id] = session
	mgr.reftable[addr] = session
	mgr.lock.Unlock()
	return session, nil
}

func (mgr *SessionManager) GetSession(sid SessionId) *Session {
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
