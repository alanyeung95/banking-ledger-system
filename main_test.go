package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
)

var testurl string = os.Getenv("API_TEST_DOMAIN")

func TestRootAPI(t *testing.T) {
	e := httpexpect.New(t, testurl)
	e.GET("/").
		Expect().
		Status(http.StatusOK)
}

func TestDepositAndWithdrawAndGetBalance(t *testing.T) {
	e := httpexpect.New(t, testurl)

	// create account
	postdata := map[string]interface{}{
		"name":     "testing_user_name_1",
		"password": "xxxxx",
	}
	contentType := "application/json;charset=utf-8"

	obj := e.POST("/accounts/create").
		WithHeader("ContentType", contentType).
		WithJSON(postdata).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("name").
		ValueEqual("name", "testing_user_name_1")
	testingAccountOneID := obj.Value("id").String().Raw()

	// test get balance
	e.GET("/accounts/"+testingAccountOneID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("id").
		ValueEqual("id", testingAccountOneID).
		ContainsKey("balance").
		ValueEqual("balance", 0)

		// test deposit
	postdata = map[string]interface{}{
		"operation": "deposit",
		"body": map[string]interface{}{
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
		ValueEqual("id", testingAccountOneID).
		ContainsKey("balance").
		ValueEqual("balance", 200)

	// test withdraw
	postdata = map[string]interface{}{
		"operation": "withdraw",
		"body": map[string]interface{}{
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
		ValueEqual("id", testingAccountOneID).
		ContainsKey("balance").
		ValueEqual("balance", 100)

	// clean up
	e.DELETE("/accounts/" + testingAccountOneID).
		Expect().
		Status(http.StatusOK)

	e.DELETE("/accounts/" + testingAccountOneID + "/transactions").
		Expect().
		Status(http.StatusOK)
}

func TestTransfer(t *testing.T) {
	e := httpexpect.New(t, testurl)
	contentType := "application/json;charset=utf-8"

	// creating first user
	postdata1 := map[string]interface{}{
		"name":     "testing_user_name_1",
		"password": "xxxxx",
	}

	obj1 := e.POST("/accounts/create").
		WithHeader("ContentType", contentType).
		WithJSON(postdata1).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("name").
		ValueEqual("name", "testing_user_name_1")

	testingAccountOneID := obj1.Value("id").String().Raw()

	// creating second user
	postdata2 := map[string]interface{}{
		"name":     "testing_user_name_2",
		"password": "xxxxx",
	}

	obj2 := e.POST("/accounts/create").
		WithHeader("ContentType", contentType).
		WithJSON(postdata2).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("name").
		ValueEqual("name", "testing_user_name_2")

	testingAccountTwoID := obj2.Value("id").String().Raw()

	// making deposit
	postdata := map[string]interface{}{
		"operation": "deposit",
		"body": map[string]interface{}{
			"amount": 100,
		},
	}

	e.POST("/transactions").
		WithHeader("ContentType", contentType).
		WithHeader("account_id", testingAccountOneID).
		WithJSON(postdata).
		Expect()

		// making transfer
	postdata = map[string]interface{}{
		"operation": "transfer",
		"body": map[string]interface{}{
			"amount": 100,
			"to":     testingAccountTwoID,
		},
	}
	contentType = "application/json;charset=utf-8"

	e.POST("/transactions").
		WithHeader("ContentType", contentType).
		WithHeader("account_id", testingAccountOneID).
		WithJSON(postdata).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("id").
		ValueEqual("id", testingAccountOneID).
		ContainsKey("balance").
		ValueEqual("balance", 0)

	e.GET("/accounts/"+testingAccountTwoID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("id").
		ValueEqual("id", testingAccountTwoID).
		ContainsKey("balance").
		ValueEqual("balance", 100)

	// clean up
	e.DELETE("/accounts/" + testingAccountOneID).
		Expect().
		Status(http.StatusOK)

	e.DELETE("/accounts/" + testingAccountOneID + "/transactions").
		Expect().
		Status(http.StatusOK)

	e.DELETE("/accounts/" + testingAccountTwoID).
		Expect().
		Status(http.StatusOK)

	e.DELETE("/accounts/" + testingAccountTwoID + "/transactions").
		Expect().
		Status(http.StatusOK)
}

func TestViewCustomerTransaction(t *testing.T) {
	e := httpexpect.New(t, testurl)
	contentType := "application/json;charset=utf-8"

	// creating  user
	postdata1 := map[string]interface{}{
		"name":     "testing_user_name_1",
		"password": "xxxxx",
	}

	obj1 := e.POST("/accounts/create").
		WithHeader("ContentType", contentType).
		WithJSON(postdata1).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("name").
		ValueEqual("name", "testing_user_name_1")

	testingAccountID := obj1.Value("id").String().Raw()

	// making deposit
	postdata := map[string]interface{}{
		"operation": "deposit",
		"body": map[string]interface{}{
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
	e := httpexpect.New(t, testurl)

	// creating customer account
	postdata1 := map[string]interface{}{
		"name":     "testing_user_name_1",
		"password": "xxxxx",
	}
	contentType := "application/json;charset=utf-8"

	obj1 := e.POST("/accounts/create").
		WithHeader("ContentType", contentType).
		WithJSON(postdata1).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("name").
		ValueEqual("name", "testing_user_name_1")

	testingCustomerAccountOneID := obj1.Value("id").String().Raw()

	// creating operation team account
	postdataAdmin := map[string]interface{}{
		"name":     "testing_admin_name_1",
		"password": "xxxxx",
	}

	adminObj1 := e.POST("/accounts/create-admin").
		WithHeader("ContentType", contentType).
		WithJSON(postdataAdmin).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("name").
		ValueEqual("name", "testing_admin_name_1")

	testingOperationAccountOneID := adminObj1.Value("id").String().Raw()

	// making deposit
	postdata2 := map[string]interface{}{
		"operation": "deposit",
		"body": map[string]interface{}{
			"amount": 150,
		},
	}
	contentType = "application/json;charset=utf-8"

	e.POST("/transactions").
		WithHeader("ContentType", contentType).
		WithHeader("account_id", testingCustomerAccountOneID).
		WithJSON(postdata2).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("id").
		ValueEqual("id", testingCustomerAccountOneID).
		ContainsKey("balance").
		ValueEqual("balance", 150)

	// getting transaction id
	objs := e.GET("/transactions").
		WithQuery("account_id", testingCustomerAccountOneID).
		WithQuery("asc", 1).
		Expect().
		Status(http.StatusOK).
		JSON().
		Array()

	objs.Element(0).Object().ValueEqual("operation", "deposit")
	fixTransactionID := objs.Element(0).Object().Value("id").String().Raw()

	e.POST("/transactions/"+fixTransactionID+"/undo").
		WithHeader("account_id", testingOperationAccountOneID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("id").
		ValueEqual("id", testingCustomerAccountOneID).
		ContainsKey("balance").
		ValueEqual("balance", 0)

	objs = e.GET("/transactions").
		WithQuery("account_id", testingOperationAccountOneID).
		WithQuery("asc", 1).
		Expect().
		Status(http.StatusOK).
		JSON().
		Array()

	objs.Element(0).Object().ValueEqual("notes", "fix transaction: "+fixTransactionID)

	// clean up
	e.DELETE("/accounts/" + testingCustomerAccountOneID).
		Expect().
		Status(http.StatusOK)

	e.DELETE("/accounts/" + testingCustomerAccountOneID + "/transactions").
		Expect().
		Status(http.StatusOK)

	e.DELETE("/accounts/" + testingOperationAccountOneID).
		Expect().
		Status(http.StatusOK)

	e.DELETE("/accounts/" + testingOperationAccountOneID + "/transactions").
		Expect().
		Status(http.StatusOK)
}
