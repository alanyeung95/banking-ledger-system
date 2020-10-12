package main

import (
	"os"
    "net/http"
    "testing"
    "github.com/gavv/httpexpect"
)

var testurl string = os.Getenv("API_TEST_DOMAIN") 

func TestRootAPI(t *testing.T) {
    e:= httpexpect.New(t, testurl)  
    e.GET("/"). 
    Expect().
    Status(http.StatusOK)
}

func TestDepositAndWithdrawAndGetBalance(t *testing.T) {
    e:= httpexpect.New(t, testurl) 

	// create account
	postdata:= map[string]interface{}{
		"name": "testing_user_name_1",
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
	ValueEqual("name","testing_user_name_1")   
	testingAccountOneID := obj.Value("id").String().Raw()

	// test get balance
    e.GET("/accounts/" + testingAccountOneID).  
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Object().
	ContainsKey("id").
	ValueEqual("id",testingAccountOneID).
	ContainsKey("balance").
	ValueEqual("balance",0)     

	// test deposit
    postdata= map[string]interface{}{
		"operation": "deposit",
		"body": map[string]interface{} {
			"amount": 200,
		},
	}

    e.POST("/transactions").  
	WithHeader("ContentType", contentType). 
	WithHeader("account_id", testingAccountOneID).
    WithJSON(postdata).  
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Object().
	ContainsKey("id").
	ValueEqual("id",testingAccountOneID).
	ContainsKey("balance").
	ValueEqual("balance",200)     

	// test withdraw
	postdata = map[string]interface{}{
		"operation": "withdraw",
		"body": map[string]interface{} {
			"amount": 100,
		},
	}

    e.POST("/transactions").  
	WithHeader("ContentType", contentType). 
	WithHeader("account_id", testingAccountOneID).
    WithJSON(postdata).  
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Object().
	ContainsKey("id").
	ValueEqual("id",testingAccountOneID).
	ContainsKey("balance").
	ValueEqual("balance",100)    

	// clean up
	e.DELETE("/accounts/" + testingAccountOneID).
	Expect().
	Status(http.StatusOK)
	
	e.DELETE("/accounts/" + testingAccountOneID + "/transactions").
	Expect().
	Status(http.StatusOK)
}


func TestTransfer(t *testing.T) {
	e:= httpexpect.New(t, testurl) 
	contentType := "application/json;charset=utf-8"
	
	// creating first user
	postdata2:= map[string]interface{}{
		"name": "testing_user_name_2",
		"password": "xxxxx",
	}

    obj2 :=  e.POST("/accounts/create").  
    WithHeader("ContentType", contentType). 
    WithJSON(postdata2).  
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Object().
	ContainsKey("name").
	ValueEqual("name","testing_user_name_2")   
	
	testingAccountTwoID := obj2.Value("id").String().Raw()

	// creating second user
	postdata3:= map[string]interface{}{
		"name": "testing_user_name_3",
		"password": "xxxxx",
	}

	obj3 :=  e.POST("/accounts/create").  
	WithHeader("ContentType", contentType). 
	WithJSON(postdata3).  
	Expect().
	Status(http.StatusOK). 
	JSON().                 
	Object().
	ContainsKey("name").
	ValueEqual("name","testing_user_name_3")   
	
	testingAccountThreeID := obj3.Value("id").String().Raw()

	// making deposit
	postdata:= map[string]interface{}{
		"operation": "deposit",
		"body": map[string]interface{} {
			"amount": 100,
		},
	}

	e.POST("/transactions").  
	WithHeader("ContentType", contentType). 
	WithHeader("account_id", testingAccountTwoID).
	WithJSON(postdata).  
	Expect()

	// making transfer
    postdata= map[string]interface{}{
		"operation": "transfer",
		"body": map[string]interface{} {
			"amount": 100,
			"to": testingAccountThreeID,
		},
	}
    contentType = "application/json;charset=utf-8"

    e.POST("/transactions").  
	WithHeader("ContentType", contentType). 
	WithHeader("account_id", testingAccountTwoID).
    WithJSON(postdata).  
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Object().
	ContainsKey("id").
	ValueEqual("id",testingAccountTwoID).
	ContainsKey("balance").
	ValueEqual("balance",0)   
	
	e.GET("/accounts/" + testingAccountThreeID).  
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Object().
	ContainsKey("id").
	ValueEqual("id",testingAccountThreeID).
	ContainsKey("balance").
	ValueEqual("balance",100)   
	
	// clean up
	e.DELETE("/accounts/" + testingAccountTwoID).
	Expect().
	Status(http.StatusOK)
	
	e.DELETE("/accounts/" + testingAccountTwoID + "/transactions").
	Expect().
	Status(http.StatusOK)
	
	e.DELETE("/accounts/" + testingAccountThreeID).
	Expect().
	Status(http.StatusOK)
	
	e.DELETE("/accounts/" + testingAccountThreeID + "/transactions").
	Expect().
	Status(http.StatusOK)
}


func TestViewCustomerTransaction(t *testing.T) {
    e:= httpexpect.New(t, testurl) 
    contentType := "application/json;charset=utf-8"

	// creating  user
	postdata2:= map[string]interface{}{
		"name": "testing_user_name_2",
		"password": "xxxxx",
	}

	obj2 :=  e.POST("/accounts/create").  
	WithHeader("ContentType", contentType). 
	WithJSON(postdata2).  
	Expect().
	Status(http.StatusOK). 
	JSON().                 
	Object().
	ContainsKey("name").
	ValueEqual("name","testing_user_name_2")   

	testingAccountID := obj2.Value("id").String().Raw()

	// making deposit
	postdata:= map[string]interface{}{
		"operation": "deposit",
		"body": map[string]interface{} {
			"amount": 100,
		},
	}

	e.POST("/transactions").  
	WithHeader("ContentType", contentType). 
	WithHeader("account_id", testingAccountID).
	WithJSON(postdata).  
	Expect()

	objs := e.GET("/transactions").
	WithQuery("account_id", testingAccountID).  
	WithQuery("asc", 1).
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Array()

	objs.Element(0).Object().ValueEqual("operation", "deposit")

	// clean up
	e.DELETE("/accounts/" + testingAccountID).
	Expect().
	Status(http.StatusOK)
	
	e.DELETE("/accounts/" + testingAccountID + "/transactions").
	Expect().
	Status(http.StatusOK)
}


func TestFixTransactionAndViewOperationTeamTransaction(t *testing.T) {
    e:= httpexpect.New(t, testurl) 

	// creating customer account
	postdata:= map[string]interface{}{
		"name": "testing_user_name_3",
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
	ValueEqual("name","testing_user_name_3")   
	
	testingCustomerAccountThreeID := obj.Value("id").String().Raw()

	// creating operation team account
	postdata4:= map[string]interface{}{
		"name": "testing_user_name_4",
		"password": "xxxxx",
	}

    obj4 :=  e.POST("/accounts/create-admin").  
    WithHeader("ContentType", contentType). 
    WithJSON(postdata4).  
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Object().
	ContainsKey("name").
	ValueEqual("name","testing_user_name_4")   
	
	testingOperationAccountOneID := obj4.Value("id").String().Raw()	

	// making deposit
	postdata= map[string]interface{}{
		"operation": "deposit",
		"body": map[string]interface{} {
			"amount": 150,
		},
	}
    contentType = "application/json;charset=utf-8"

    e.POST("/transactions").  
	WithHeader("ContentType", contentType). 
	WithHeader("account_id", testingCustomerAccountThreeID).
    WithJSON(postdata).  
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Object().
	ContainsKey("id").
	ValueEqual("id",testingCustomerAccountThreeID).
	ContainsKey("balance").
	ValueEqual("balance",150)   

	// getting transaction id
	objs := e.GET("/transactions").
	WithQuery("account_id", testingCustomerAccountThreeID).  
	WithQuery("asc", 1).
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Array()

	objs.Element(0).Object().ValueEqual("operation", "deposit")
	fixTransactionID := objs.Element(0).Object().Value("id").String().Raw()

	e.POST("/transactions/" + fixTransactionID + "/undo").
	WithHeader("account_id", testingOperationAccountOneID).
    Expect().
    Status(http.StatusOK). 
	JSON().    
	Object().    
	ContainsKey("id").
	ValueEqual("id",testingCustomerAccountThreeID).
	ContainsKey("balance").
	ValueEqual("balance",0)   

	objs = e.GET("/transactions").
	WithQuery("account_id", testingOperationAccountOneID).  
	WithQuery("asc", 1).
    Expect().
    Status(http.StatusOK). 
    JSON().                 
	Array()

	objs.Element(0).Object().ValueEqual("notes", "fix transaction: " + fixTransactionID)	

	// clean up
	e.DELETE("/accounts/" + testingCustomerAccountThreeID).
	Expect().
	Status(http.StatusOK)
	
	e.DELETE("/accounts/" + testingCustomerAccountThreeID + "/transactions").
	Expect().
	Status(http.StatusOK)

	e.DELETE("/accounts/" + testingOperationAccountOneID).
	Expect().
	Status(http.StatusOK)
	
	e.DELETE("/accounts/" + testingOperationAccountOneID + "/transactions").
	Expect().
	Status(http.StatusOK)
}
