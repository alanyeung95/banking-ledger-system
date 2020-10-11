package main

import (
// /	"fmt"
    "net/http"
    "testing"
    "github.com/gavv/httpexpect"
)


var testurl string ="http://127.0.0.1:3000"
var testingAccountID string = ""

func TestRootAPI(t *testing.T) {
    e:= httpexpect.New(t, testurl)  
    e.GET("/").   //ge請求
        Expect().
        Status(http.StatusOK)
}

func TestCreatingAccount(t *testing.T) {
    e:= httpexpect.New(t, testurl) 
    postdata:= map[string]interface{}{
		"name": "testing_user_name",
		"password": "xxxxx",
	}
    contentType := "application/json;charset=utf-8"

    obj :=  e.POST("/accounts/create").  
    WithHeader("ContentType", contentType). 
    WithJSON(postdata).  
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Object().
	ContainsKey("name").
	ValueEqual("name","testing_user_name")   
	
	testingAccountID = obj.Value("id").String().Raw()
}

func TestDeposit(t *testing.T) {
    e:= httpexpect.New(t, testurl) 
    postdata:= map[string]interface{}{
		"operation": "deposit",
		"body": map[string]interface{} {
			"amount": 200,
		},
	}
    contentType := "application/json;charset=utf-8"

    e.POST("/transactions").  
	WithHeader("ContentType", contentType). 
	WithHeader("account_id", testingAccountID).
    WithJSON(postdata).  
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Object().
	ContainsKey("id").
	ValueEqual("id",testingAccountID)   
}

func TestCleanUp(t *testing.T) {
	e:= httpexpect.New(t, testurl) 

	e.DELETE("/accounts/" + testingAccountID).
	Expect().
	Status(http.StatusOK)
	
	e.DELETE("/accounts/" + testingAccountID + "/transactions").
	Expect().
    Status(http.StatusOK)
}