package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kai-munekuni/user-api/internal/domain/model"
)

var userKey = "users"

// Firestore db
type Firestore struct {
	c *firestore.Client
}

// NewFirestore intilaize Firestore
func NewFirestore(c *firestore.Client) *Firestore {
	return &Firestore{c: c}
}

// CreateUser createuser
func (s *Firestore) CreateUser(ctx context.Context, ID, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	ref := s.c.Collection(userKey).Doc(ID)
	err = s.c.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if doc.Exists() {
			return model.ErrUserAlreadyExists
		}
		if err != nil && !isNotFoundError(err) {
			return fmt.Errorf("%w", err)
		}
		return tx.Set(ref, firestoreUser{
			Password: hashedPassword,
			Nickname: ID,
			Comment:  "",
		}, firestore.MergeAll)
	})
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

type firestoreUser struct {
	Password string `firestore:"password"`
	Nickname string `firestore:"nickname"`
	Comment  string `firestore:"comment"`
}

func isNotFoundError(err error) bool {
	return status.Code(err) == codes.NotFound
}

// Authorize authorize user
func (s *Firestore) Authorize(ctx context.Context, ID, password string) error {
	dsnap, err := s.c.Collection(userKey).Doc(ID).Get(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	var u firestoreUser
	if err := dsnap.DataTo(&u); err != nil {
		return fmt.Errorf("%w", err)
	}
	if !checkPasswordHash(password, u.Password) {
		return fmt.Errorf("password is incorrect")
	}
	return nil
}

// Find find user by id
func (s *Firestore) Find(ctx context.Context, ID string) (*model.User, error) {
	dsnap, err := s.c.Collection(userKey).Doc(ID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var u firestoreUser
	if err := dsnap.DataTo(&u); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	result := model.User{
		ID:       ID,
		Nickname: u.Nickname,
		Comment:  u.Comment,
	}

	return &result, nil
}

// UpdateField update user field by cas
func (s *Firestore) UpdateField(ctx context.Context, ID string, nickname, comment *string) (*model.User, error) {
	u := model.User{}
	if nickname != nil {
		if *nickname == "" {
			u.Nickname = ID
		} else {
			u.Nickname = *nickname
		}
	}
	if comment != nil {
		u.Comment = *comment
	}
	_, err := s.c.Collection(userKey).Doc(ID).Update(ctx, []firestore.Update{
		{
			Path:  "nickname",
			Value: u.Nickname,
		}, {
			Path:  "comment",
			Value: u.Comment,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	u.ID = ID
	return &u, nil
}

// Delete delete user
func (s *Firestore) Delete(ctx context.Context, ID string) error {
	_, err := s.c.Collection(userKey).Doc(ID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
