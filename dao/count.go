package dao

import (
	"fmt"
	"ledger/model"
	"time"
)

func AddCount(count *model.Count) (int64, error) {
	if count.UserId == 0 || count.Name == "" || count.Money == 0 || count.Type == 0 {
		return 0, fmt.Errorf("add count feild err")
	}
	createTime := time.Now().UnixNano()

	sqlstr := `insert into 
					count(user_id,name,money,type,kind,memo,createtime,updatetime,status)
				values(?,?,?,?,?,?,?,?,?)`
	res, err := DB.Exec(sqlstr, count.UserId, count.Name, count.Money, count.Type, count.Kind, count.Memo, createTime, createTime, count.Status)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, err
}

func GetCount(id int32) (count *model.Count, err error) {
	sqlstr := `select * from count where id = ? and status =1`
	err = DB.Get(count, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return count, nil
}

func UpdateCount(count model.Count) (int64, error) {
	updateTime := time.Now().UnixNano()
	sqlstr := `update count set user_id=?,name=?,money=?,type=?,kind=?,memo=?,updatetime=?,status=? where id = ?`
	result, err := DB.Exec(sqlstr, count.UserId, count.Name, count.Money, count.Kind, count.Memo, updateTime, count.Status, count.Id)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func GetAllNameByUserId(userId int32) ([]*model.Count, error) {
	var counts []*model.Count
	sqlstr := `select distinct (name) from count where userId = ? order by kind`
	err := DB.Get(counts, sqlstr, userId)
	if err != nil {
		return nil, err
	}
	return counts, nil

}
