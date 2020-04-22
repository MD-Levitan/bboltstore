package store

import (
	"bytes"
	"encoding/base32"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"go.etcd.io/bbolt"
)

// Store represents a session store.
type Store struct {
	codecs []securecookie.Codec
	config Config
	db     *bbolt.DB
}

// Get returns a session for the given name after adding it to the registry.
//
// See gorilla/sessions FilesystemStore.Get().
func (s *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See gorilla/sessions FilesystemStore.New().
func (s *Store) New(r *http.Request, name string) (*sessions.Session, error) {
	var err error
	session := sessions.NewSession(s, name)
	options := s.config.SessionOptions
	session.Options = &options
	session.IsNew = true
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.codecs...)
		if err == nil {
			ok, err := s.load(session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

// Save adds a single session to the response.
func (s *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.Options.MaxAge < 0 {
		s.delete(session)
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		// Build an alphanumeric ID.
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		if err := s.save(session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

// load loads a session data from the database.
// True is returned if there is a session data in the database.
func (s *Store) load(session *sessions.Session) (bool, error) {
	// exists represents whether a session data exists or not.
	var exists bool
	err := s.db.View(func(tx *bbolt.Tx) error {
		id := []byte(session.ID)
		bucket := tx.Bucket(s.config.DBOptions.BucketName)
		// Get the session data.
		data := bucket.Get(id)
		if data == nil {
			return nil
		}
		sessionData, err := GetSession(data)
		if err != nil {
			return err
		}
		// Check the expiration of the session data.
		if Expired(sessionData) {
			err := s.db.Update(func(txu *bbolt.Tx) error {
				return txu.Bucket(s.config.DBOptions.BucketName).Delete(id)
			})
			return err
		}
		exists = true
		dec := gob.NewDecoder(bytes.NewBuffer(sessionData.Values))
		return dec.Decode(&session.Values)
	})
	return exists, err
}

// delete removes the key-value from the database.
func (s *Store) delete(session *sessions.Session) error {
	err := s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(s.config.DBOptions.BucketName).Delete([]byte(session.ID))
	})
	if err != nil {
		return err
	}
	return nil
}

// save stores the session data in the database.
func (s *Store) save(session *sessions.Session) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(session.Values)
	if err != nil {
		return err
	}
	data, err := json.Marshal(NewSession(buf.Bytes(), session.Options.MaxAge))
	if err != nil {
		return err
	}
	err = s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(s.config.DBOptions.BucketName).Put([]byte(session.ID), data)
	})
	return err
}

// Creates new sesion storage using opened db.
func NewStore(db *bbolt.DB, config Config, keyPairs ...[]byte) (*Store, error) {
	config.setDefault()
	store := &Store{
		codecs: securecookie.CodecsFromPairs(keyPairs...),
		config: config,
		db:     db,
	}
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(config.DBOptions.BucketName)
		return err
	})
	if err != nil {
		return nil, err
	}
	return store, nil
}

// Opens a db and creates new sesion storage.
func NewStoreWithDB(dbPath string, config Config, keyPairs ...[]byte) (*Store, error) {
	config.setDefault()

	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return nil, err
	}

	store := &Store{
		codecs: securecookie.CodecsFromPairs(keyPairs...),
		config: config,
		db:     db,
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(config.DBOptions.BucketName)
		return err
	})
	if err != nil {
		return nil, err
	}
	return store, nil
}
