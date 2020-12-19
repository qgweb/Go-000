package dao

import (
	"Go-000/Week04/internal/conf"
	"Go-000/Week04/internal/model"
	"context"
	"database/sql"
	"fmt"

	xerror "errors"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

var ErrNotFound = xerror.New("record not found")

type Dao struct {
	db *sql.DB
}

func connectMysql(config conf.Database) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?loc=Local&"+
		"charset=utf8&parseTime=true",
		config.User, config.Password, config.Host,
		config.Port, config.Name,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(config.Idle)
	db.SetMaxOpenConns(config.Active)
	return db
}

func New(config conf.Config) *Dao {
	dao := &Dao{
		db: connectMysql(config.DB),
	}
	return dao
}

func (d *Dao) Close() {
	if d.db != nil {
		_ = d.db.Close()
	}
}

func (d *Dao) GetUser(ctx context.Context, name, pwd string) (user model.User, err error) {
	loginSql := "select * from x_user where  name=? and passwd=? limit 1"
	rows := d.db.QueryRowContext(ctx, loginSql, name, pwd)
	if err := rows.Err(); err != nil {
		return user, errors.Wrapf(err, "sql:%s,name:%s,pwd:%s", loginSql, name, pwd)
	}
	if err := rows.Scan(&user.ID, &user.Name, &user.Passwd); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.Wrapf(ErrNotFound, "name:%s,pwd:%s", name, pwd)
		}
		return user, errors.Wrapf(err, "scan xerror,sql:%s,name:%s,pwd:%s", loginSql, name, pwd)
	}
	return user, nil
}
