package repo

import (
	"context"
	"time"
)

type IRepo interface {
	Ping() map[string]string
}

// mongo sample di
type Repo struct {
	UserRepo IUserRepo
}

func NewRepo(userRepo IUserRepo) *Repo {
	return &Repo{
		UserRepo: userRepo,
	}
}

func (m *Repo) Ping() map[string]string {
	return map[string]string{"text": "pong"}
}

func NewCTX() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
