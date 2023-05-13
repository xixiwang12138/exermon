package _examples

import (
	"github.com/xixiwang12138/exermon/db"
	"github.com/xixiwang12138/exermon/db/op"
	"github.com/xixiwang12138/xlog"
)

type User struct {
	Id   string
	Name string

	Age int
}

type Grade struct {
	Id     string
	UserId string
	Score  int
}

type UserService struct {
	*db.BaseDao[User]
}

func NewUserService() *UserService {
	// TODO
	return &UserService{BaseDao: db.NewBaseDao[User](nil)}
}

// UpdateUser 不涉及事务操作
func (u *UserService) UpdateUser(xl *xlog.XLogger, id string) (user *User, err error) {
	// Tip1: 使用Instance方法注入context
	instance := u.Instance(xl) //可以将BaseDao理解为一个模板，注入context后变为实例，注入之前不可进行数据库操作
	if user, err = instance.Get(op.Eq("id", id)); err != nil {
		xl.Error("xxx")
		return nil, err
	}

	if user.Name == "hh" {
		user.Name = "xxx"
	}

	// Tip2: 这里可以直接使用实例
	err = instance.Save(user)
	return nil, err
}

// ListUsers 列表查询
func (u *UserService) ListUsers(xl *xlog.XLogger, limit, offset int, fuzzyName string) (users []*User, err error) {
	instance := u.Instance(xl)
	users, err = instance.List(
		op.Filter(op.Like("name", fuzzyName)),
		op.Offset(offset),
		op.Limit(limit),
		op.Order("age asc"))
	return
}

func (u *UserService) DeleteUser(xl *xlog.XLogger, userId string) (err error) {
	tx := u.BaseDao.Begin()
	instance := tx.Instance(xl)

	if err = instance.Delete(op.Eq("id", userId)); err != nil {
		err = tx.Rollback()
		return
	}

	//继承自原有的
	gradeDaoInstance := db.Extend[Grade](xl, tx.GetTransaction())
	if err = gradeDaoInstance.Delete(op.Eq("user_id", userId)); err != nil {
		err = tx.Rollback()
	}
	err = tx.Commit()
	return
}
