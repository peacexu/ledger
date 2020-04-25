package dao

import (
	"fmt"
	"ledger/model"
	"time"
)

//CREATE TABLE `user` (
//`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
//`phone` varchar(11) CHARACTER SET utf8mb4 NOT NULL DEFAULT '',
//`pwd` varchar(11) CHARACTER SET utf8 NOT NULL DEFAULT '',
//`createtime` bigint(20) NOT NULL,
//`updatetime` bigint(20) NOT NULL,
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=latin1;;
func AddUser(user *model.User) (int64, error) {
	if user == nil || user.Phone == "" || user.Pwd == "" {
		err := fmt.Errorf("invalid user parameter")
		return -1, err
	}

	createTime := time.Now().UnixNano()

	sqlstr := `insert into
					user(phone,pwd,createtime,updatetime,status)
				values (?,?,?,?,?)`
	result, err := DB.Exec(sqlstr, user.Phone, user.Pwd, createTime, createTime, user.Status)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func GetUserById(id int32) (*model.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invaild id")
	}
	sqlstr := `select * from user where id = ? and status = 1`
	var user *model.User
	err := DB.Get(&user, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
