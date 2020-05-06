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

	r.GET("/contact",contact)
	r.POST("/addContacter",addContacter)
	r.POST("/editContact",editContact)
	r.POST("/deleteContact",deleteContact)

	r.GET("/getContactInfo",getContactInfo)
	r.POST("/addRecord",addRecord)
	r.POST("/editRecord",editRecord)
	r.POST("/deleteRecord",deleteRecord)



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
func contact(c *gin.Context) {
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
				logrus.Debug("contacts.Kind:",kind)
				logrus.Debug("contacts.Entitys = ",entityList[0])
				contacts.Kind = kind
				contacts.Entitys = entityList
				contactsList = append(contactsList, contacts)
			}
			contacts = new(common.Contacts)
			entityList = make([]*common.Entity,0)
			kind = p.Kind
		}
		entity := common.Entity{
			Id:      strconv.FormatInt(p.Id,10),
			Initial: p.Kind,
			Name:    p.Name,
			Info:*p.Info,
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

	logrus.Debug("concat name:",contact.Name)
	contactPersonModel.Kind = getKind(contact.Name)
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

//editContact
func editContact(c *gin.Context)  {
	contactId, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		logrus.Error("cat not get id err:",err)
		c.JSON(http.StatusBadRequest,"can not get id")
		return
	}

	var entity common.Entity
	err = c.ShouldBindJSON(&entity)
	if err != nil {
		logrus.Error("can not get data err:",err)
		c.JSON(http.StatusBadRequest,"can not get data")
		return
	}
	var PersonModel ledger_dsl.ContactPersonModel
	PersonModel.Id = contactId
	if entity.Name != "" {
		PersonModel.Name = entity.Name
		PersonModel.Kind = getKind(entity.Name)
	}
	if entity.Info != "" {
		PersonModel.Info =&entity.Info
	}
	_, err = dao.UpdateConcatPerson(&PersonModel)
	if err != nil {
		logrus.Error("call dao.UpdateConcatPerson err:",err)
		c.JSON(http.StatusExpectationFailed,"update contact person err")
		return
	}
	c.JSON(http.StatusOK,"update success !!")
}

func deleteContact(c *gin.Context) {
	contactId, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		logrus.Error("cat not get id err:",err)
		c.JSON(http.StatusBadRequest,"can not get id")
		return
	}
	personModel := ledger_dsl.ContactPersonModel{Id: contactId, Status: 2}
	_, err = dao.UpdateConcatPerson(&personModel)
	if err != nil {
		logrus.Error("update contact person err:",err)
		c.JSON(http.StatusExpectationFailed,"delete contact person err")
		return
	}
	c.JSON(http.StatusOK,"delete success !!")
}

//getContactInfo
func getContactInfo(c *gin.Context) {
	contactId, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		logrus.Error("cat not get id err:",err)
		c.JSON(http.StatusBadRequest,"can not get id")
		return
	}
	counts, err := dao.GetUserCount(contactId)
	if err != nil {
		logrus.Error("call dao.GetUserCount err:",err)
		c.JSON(http.StatusBadRequest,"can not get user count")
		return
	}

	var cDetailList []*common.CDetail
	var summary common.Summary
	var summaryList []*common.Summary
	var contactInfo common.GetContactInfo

	for _, count := range counts {
		cDetail := common.CDetail{
			DetailId: strconv.FormatInt(count.Id,10),
			Time:     count.Time,
			Money:    count.Money,
			Memo:     *count.Memo,
			Type:     count.Type,
		}
		if count.Type == 1 {
			summary.CCount += count.Money
		}else if count.Type == 2 {
			summary.GCount += count.Money
		}
		cDetailList = append(cDetailList,&cDetail)
	}
	summary.Difference = summary.CCount - summary.GCount
	summary.SId = 1
	summaryList = append(summaryList,&summary)
	contactInfo.CDetail = cDetailList
	contactInfo.Summary = summaryList

	js, err := json.Marshal(contactInfo)
	if err != nil {
		logrus.Error("json.Marshal contactInfo err:",err)
		c.JSON(http.StatusBadRequest,"json.Marshal contact err")
		return
	}

	c.JSON(http.StatusOK,string(js))

}

func addRecord(c *gin.Context) {
	var cDetail common.CDetail
	err := c.ShouldBindJSON(&cDetail)
	if err != nil {
		logrus.Error("can not get data err:",err)
		c.JSON(http.StatusBadRequest,"can not get data")
		return
	}
	if cDetail.Id == "" || cDetail.Money <=0 || cDetail.Type <= 0{
		logrus.Error("can not get id or money or type")
		c.JSON(http.StatusBadRequest,"can not get id or money or type")
		return
	}
	id, err := strconv.ParseInt(cDetail.Id, 10, 64)
	if err != nil {
		logrus.Error("can not get id err:",err)
		c.JSON(http.StatusBadRequest,"can not get id")
		return
	}
	userCountModel := ledger_dsl.UserCountModel{
		Id:              common.GenerateIds(1)[0],
		ContactPersonId: id,
		Money:           cDetail.Money,
		Type:            cDetail.Type,
		Createtime:      time.Now().UnixNano(),
		Updatetime:      time.Now().UnixNano(),
		Status:          1,
		Time:            cDetail.Time,
		Memo:            &cDetail.Memo,
	}
	_, err = dao.AddUserCount(&userCountModel)
	if err != nil {
		logrus.Error("call dao.AddUserCount err:",err)
		c.JSON(http.StatusExpectationFailed,"add user count err")
		return
	}
	c.JSON(http.StatusOK,"add success!!")
}

func editRecord(c *gin.Context) {
	userCountId, err := strconv.ParseInt(c.Query("detail_id"), 10, 64)
	if err != nil {
		logrus.Error("cat not get detail_id err:",err)
		c.JSON(http.StatusBadRequest,"can not get detail_id")
		return
	}

	var cDetail common.CDetail
	err = c.ShouldBindJSON(&cDetail)
	if err != nil {
		logrus.Error("can not get data err:",err)
		c.JSON(http.StatusBadRequest,"can not get data")
		return
	}
	userCountModel := ledger_dsl.UserCountModel{
		Id:              userCountId,
		Money:           cDetail.Money,
		Type:            cDetail.Type,
		Updatetime:      time.Now().UnixNano(),
		Time:            cDetail.Time,
		Memo:            &cDetail.Memo,
	}
	_, err = dao.UpdateUserCount(&userCountModel)
	if err != nil {
		logrus.Error("call dao.UpdateUserCount err:",err)
		c.JSON(http.StatusExpectationFailed,"edit user count err")
		return
	}
	c.JSON(http.StatusOK,"edit success !!")

}

func deleteRecord(c *gin.Context)  {
	userCountId, err := strconv.ParseInt(c.Query("detail_id"), 10, 64)
	if err != nil {
		logrus.Error("cat not get detail_id err:",err)
		c.JSON(http.StatusBadRequest,"can not get detail_id")
		return
	}

	userCountModel := ledger_dsl.UserCountModel{Id: userCountId, Status: 2}

	_, err = dao.UpdateUserCount(&userCountModel)
	if err != nil {
		logrus.Error("call dao.UpdateUserCount err:",err)
		c.JSON(http.StatusExpectationFailed,"delete user count err")
		return
	}
	c.JSON(http.StatusOK,"delete success !!")
}

func getKind(name string) string  {
	a := pinyin.NewArgs()
	pin := pinyin.Pinyin(name, a)
	logrus.Debug("pinyin:",pin)
	var kind string
	if len(pin) > 0 && len(pin[0]) >0 && len(pin[0][0]) > 0{
		kind = strings.ToUpper(pin[0][0][0:1])
		logrus.Info("parse chinese: contact person kind:",kind)
	}else {
		kind = strings.ToUpper(name[0:1])
		logrus.Info("can not parse :contact person kind:",kind)
	}
	return kind
}
