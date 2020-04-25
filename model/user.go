package model

type User struct {
	Id         int32  `json:"id"`
	Phone      string `json:"phone"`
	Pwd        string `json:"pwd"`
	Createtime int64  `json:"createtime"`
	Updatetime int64  `json:"updatetime"`
	Status     int32  `json:"status"`
}

type Count struct {
	Id         int32   `json:"id"`
	UserId     int32   `json:"user_id"`
	Name       string  `json:"name"`
	Money      int32   `json:"money"`
	Type       int32   `json:"type"`
	Kind       byte    `json:"kind"`
	Memo       *string `json:"memo"`
	Createtime int64   `json:"createtime"`
	Updatetime int64   `json:"updatetime"`
	Status     int32   `json:"status"`
	Info       *string `json:"info"`
	Time       int64   `json:"time"`
}
