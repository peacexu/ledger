package test

import (
	"ledger/service"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	router := service.SetupRouter()
	w := httptest.NewRecorder()
	reader := strings.NewReader(`{"phone":"13269338513","password":"123456"}`)
	req, _ := http.NewRequest("POST", "/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	router.ServeHTTP(w,req)
	log.Println(w.Body.String())
}

func TestRegister(t *testing.T) {
	router := service.SetupRouter()
	w := httptest.NewRecorder()
	reader := strings.NewReader(`{"phone":"13269338513","password":"123456"}`)
	req, _ := http.NewRequest("POST", "/register", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	router.ServeHTTP(w,req)
	log.Println(w.Body.String())
}

func TestAddContacter(t *testing.T) {
	router := service.SetupRouter()
	w := httptest.NewRecorder()
	reader := strings.NewReader(`{"name":"朱大肠","info":"又大又长","time":"2020-05-01","memo":"杀猪","money":1000,"type":1}`)
	req, _ := http.NewRequest("POST", "/addContacter?user_id=441629910608379904", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	router.ServeHTTP(w,req)
	log.Println(w.Body.String())
}

func TestGetContacterInfo(t *testing.T) {
	router := service.SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/contact?user_id=441813108315389950&page=0", nil)
	router.ServeHTTP(w,req)
	log.Println(w.Body.String())
}

