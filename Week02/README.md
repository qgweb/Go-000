## 问题
我们在数据库操作的时候，比如 `dao` 层中当遇到一个 `sql.ErrNoRows` 的时候，是否应该 `Wrap` 这个 `error`，抛给上层。为什么？应该怎么做请写出代码

## 总结
在DAO层的sql.ErrNoRows 需要特殊处理， DAO层是数据操作抽象层，可能底层数据存储是mysql，也许后面就变成mongodb，
sql.ErrNoRows的报错强依赖于某种关系型数据库，不具有抽象性，所以单独定义一个数据不存在的error来规避依赖性。
所以不需要warp这个错误， 如果不是ErrNoRows 就需要包一层，把参数给包进去，方便后续排查问题。



## 代码
```go

package global

type ErrNotFound = errors.New("record not found")

package Dao

type AccountOperator struct {}

func (a *AccountOperator) FindUser(uid int) (user *model.Account, err error) {
	err = DB.Table("t_account").Where("id = ?", uid).Find(user).Error
	if errors.Is(err, sql.ErrNoRows) {
    return nil, fmt.Errorf("%w", global.ErrNotFound)
	}
  
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("userId: %d", uid))
    return nil, err
	}
  
	return user, nil
}

package Biz

import "DAO"

type AccountBusiness struct {
  ao *Dao.AccountOperator
}

func (ab *AccountBusiness) GetAccountDetail(uid int) (user *model.Account,err error) {
  user, err := ab.ao.FindUser(uid)
  if errors.Is(err, global.ErrNotFound) { 
      //如果需要降级
      return mock.NewAccount(),nil 
      //不需要
      return user, errors.Wrap(err, fmt.Sprintf("uid:%d", uid))
  }
  if err !=nil {
    return user, err
  }
}


