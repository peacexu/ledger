package ledger_dsl

import (
	. "github.com/lqs/sqlingo"
	"reflect"
)

type tContactPerson struct {
	Table
	Id         fContactPersonId
	UserId     fContactPersonUserId
	Name       fContactPersonName
	Kind       fContactPersonKind
	Info       fContactPersonInfo
	Status     fContactPersonStatus
	Createtime fContactPersonCreatetime
	Updatetime fContactPersonUpdatetime
}

type fContactPersonId struct{ NumberField }
type fContactPersonUserId struct{ NumberField }
type fContactPersonName struct{ StringField }
type fContactPersonKind struct{ StringField }
type fContactPersonInfo struct{ StringField }
type fContactPersonStatus struct{ NumberField }
type fContactPersonCreatetime struct{ NumberField }
type fContactPersonUpdatetime struct{ NumberField }

var ContactPerson = tContactPerson{
	Table:      NewTable("contact_person"),
	Id:         fContactPersonId{NewNumberField("contact_person", "id")},
	UserId:     fContactPersonUserId{NewNumberField("contact_person", "user_id")},
	Name:       fContactPersonName{NewStringField("contact_person", "name")},
	Kind:       fContactPersonKind{NewStringField("contact_person", "kind")},
	Info:       fContactPersonInfo{NewStringField("contact_person", "info")},
	Status:     fContactPersonStatus{NewNumberField("contact_person", "status")},
	Createtime: fContactPersonCreatetime{NewNumberField("contact_person", "createtime")},
	Updatetime: fContactPersonUpdatetime{NewNumberField("contact_person", "updatetime")},
}

func (t tContactPerson) GetFields() []Field {
	return []Field{t.Id, t.UserId, t.Name, t.Kind, t.Info, t.Status, t.Createtime, t.Updatetime}
}

func (t tContactPerson) GetFieldByName(name string) Field {
	r := reflect.ValueOf(t)
	f := reflect.Indirect(r).FieldByName(CamelName(name))
	if !f.IsValid() {
		return nil
	}
	if field, ok := f.Interface().(Field); ok {
		return field
	}
	return nil
}

func (t tContactPerson) GetFieldsSQL() string {
	return "`id`, `user_id`, `name`, `kind`, `info`, `status`, `createtime`, `updatetime`"
}

func (t tContactPerson) GetFullFieldsSQL() string {
	return "`contact_person`.`id`, `contact_person`.`user_id`, `contact_person`.`name`, `contact_person`.`kind`, `contact_person`.`info`, `contact_person`.`status`, `contact_person`.`createtime`, `contact_person`.`updatetime`"
}

type ContactPersonModel struct {
	Id         int64
	UserId     int64
	Name       string
	Kind       string
	Info       *string
	Status     int32
	Createtime int64
	Updatetime int64
}

func (m ContactPersonModel) GetTable() Table {
	return ContactPerson
}

func (m ContactPersonModel) GetValues() []interface{} {
	return []interface{}{m.Id, m.UserId, m.Name, m.Kind, m.Info, m.Status, m.Createtime, m.Updatetime}
}

type tUser struct {
	Table
	Id         fUserId
	Phone      fUserPhone
	Pwd        fUserPwd
	Createtime fUserCreatetime
	Updatetime fUserUpdatetime
	Status     fUserStatus
}

type fUserId struct{ NumberField }
type fUserPhone struct{ StringField }
type fUserPwd struct{ StringField }
type fUserCreatetime struct{ NumberField }
type fUserUpdatetime struct{ NumberField }
type fUserStatus struct{ NumberField }

var User = tUser{
	Table:      NewTable("user"),
	Id:         fUserId{NewNumberField("user", "id")},
	Phone:      fUserPhone{NewStringField("user", "phone")},
	Pwd:        fUserPwd{NewStringField("user", "pwd")},
	Createtime: fUserCreatetime{NewNumberField("user", "createtime")},
	Updatetime: fUserUpdatetime{NewNumberField("user", "updatetime")},
	Status:     fUserStatus{NewNumberField("user", "status")},
}

func (t tUser) GetFields() []Field {
	return []Field{t.Id, t.Phone, t.Pwd, t.Createtime, t.Updatetime, t.Status}
}

func (t tUser) GetFieldByName(name string) Field {
	r := reflect.ValueOf(t)
	f := reflect.Indirect(r).FieldByName(CamelName(name))
	if !f.IsValid() {
		return nil
	}
	if field, ok := f.Interface().(Field); ok {
		return field
	}
	return nil
}

func (t tUser) GetFieldsSQL() string {
	return "`id`, `phone`, `pwd`, `createtime`, `updatetime`, `status`"
}

func (t tUser) GetFullFieldsSQL() string {
	return "`user`.`id`, `user`.`phone`, `user`.`pwd`, `user`.`createtime`, `user`.`updatetime`, `user`.`status`"
}

type UserModel struct {
	Id         int64
	Phone      string
	Pwd        string
	Createtime int64
	Updatetime int64
	Status     int32
}

func (m UserModel) GetTable() Table {
	return User
}

func (m UserModel) GetValues() []interface{} {
	return []interface{}{m.Id, m.Phone, m.Pwd, m.Createtime, m.Updatetime, m.Status}
}

type tUserCount struct {
	Table
	Id              fUserCountId
	ContactPersonId fUserCountContactPersonId
	Money           fUserCountMoney
	Type            fUserCountType
	Createtime      fUserCountCreatetime
	Updatetime      fUserCountUpdatetime
	Status          fUserCountStatus
	Time            fUserCountTime
	Memo            fUserCountMemo
}

type fUserCountId struct{ NumberField }
type fUserCountContactPersonId struct{ NumberField }
type fUserCountMoney struct{ NumberField }
type fUserCountType struct{ NumberField }
type fUserCountCreatetime struct{ NumberField }
type fUserCountUpdatetime struct{ NumberField }
type fUserCountStatus struct{ NumberField }
type fUserCountTime struct{ StringField }
type fUserCountMemo struct{ StringField }

var UserCount = tUserCount{
	Table:           NewTable("user_count"),
	Id:              fUserCountId{NewNumberField("user_count", "id")},
	ContactPersonId: fUserCountContactPersonId{NewNumberField("user_count", "contact_person_id")},
	Money:           fUserCountMoney{NewNumberField("user_count", "money")},
	Type:            fUserCountType{NewNumberField("user_count", "type")},
	Createtime:      fUserCountCreatetime{NewNumberField("user_count", "createtime")},
	Updatetime:      fUserCountUpdatetime{NewNumberField("user_count", "updatetime")},
	Status:          fUserCountStatus{NewNumberField("user_count", "status")},
	Time:            fUserCountTime{NewStringField("user_count", "time")},
	Memo:            fUserCountMemo{NewStringField("user_count", "memo")},
}

func (t tUserCount) GetFields() []Field {
	return []Field{t.Id, t.ContactPersonId, t.Money, t.Type, t.Createtime, t.Updatetime, t.Status, t.Time, t.Memo}
}

func (t tUserCount) GetFieldByName(name string) Field {
	r := reflect.ValueOf(t)
	f := reflect.Indirect(r).FieldByName(CamelName(name))
	if !f.IsValid() {
		return nil
	}
	if field, ok := f.Interface().(Field); ok {
		return field
	}
	return nil
}

func (t tUserCount) GetFieldsSQL() string {
	return "`id`, `contact_person_id`, `money`, `type`, `createtime`, `updatetime`, `status`, `time`, `memo`"
}

func (t tUserCount) GetFullFieldsSQL() string {
	return "`user_count`.`id`, `user_count`.`contact_person_id`, `user_count`.`money`, `user_count`.`type`, `user_count`.`createtime`, `user_count`.`updatetime`, `user_count`.`status`, `user_count`.`time`, `user_count`.`memo`"
}

type UserCountModel struct {
	Id              int64
	ContactPersonId int64
	Money           int32
	Type            int32
	Createtime      int64
	Updatetime      int64
	Status          int32
	Time            string
	Memo            *string
}

func (m UserCountModel) GetTable() Table {
	return UserCount
}

func (m UserCountModel) GetValues() []interface{} {
	return []interface{}{m.Id, m.ContactPersonId, m.Money, m.Type, m.Createtime, m.Updatetime, m.Status, m.Time, m.Memo}
}

var tableMap = map[string]Table{
	"contact_person": ContactPerson,
	"user":           User,
	"user_count":     UserCount,
}

func GetTable(name string) Table {
	if table, ok := tableMap[name]; ok {
		return table
	}
	return nil
}
