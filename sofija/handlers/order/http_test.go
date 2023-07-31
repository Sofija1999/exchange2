package order

import (
	"egw-be/sofija/core/usecases"
	"egw-be/sofija/handlers/product"
	"egw-be/sofija/handlers/user"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Bloxico/exchange-gateway/sofija/app"
	"github.com/Bloxico/exchange-gateway/sofija/repo"
	"github.com/Bloxico/exchange-gateway/sofija/testutil"
	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testApp *app.App

type HttpSuite struct {
	suite.Suite
	orderHttpSvc   EgwOrderHttpHandler
	userHttpSvc    user.EgwUserHttpHandler
	productHttpSvc product.EgwProductHttpHandler
	wsContainer    *restful.Container
}

func (suite *HttpSuite) SetupTest() {
}

func (suite *HttpSuite) TearDownTest() {
	fmt.Println("TearDownTest - cleaning up tables")
	testutil.CleanUpTables(*testApp.DB)
}

func (suite *HttpSuite) SetupSuite() {

	testApp = testutil.InitTestApp()
	suite.wsContainer = restful.NewContainer()
	http.Handle("/", suite.wsContainer)

	realOrderRep := repo.NewEgwOrderRepository(testApp.DB)
	realOrderSvc := usecases.NewEgwOrderService(realOrderRep)
	suite.orderHttpSvc = *NewEgwOrderHandler(realOrderSvc, suite.wsContainer)
	//orderHandler := NewEgwOrderHandler(realOrderSvc, suite.wsContainer)
	//suite.orderHttpSvc = orderHandler

	realUserRep := repo.NewEgwUserRepository(testApp.DB)
	realUserSvc := usecases.NewEgwUserService(realUserRep)
	suite.userHttpSvc = *user.NewEgwUserHandler(realUserSvc, suite.wsContainer)

	realProductRep := repo.NewEgwProductRepository(testApp.DB)
	realProductSvc := usecases.NewEgwProductService(realProductRep)
	suite.productHttpSvc = *product.NewEgwProductHandler(realProductSvc, suite.wsContainer)

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(HttpSuite))
}

func (suite *HttpSuite) TestInsertOrder() {

	// Prepare user data
	postUserData := user.RegisterRequestData{
		Email:    "testy87@email.com",
		Name:     "First name",
		Surname:  "Last name",
		Password: "testpassword",
	}
	// Make request to register user
	responseUserRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/user/register", postUserData, nil)

	fmt.Println(responseUserRec.Body.String())
	// Validate user registration response
	assert.Equal(suite.T(), http.StatusOK, responseUserRec.Code)
	var returnedUser user.RegisterResponseData
	err := json.Unmarshal(responseUserRec.Body.Bytes(), &returnedUser)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling user profile to json: %s", err)
	}

	// Check that user is not nil and the registration data is correct
	assert.NotNil(suite.T(), returnedUser.AuthToken)
	assert.NotNil(suite.T(), returnedUser.User)
	assert.Equal(suite.T(), returnedUser.User.Email, postUserData.Email)
	assert.Equal(suite.T(), returnedUser.User.Surname, postUserData.Surname)
	assert.Equal(suite.T(), returnedUser.User.Name, postUserData.Name)

	// Prepare product data
	postProductData := product.InsertRequestData{
		Name:             "borovnica",
		ShortDescription: "sitna borovnica",
		Description:      "sitna tamna borovnica iz sume",
		Price:            400,
	}
	// Make request to insert product
	responseProductRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/product/insert", postProductData, nil)

	// Validate product insertion response
	assert.Equal(suite.T(), http.StatusOK, responseProductRec.Code)
	var returnedProduct product.InsertResponseData
	err = json.Unmarshal(responseProductRec.Body.Bytes(), &returnedProduct)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling product to json: %s", err)
	}

	// Check that product is not nil and the insertion data is correct
	assert.NotNil(suite.T(), returnedProduct.Product)
	assert.Equal(suite.T(), returnedProduct.Product.Name, postProductData.Name)
	assert.Equal(suite.T(), returnedProduct.Product.ShortDescription, postProductData.ShortDescription)
	assert.Equal(suite.T(), returnedProduct.Product.Description, postProductData.Description)
	assert.Equal(suite.T(), returnedProduct.Product.Price, postProductData.Price)

	// Prepare order data
	postOrderData := InsertRequestData{
		UserID: returnedUser.User.ID, // Use the user ID from the registered user
		Status: "CREATED",
		Items: []*InsertOrderItemRequest{
			{
				ProductID:   returnedProduct.Product.ID, // Use the product ID from the inserted product
				ProductName: returnedProduct.Product.Name,
				Quantity:    10,
			},
			// Add more items if needed
		},
	}
	// Make request to insert order
	responseOrderRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/order/insert", postOrderData, nil)

	// Validate order insertion response
	assert.Equal(suite.T(), http.StatusOK, responseOrderRec.Code)
	var returnedOrder InsertResponseData
	fmt.Println(returnedOrder.ID)
	err = json.Unmarshal(responseOrderRec.Body.Bytes(), &returnedOrder)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling order to json: %s", err)
	}
}

func (suite *HttpSuite) TestDeleteOrder() {

	// Prepare user data
	postUserData := user.RegisterRequestData{
		Email:    "testy55@email.com",
		Name:     "First name",
		Surname:  "Last name",
		Password: "testpassword",
	}
	// Make request to register user
	responseUserRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/user/register", postUserData, nil)

	fmt.Println(responseUserRec.Body.String())
	// Validate user registration response
	assert.Equal(suite.T(), http.StatusOK, responseUserRec.Code)
	var returnedUser user.RegisterResponseData
	err := json.Unmarshal(responseUserRec.Body.Bytes(), &returnedUser)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling user profile to json: %s", err)
	}

	// Check that user is not nil and the registration data is correct
	assert.NotNil(suite.T(), returnedUser.AuthToken)
	assert.NotNil(suite.T(), returnedUser.User)
	assert.Equal(suite.T(), returnedUser.User.Email, postUserData.Email)
	assert.Equal(suite.T(), returnedUser.User.Surname, postUserData.Surname)
	assert.Equal(suite.T(), returnedUser.User.Name, postUserData.Name)

	// Prepare product data
	postProductData := product.InsertRequestData{
		Name:             "kupina",
		ShortDescription: "sitna kupina",
		Description:      "sitna tamna kupina iz sume",
		Price:            330,
	}
	// Make request to insert product
	responseProductRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/product/insert", postProductData, nil)

	// Validate product insertion response
	assert.Equal(suite.T(), http.StatusOK, responseProductRec.Code)
	var returnedProduct product.InsertResponseData
	err = json.Unmarshal(responseProductRec.Body.Bytes(), &returnedProduct)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling product to json: %s", err)
	}

	// Check that product is not nil and the insertion data is correct
	assert.NotNil(suite.T(), returnedProduct.Product)
	assert.Equal(suite.T(), returnedProduct.Product.Name, postProductData.Name)
	assert.Equal(suite.T(), returnedProduct.Product.ShortDescription, postProductData.ShortDescription)
	assert.Equal(suite.T(), returnedProduct.Product.Description, postProductData.Description)
	assert.Equal(suite.T(), returnedProduct.Product.Price, postProductData.Price)

	// Prepare order data
	postOrderData := InsertRequestData{
		UserID: returnedUser.User.ID, // Use the user ID from the registered user
		Status: "CREATED",
		Items: []*InsertOrderItemRequest{
			{
				ProductID:   returnedProduct.Product.ID, // Use the product ID from the inserted product
				ProductName: returnedProduct.Product.Name,
				Quantity:    10,
			},
			// Add more items if needed
		},
	}
	// Make request to insert order
	responseOrderRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/order/insert", postOrderData, nil)

	// Validate order insertion response
	assert.Equal(suite.T(), http.StatusOK, responseOrderRec.Code)
	var returnedOrder InsertResponseData
	fmt.Println(returnedOrder.ID)
	err = json.Unmarshal(responseOrderRec.Body.Bytes(), &returnedOrder)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling order to json: %s", err)
	}

	id := returnedOrder.ID

	endpoint := fmt.Sprintf("/order/delete/%s", id)
	responseRec2 := testutil.MakeRequest(*suite.wsContainer, "DELETE", endpoint, nil, nil)

	assert.Equal(suite.T(), http.StatusOK, responseRec2.Code, "Error while delete order")
}

func (suite *HttpSuite) TestUpdateOrder() {
	// Prepare user data
	postUserData := user.RegisterRequestData{
		Email:    "testy1689@email.com",
		Name:     "First name",
		Surname:  "Last name",
		Password: "testpassword",
	}
	// Make request to register user
	responseUserRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/user/register", postUserData, nil)

	fmt.Println(responseUserRec.Body.String())
	// Validate user registration response
	assert.Equal(suite.T(), http.StatusOK, responseUserRec.Code)
	var returnedUser user.RegisterResponseData
	err := json.Unmarshal(responseUserRec.Body.Bytes(), &returnedUser)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling user profile to json: %s", err)
	}

	// Check that user is not nil and the registration data is correct
	assert.NotNil(suite.T(), returnedUser.AuthToken)
	assert.NotNil(suite.T(), returnedUser.User)
	assert.Equal(suite.T(), returnedUser.User.Email, postUserData.Email)
	assert.Equal(suite.T(), returnedUser.User.Surname, postUserData.Surname)
	assert.Equal(suite.T(), returnedUser.User.Name, postUserData.Name)

	// Prepare product data
	postProductData := product.InsertRequestData{
		Name:             "malina",
		ShortDescription: "sitna malina",
		Description:      "sitna roze malina iz sume",
		Price:            850,
	}
	// Make request to insert product
	responseProductRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/product/insert", postProductData, nil)

	// Validate product insertion response
	assert.Equal(suite.T(), http.StatusOK, responseProductRec.Code)
	var returnedProduct product.InsertResponseData
	err = json.Unmarshal(responseProductRec.Body.Bytes(), &returnedProduct)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling product to json: %s", err)
	}

	// Check that product is not nil and the insertion data is correct
	assert.NotNil(suite.T(), returnedProduct.Product)
	assert.Equal(suite.T(), returnedProduct.Product.Name, postProductData.Name)
	assert.Equal(suite.T(), returnedProduct.Product.ShortDescription, postProductData.ShortDescription)
	assert.Equal(suite.T(), returnedProduct.Product.Description, postProductData.Description)
	assert.Equal(suite.T(), returnedProduct.Product.Price, postProductData.Price)

	// Prepare order data
	postOrderData := InsertRequestData{
		UserID: returnedUser.User.ID, // Use the user ID from the registered user
		Status: "CREATED",
		Items: []*InsertOrderItemRequest{
			{
				ProductID:   returnedProduct.Product.ID, // Use the product ID from the inserted product
				ProductName: returnedProduct.Product.Name,
				Quantity:    10,
			},
			// Add more items if needed
		},
	}
	// Make request to insert order
	responseOrderRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/order/insert", postOrderData, nil)

	// Validate order insertion response
	assert.Equal(suite.T(), http.StatusOK, responseOrderRec.Code)
	var returnedOrder InsertResponseData
	fmt.Println(returnedOrder.ID)
	err = json.Unmarshal(responseOrderRec.Body.Bytes(), &returnedOrder)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling order to json: %s", err)
	}

	id := returnedOrder.ID

	// make request to update a order's handle
	updateData := UpdateRequestData{
		Status: "COMPLETED",
	}

	endpoint := fmt.Sprintf("/order/update/%s", id)
	responseRec2 := testutil.MakeRequest(*suite.wsContainer, "PUT", endpoint, updateData, nil)

	assert.Equal(suite.T(), http.StatusOK, responseRec2.Code)

	var updatedOrder EgwOrderModel
	err = json.Unmarshal(responseRec2.Body.Bytes(), &updatedOrder)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling product to json: %s", err)
	}

	assert.Equal(suite.T(), updateData.Status, updatedOrder.Status)

}
