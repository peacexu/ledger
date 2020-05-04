package common

//login
type Login struct {
	Phone string `json:"phone"`
	Password string `json:"password"`
}

//register
type Register struct {
	Phone string	`json:"phone"`
	Password string `json:"password"`
}

//联系人列表
type Entity struct {
	Id      int64	`json:"id"`
	Initial string	`json:"initial"`
	Name    string	`json:"name"`
}
type Contacts struct {
	Kind    string	`json:"kind"`
	Entitys []*Entity `json:"entitys"`
}
type ContactsResponse struct {
	ContactsList []*Contacts `json:"contacts_list"`
	Total int	`json:"total"`
}

//添加联系人
type AddContacter struct {
	Name string `json:"name"`
	Info string `json:"info"`
	Time string `json:"time"`
	Memo string `json:"memo"`
	Money int32 `json:"money"`
	Type int32 `json:"type"`
}


//分页
type LimitOffset struct {
	Limit int
	Offset int
}

//返回消息
type MessageResponse struct {
	Msg string `json:"msg"`
}

func GenerateIds(num int64) []int64 {
	snowflake := new(Snowflake)

	ids:= snowflake.BatchGenerate(num)
	return ids

}