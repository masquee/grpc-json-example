package server

import (
	"context"
	"sync"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "grpc-json-example/proto"
)

type Backend struct {
	mu    *sync.RWMutex
	users []*pb.User
}

var _ pb.UserServiceServer = (*Backend)(nil)

func New() *Backend {
	return &Backend{
		mu: &sync.RWMutex{},
	}
}

func (b *Backend) AddUser(ctx context.Context, user *pb.User) (*empty.Empty, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Check user ID doesn't already exist
	for _, u := range b.users {
		if u.GetId() == user.GetId() {
			return nil, status.Error(codes.FailedPrecondition, "user already exists")
		}
	}

	if user.GetCreateDate() == nil {
		user.CreateDate = ptypes.TimestampNow()
	}

	b.users = append(b.users, user)

	return new(empty.Empty), nil
}

func (b *Backend) ListUsers(_ *empty.Empty, srv pb.UserService_ListUsersServer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, user := range b.users {
		err := srv.Send(user)
		if err != nil {
			return err
		}
	}

	return nil
}
