package user_application

import (
	"context"
	"errors"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"github.com/mik3lon/starter-template/pkg/bus"
)

type FindUserQuery struct {
	Email string
}

func (c FindUserQuery) Id() string {
	return "find-user-query-handler"
}

type FindUserQueryHandler struct {
	r user_domain.UserRepository
}

func NewFindUserQueryHandler(r user_domain.UserRepository) *FindUserQueryHandler {
	return &FindUserQueryHandler{r: r}
}

func (fuqh FindUserQueryHandler) Handle(ctx context.Context, query bus.Dto) (interface{}, error) {
	q, ok := query.(*FindUserQuery)
	if !ok {
		return nil, errors.New("invalid query")
	}

	return fuqh.r.FindByEmail(ctx, q.Email)
}
