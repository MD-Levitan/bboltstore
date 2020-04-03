package store

import (
	"encoding/json"
	"time"
)

type Session struct {
	Values           []byte `json:"Values"`
	ExpiresAt        int64  `json:"ExpiresAt"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Session) Reset()         { *m = Session{} }
func (m *Session) String() string { r, _ := json.Marshal(m); return string(r) }
func (*Session) ProtoMessage()    {}

func (m *Session) GetValues() []byte {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *Session) GetExpiresAt() int64 {
	if m != nil {
		return m.ExpiresAt
	}
	return 0
}

// Session converts the byte slice to the session struct value.
func GetSession(data []byte) (Session, error) {
	session := Session{}
	err := json.Unmarshal(data, &session)
	return session, err
}

// Expired checks if the session is expired.
func Expired(session Session) bool {
	return session.ExpiresAt > 0 && session.ExpiresAt <= time.Now().Unix()
}

// NewSession creates and returns a session data.
func NewSession(values []byte, maxAge int) *Session {
	expiresAt := time.Now().Unix() + int64(maxAge)
	return &Session{Values: values, ExpiresAt: expiresAt}
}
