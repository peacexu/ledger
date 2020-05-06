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
	Id      string	`json:"id"`
	Initial string	`json:"initial"`
	Name    string	`json:"name"`
	Info string `json:"info"`
}
type Contacts struct {
	Kind    string	`json:"kind"`
	Entitys []*Entity `json:"entitys"`
}
type ContactsResponse struct {
	ContactsList []*Contacts `json:"contactsList"`
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

//联系人详情
type CDetail struct {
	DetailId string `json:"detail_id"`
	Id string `json:"id"`
	Time string	`json:"time"`
	Money int32	`json:"money"`
	Memo string	`json:"memo"`
	Type int32	`json:"type"`
}

type Summary struct {
	SId int32 `json:"s_id"`
	CCount int32 `json:"c_count"`
	GCount int32 `json:"g_count"`
	Difference int32 `json:"difference"`
}

type GetContactInfo struct {
	CDetail []*CDetail `json:"c_detail"`
	Summary []*Summary `json:"summary"`
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