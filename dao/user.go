package dao

import (
	"ledger/common"
	ledger_dsl "ledger/generated/sqlingo"
	"github.com/lqs/sqlingo"
	"github.com/sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var db sqlingo.Database
func init() {
	var err error
	db, err = sqlingo.Open("mysql", "root:123456@tcp(122.112.242.150:3306)/ledger")
	if err != nil {
		logrus.Error("open mysql err:",err)
		return
	}
}

func GetUserByPhone(phone string) (*ledger_dsl.UserModel, error) {
	expression := ledger_dsl.User.Phone.Equals(phone).And(ledger_dsl.User.Status.Equals(1))
	var user *ledger_dsl.UserModel
	_, err := db.SelectFrom(ledger_dsl.User).Where(expression).FetchFirst(&user)
	if err != nil {
		logrus.Error("select pwd from user err:",err)
		return nil,err
	}

	return user,nil
}

func GetConcatPersonByUserId(userId int64,limitOffset *common.LimitOffset) ([]*ledger_dsl.ContactPersonModel,int,error)  {
	expression := ledger_dsl.ContactPerson.UserId.Equals(userId).And(ledger_dsl.ContactPerson.Status.Equals(1))
	var concatPersons []*ledger_dsl.ContactPersonModel

	count, err := db.SelectFrom(ledger_dsl.ContactPerson).Where(expression).Count()
	if err != nil {
		logrus.Error("select count concat person err:",err)
		return nil,0,err
	}
	_, err = db.SelectFrom(ledger_dsl.ContactPerson).Where(expression).OrderBy(ledger_dsl.ContactPerson.Kind).Limit(limitOffset.Limit).Offset(limitOffset.Offset).FetchAll(&concatPersons)
	if err != nil {
		logrus.Error("select concat persons err:",err)
		return nil,0,err
	}

	if concatPersons == nil || len(concatPersons) == 0 {
		return nil,count,nil
	}else {
		return concatPersons,count,nil
	}
}

func AddConcatPerson(contactPerson *ledger_dsl.ContactPersonModel) (bool ,error) {
	_,err := db.InsertInto(ledger_dsl.ContactPerson).Models(contactPerson).Execute()
	if err != nil {
		logrus.Errorf("insert %v to contact person err:%v",contactPerson,err)
		return false, err
	}
	logrus.Infof("insert %v into contact person success",contactPerson)
	return true,nil
}

func AddUserCount(userCount *ledger_dsl.UserCountModel) (bool,error)  {
	_, err := db.InsertInto(ledger_dsl.UserCount).Models(userCount).Execute()
	if err != nil {
		logrus.Errorf("insert %v into userCount err:%v",userCount,err)
		return false,err
	}
	logrus.Infof("insert %v into userCount success",userCount)
	return true,nil

}

func AddUser(user *ledger_dsl.UserModel) (bool, error) {
	_, err := db.InsertInto(ledger_dsl.User).Models(user).Execute()
	if err != nil {
		logrus.Errorf("insert %v into user err:%v",user,err)
		return false,err
	}

	logrus.Infof("insert %v into user success",user)
	return true,nil
}