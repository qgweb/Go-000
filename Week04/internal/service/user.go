package service

import (
	"Go-000/Week04/internal/conf"
	"Go-000/Week04/internal/dao"
	"Go-000/Week04/internal/model"
	"Go-000/Week04/internal/xerror"
	"context"

	"github.com/pkg/errors"
)

type Service struct {
	dao *dao.Dao
}

func New(config conf.Config) *Service {
	return &Service{
		dao: dao.New(config),
	}
}

func (s *Service) SimpleLogin(ctx context.Context, name, pwd string) (user model.User, err error) {
	if name == "" || pwd == "" {
		return user, errors.WithStack(xerror.NotFound("用户名或密码为空", ""))
	}
	user, err = s.dao.GetUser(ctx, name, pwd)
	if err != nil {
		if errors.Is(err, dao.ErrNotFound) {
			return user, errors.WithStack(
				xerror.Unauthorized("用户名或密码错误", "用户名:%s,密码:%s", name, pwd))
		}
		return user, errors.WithStack(xerror.Internal("系统错误", "%w", err))
	}

	return user, nil
}

func (s *Service) Close(ctx context.Context) {
	s.dao.Close()
}
