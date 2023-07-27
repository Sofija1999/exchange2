package order

import (
	"bytes"
	"egw-be/sofija/core/usecases"
	"egw-be/sofija/testutil"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"egw-be/sofija/handlers/product"
	"egw-be/sofija/handlers/user"

	"github.com/Bloxico/exchange-gateway/sofija/repo"

	"github.com/Bloxico/exchange-gateway/sofija/app"

	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/suite"
)

var testApp *app.App

type HttpSuite struct {
	suite.Suite
	orderHttpSvc EgwOrderHttpHandler
	wsContainer  *restful.Container
}

func (suite *HttpSuite) SetupTest() {
	// Priprema podataka za registraciju korisnika
	postData := user.RegisterRequestData{
		Email:   "testy1@email.com",
		Name:    "First name",
		Surname: "Last name",
	}

	// Serijalizacija podataka u JSON format
	reqBody, err := json.Marshal(postData)
	if err != nil {
		suite.T().Fatalf("Error serializing user registration data to JSON: %s", err)
	}

	// Kreiranje HTTP zahteva
	req, err := http.NewRequest("POST", "/user/register", bytes.NewReader(reqBody))
	if err != nil {
		suite.T().Fatalf("Error creating HTTP request: %s", err)
	}

	// Postavljanje zaglavlja zahteva (opciono)
	req.Header.Set("Content-Type", "application/json")

	// Izvršavanje HTTP zahteva
	responseRec := httptest.NewRecorder()
	suite.wsContainer.ServeHTTP(responseRec, req)

	// Provera odgovora
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)

	// Deserijalizacija odgovora u strukturu RegisterResponseData
	var returnedUser user.RegisterResponseData
	err = json.Unmarshal(responseRec.Body.Bytes(), &returnedUser)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling user profile to JSON: %s", err)
	}

	// Provera da li je odgovor ispravan
	assert.NotNil(suite.T(), returnedUser.AuthToken)
	assert.NotNil(suite.T(), returnedUser.User)
	assert.Equal(suite.T(), returnedUser.User.Email, postData.Email)
	assert.Equal(suite.T(), returnedUser.User.Surname, postData.Surname)
	assert.Equal(suite.T(), returnedUser.User.Name, postData.Name)

	//prepare data for insert product
	postData1 := product.InsertRequestData{
		Name:             "borovnica",
		ShortDescription: "sitna borovnica",
		Description:      "sitna tamna borovnica iz sume",
		Price:            400,
	}

	// make request
	responseRec1 := testutil.MakeRequest(*suite.wsContainer, "POST", "/product/insert", postData1, nil)

	// validate response
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	var returnedProduct product.InsertResponseData

	fmt.Println(responseRec1.Body.String())

	err1 := json.Unmarshal(responseRec1.Body.Bytes(), &returnedProduct)
	if err1 != nil {
		suite.T().Fatalf("Error unmarshalling product to json: %s", err1)
	}

	assert.NotNil(suite.T(), returnedProduct.Product)
	assert.Equal(suite.T(), returnedProduct.Product.Name, postData1.Name)
	assert.Equal(suite.T(), returnedProduct.Product.ShortDescription, postData1.ShortDescription)
	assert.Equal(suite.T(), returnedProduct.Product.Description, postData1.Description)
	assert.Equal(suite.T(), returnedProduct.Product.Price, postData1.Price)
}

func (suite *HttpSuite) TearDownTest() {
	testutil.CleanUpTables(*testApp.DB)
}

func (suite *HttpSuite) SetupSuite() {

	testApp = testutil.InitTestApp()
	suite.wsContainer = restful.NewContainer()
	http.Handle("/", suite.wsContainer)

	realOrderRep := repo.NewEgwOrderRepository(testApp.DB)
	realOrderSvc := usecases.NewEgwOrderService(realOrderRep)
	suite.orderHttpSvc = *NewEgwOrderHandler(realOrderSvc, suite.wsContainer)

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(HttpSuite))
}

func (suite *HttpSuite) TestInsertOrder() {
	//data for insert order
	postData := InsertRequestData{
		UserID: "e23df3a8-4c06-4652-b432-0e2ee514575c",
		Status: "CREATED",
		Items: []*InsertOrderItemRequest{
			{
				ProductID:   "6e692fcc-202b-487d-b243-b52f0031b338",
				ProductName: "jabuka",
				Quantity:    10,
			},
			{
				ProductID:   "f01b7274-f2cd-4970-b917-572af43600c0",
				ProductName: "kupine",
				Quantity:    30,
			},
		},
	}

	//make request
	responseRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/order/insert", postData, nil)
	fmt.Println("greskica1")

	// validate response
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	var returnedOrder InsertRequestData

	fmt.Println(responseRec.Body.String())

	err := json.Unmarshal(responseRec.Body.Bytes(), &returnedOrder)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling order to json: %s", err)
	}

	//check if order is nil
	assert.NotNil(suite.T(), returnedOrder)

	//check userID and status of order
	assert.Equal(suite.T(), returnedOrder.UserID, postData.UserID)
	assert.Equal(suite.T(), returnedOrder.Status, postData.Status)

	//check if order items is nil and len of items
	assert.NotNil(suite.T(), returnedOrder.Items)
	assert.Len(suite.T(), returnedOrder.Items, len(postData.Items))

	//assertions for each order item in the returned order
	for i, item := range returnedOrder.Items {
		assert.Equal(suite.T(), item.ProductID, postData.Items[i].ProductID)
		assert.Equal(suite.T(), item.ProductName, postData.Items[i].ProductName)
		assert.Equal(suite.T(), item.Quantity, postData.Items[i].Quantity)
	}

}
