package service

import (
	"encoding/json"
	"ledger/dao"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

}

//getUserInfo
func getUserInfo(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := dao.GetUserById(int32(userId))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}

	jsonUser, _ := json.Marshal(user)
	c.JSON(http.StatusOK, string(jsonUser))
}

//get contacter info
func getContacterInfo(c gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}

	counts, err := dao.GetAllNameByUserId(int32(userId))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}

	kind := byte('A')
	var contactsList []*Contacts
	var contacts *Contacts = new(Contacts)
	var entityList []*Entity
	for _, count := range counts {
		if count.Kind != kind {
			if entityList != nil && len(entityList) > 0 {
				contacts.Kind = kind
				contacts.Entitys = entityList
				contactsList = append(contactsList, contacts)
			}
			contacts = new(Contacts)
			entityList = entityList[0:0]
			kind = count.Kind
		}
		entity := Entity{
			Id:      count.Id,
			Initial: count.Kind,
			Name:    count.Name,
		}
		entityList = append(entityList, &entity)
	}

	jsonList, _ := json.Marshal(contactsList)
	c.JSON(http.StatusOK, string(jsonList))
}

//add con
