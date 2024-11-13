package user_application

import (
	"context"
	"errors"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/mik3lon/starter-template/pkg/bus"
)

type ListUsersQuery struct {
	Page   int
	Size   int
	Filter map[string]string
}

func (c ListUsersQuery) Id() string {
	return "list_users_query"
}

type ListUsersQueryHandler struct {
	r user_domain.UserRepository
}

func NewListUsersQueryHandler(r user_domain.UserRepository) *ListUsersQueryHandler {
	return &ListUsersQueryHandler{r: r}
}

func (fuqh ListUsersQueryHandler) Handle(ctx context.Context, query bus.Dto) (interface{}, error) {
	q, ok := query.(*ListUsersQuery)
	if !ok {
		return nil, errors.New("invalid query")
	}

	userList, err := fuqh.r.FindAll(ctx, q.Page, q.Size)
	if err != nil {
		panic(err)
	}

	return userList, nil
}
