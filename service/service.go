package service

import (
	"encoding/json"
	"github.com/mozillazg/go-pinyin"
	"github.com/sirupsen/logrus"
	"ledger/common"
	"ledger/dao"
	ledger_dsl "ledger/generated/sqlingo"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	//日志初始化
	logrus.SetLevel(logrus.DebugLevel)
	//logrus.SetLevel(logrus.InfoLevel)
	myFormatter := new(logrus.TextFormatter)
	myFormatter.FullTimestamp = true
	myFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(myFormatter)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/login",login)
	r.POST("/register",register)
	r.GET("/getContacterInfo",getContacterInfo)
	r.POST("/addContacter",addContacter)
	return r
}

//login
func login(c *gin.Context) {
	var loginMedel common.Login
	err := c.ShouldBindJSON(&loginMedel)
	if err != nil {
		logrus.Errorf("login data err:",err)
		c.JSON(http.StatusBadRequest,"invalid data")
	}

	user, err := dao.GetUserByPhone(loginMedel.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.MessageResponse{Msg:"invalid phone"})
		return
	}

	if user == nil {
		logrus.Info("can not get this user")
		c.JSON(http.StatusBadRequest,common.MessageResponse{Msg:"this phone can not register"})
		return
	}

	if loginMedel.Password== user.Pwd {
		m := make(map[string]int64)
		m["user_id"] = user.Id
		js, _ := json.Marshal(m)
		logrus.Info("login success user:",user)
		c.JSON(http.StatusOK, string(js))
		return
	}else {
		logrus.Info("error password:",loginMedel.Password)
		c.JSON(http.StatusBadRequest,common.MessageResponse{Msg:"error password"})
		return
	}

}

//register
func register(c *gin.Context) {
	var registerMedel common.Register
	err := c.ShouldBindJSON(&registerMedel)
	if err != nil {
		logrus.Errorf("login data err:",err)
		c.JSON(http.StatusBadRequest,"invalid data")
	}

	user, err := dao.GetUserByPhone(registerMedel.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.MessageResponse{Msg:"invalid phone"})
		return
	}

	if user != nil {
		logrus.Info("this phone is already register")
		c.JSON(http.StatusBadRequest,common.MessageResponse{Msg:"this phone is already register"})
		return
	}

	userModel := ledger_dsl.UserModel{
		Id:         common.GenerateIds(1)[0],
		Phone:      registerMedel.Phone,
		Pwd:        registerMedel.Password,
		Createtime: time.Now().UnixNano(),
		Updatetime: time.Now().UnixNano(),
		Status:     1,
	}
	_, err = dao.AddUser(&userModel)
	if err != nil {
		logrus.Error("call dao.AddUser err:",err)
		c.JSON(http.StatusExpectationFailed,"add user err")
		return
	}
	c.JSON(http.StatusOK,"register success" )
}

//get contacter info
func getContacterInfo(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}

	page, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid page")
		return
	}

	limitOffset := common.LimitOffset{
		Limit:  10,
		Offset: int(page)*10,
	}
	contactPersons,total, err := dao.GetConcatPersonByUserId(userId,&limitOffset)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}

	kind := "A"
	var contactsList []*common.Contacts
	var contacts *common.Contacts = new(common.Contacts)
	var entityList []*common.Entity
	for _, p := range contactPersons {
		if p.Kind != kind {
			if entityList != nil && len(entityList) > 0 {
				contacts.Kind = kind
				contacts.Entitys = entityList
				contactsList = append(contactsList, contacts)
			}
			contacts = new(common.Contacts)
			entityList = entityList[0:0]
			kind = p.Kind
		}
		entity := common.Entity{
			Id:      p.Id,
			Initial: p.Kind,
			Name:    p.Name,
		}
		entityList = append(entityList, &entity)
	}
	//最后一个添加
	if entityList != nil && len(entityList) > 0 {
		contacts.Kind = kind
		contacts.Entitys = entityList
		contactsList = append(contactsList, contacts)
	}

	resp := common.ContactsResponse{ContactsList: contactsList, Total: total}
	js, _ := json.Marshal(resp)
	c.JSON(http.StatusOK, string(js))
	return
}

//add Contacter
func addContacter(c *gin.Context)  {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}
	logrus.Debug("userId:",userId)

	var contact common.AddContacter
	if err = c.ShouldBindJSON(&contact);err!=nil {
		c.JSON(http.StatusBadRequest,"invalid data")
		return
	}
	logrus.Debug("contact:",contact)

	ids := common.GenerateIds(2)
	contactPersonModel := ledger_dsl.ContactPersonModel{
		Id:         ids[0],
		UserId:     userId,
		Name:       contact.Name,
		Info:       &contact.Info,
		Status:     1,
		Createtime: time.Now().UnixNano(),
		Updatetime: time.Now().UnixNano(),
	}

	a := pinyin.NewArgs()
	pin := pinyin.Pinyin(contact.Name, a)
	logrus.Debug("pinyin:",pin)
	if len(pin) > 0 && len(pin[0]) >0 && len(pin[0][0]) > 0{
		contactPersonModel.Kind = strings.ToUpper(pin[0][0][0:1])
		logrus.Info("contact person kind:",contactPersonModel.Kind)
	}else {
		contactPersonModel.Kind = "Z"
	}
	_, err = dao.AddConcatPerson(&contactPersonModel)
	if err != nil {
		logrus.Errorf("call dao.AddConcatPerson err:",err)
		c.JSON(http.StatusExpectationFailed,"添加联系人错误")
		return
	}

	userCountModel := ledger_dsl.UserCountModel{
		Id:              ids[1],
		ContactPersonId: ids[0],
		Money:           contact.Money,
		Type:            contact.Type,
		Createtime:      time.Now().UnixNano(),
		Updatetime:      time.Now().UnixNano(),
		Status:          1,
		Time:            contact.Time,
		Memo:            &contact.Memo,
	}

	_, err = dao.AddUserCount(&userCountModel)
	if err != nil {
		logrus.Errorf("call dao.AddUserCount err:",err)
		c.JSON(http.StatusExpectationFailed,"添加联系人记录错误")
		return
	}

	c.JSON(http.StatusOK,"add success!")
	return

}

