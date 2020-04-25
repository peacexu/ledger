package service

type Entity struct {
	Id      int32
	Initial byte
	Name    string
}

type Contacts struct {
	Kind    byte
	Entitys []*Entity
}
