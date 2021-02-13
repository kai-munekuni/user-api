package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/kai-munekuni/user-api/internal/domain/model"
)

// MemoryDB in memory database
// deprecated, should use firestore
type MemoryDB struct {
	db   map[string]*model.User
	lock sync.RWMutex
}

// NewMemoryDB initialize MemoryDB
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{db: map[string]*model.User{}}
}

// CreateUser ceate new user, if user exists, return error
func (m *MemoryDB) CreateUser(ctx context.Context, id, password string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.db[id] != nil {
		return model.ErrUserAlreadyExists
	}
	hashedPassword, err := hashPassword(password)

	if err != nil {
		return fmt.Errorf("error! : %w", err)
	}
	m.db[id] = &model.User{
		ID:       id,
		Password: hashedPassword,
		Nickname: id,
	}

	return nil
}

// Authorize authorize user
func (m *MemoryDB) Authorize(ctx context.Context, id, password string) error {
	m.lock.RLock()
	defer m.lock.RUnlock()
	u := m.db[id]
	if u == nil {
		return fmt.Errorf("cannot find user")
	}
	if !checkPasswordHash(password, u.Password) {
		return fmt.Errorf("password is incorrect")
	}

	return nil
}

// Find find user by id
func (m *MemoryDB) Find(ctx context.Context, ID string) (*model.User, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	result := m.db[ID]
	if result == nil {
		return nil, fmt.Errorf("cannot find user")
	}

	return result, nil
}

// UpdateField update user field by cas
func (m *MemoryDB) UpdateField(ctx context.Context, ID string, nickname, comment *string) (*model.User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	user := m.db[ID]
	if user == nil {
		return nil, fmt.Errorf("cannot find user")
	}
	if nickname != nil {
		if *nickname == "" {
			user.Nickname = ID
		} else {
			user.Nickname = *nickname
		}
	}
	if comment != nil {
		user.Comment = *comment
	}
	m.db[ID] = user
	return user, nil
}

// Delete delete user
func (m *MemoryDB) Delete(ctx context.Context, ID string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.db, ID)
	return nil
}
