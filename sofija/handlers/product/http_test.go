package product

import (
	"egw-be/sofija/server/auth"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Bloxico/exchange-gateway/sofija/app"
	"github.com/Bloxico/exchange-gateway/sofija/core/usecases"
	"github.com/Bloxico/exchange-gateway/sofija/repo"
	"github.com/Bloxico/exchange-gateway/sofija/testutil"

	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testApp *app.App

type HttpSuite struct {
	suite.Suite
	productHttpSvc EgwProductHttpHandler
	wsContainer    *restful.Container
}

func (suite *HttpSuite) SetupTest() {
}

func (suite *HttpSuite) TearDownTest() {
	testutil.CleanUpTables(*testApp.DB)
}

func (suite *HttpSuite) SetupSuite() {

	testApp = testutil.InitTestApp()
	suite.wsContainer = restful.NewContainer()
	http.Handle("/", suite.wsContainer)

	realProductRep := repo.NewEgwProductRepository(testApp.DB)
	realProductSvc := usecases.NewEgwProductService(realProductRep)
	suite.productHttpSvc = *NewEgwProductHandler(realProductSvc, suite.wsContainer)

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(HttpSuite))
}

func (suite *HttpSuite) TestInsertProduct() {
	//prepare data for insert product
	postData := InsertRequestData{
		Name:             "borovnica",
		ShortDescription: "sitna borovnica",
		Description:      "sitna tamna borovnica iz sume",
		Price:            400,
	}

	// make request
	responseRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/product/insert", postData, nil)

	// validate response
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	var returnedProduct InsertResponseData

	fmt.Println(responseRec.Body.String())

	err := json.Unmarshal(responseRec.Body.Bytes(), &returnedProduct)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling product to json: %s", err)
	}

	assert.NotNil(suite.T(), returnedProduct.Product)
	assert.Equal(suite.T(), returnedProduct.Product.Name, postData.Name)
	assert.Equal(suite.T(), returnedProduct.Product.ShortDescription, postData.ShortDescription)
	assert.Equal(suite.T(), returnedProduct.Product.Description, postData.Description)
	assert.Equal(suite.T(), returnedProduct.Product.Price, postData.Price)
}

func (suite *HttpSuite) TestUpdateProduct() {
	// prepare insert data
	postData := InsertRequestData{
		Name:             "borovnica",
		ShortDescription: "sitna borovnica",
		Description:      "sitna tamna borovnica iz sume",
		Price:            400,
	}
	// make request
	responseRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/product/insert", postData, nil)

	// validate response
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code, "Error insert product")

	fmt.Println(responseRec.Body.String())

	var returnedProduct InsertResponseData
	err := json.Unmarshal(responseRec.Body.Bytes(), &returnedProduct)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling product data to json: %s", err)
	}

	id := returnedProduct.Product.ID

	realJWTToken, err := auth.CreateJWT("ana@gmail.com", "b221ca38-7b76-45cd-9c64-1d3c5a88220e")
	if err != nil {
		suite.T().Fatalf("Error creating JWT token: %s", err)
	}

	// make request to update a product's handle
	updateData := UpdateRequestData{
		Name:             "updated borovnica",
		ShortDescription: "updated short description",
		Description:      "updated description",
		Price:            500,
	}

	endpoint := fmt.Sprintf("/product/update/%s", id)
	responseRec2 := testutil.MakeRequest(*suite.wsContainer, "PUT", endpoint, updateData, &realJWTToken)

	assert.Equal(suite.T(), http.StatusOK, responseRec2.Code)

	var updatedProduct EgwProductModel
	err = json.Unmarshal(responseRec2.Body.Bytes(), &updatedProduct)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling product to json: %s", err)
	}

	assert.Equal(suite.T(), updateData.Name, updatedProduct.Name)
	assert.Equal(suite.T(), updateData.ShortDescription, updatedProduct.ShortDescription)
	assert.Equal(suite.T(), updateData.Description, updatedProduct.Description)
	assert.Equal(suite.T(), updateData.Price, updatedProduct.Price)

}

func (suite *HttpSuite) TestDeleteProduct() {

	postData := InsertRequestData{
		Name:             "borovnica",
		ShortDescription: "sitna borovnica",
		Description:      "sitna tamna borovnica iz sume",
		Price:            400,
	}

	responseRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/product/insert", postData, nil)

	assert.Equal(suite.T(), http.StatusOK, responseRec.Code, "Error insert product")

	var returnedProduct InsertResponseData
	err := json.Unmarshal(responseRec.Body.Bytes(), &returnedProduct)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling product data to json: %s", err)
	}

	id := returnedProduct.Product.ID

	realJWTToken, err := auth.CreateJWT("ana@gmail.com", "b221ca38-7b76-45cd-9c64-1d3c5a88220e")
	if err != nil {
		suite.T().Fatalf("Error creating JWT token: %s", err)
	}

	endpoint := fmt.Sprintf("/product/delete/%s", id)
	responseRec2 := testutil.MakeRequest(*suite.wsContainer, "DELETE", endpoint, nil, &realJWTToken)

	assert.Equal(suite.T(), http.StatusOK, responseRec2.Code, "Greška pri brisanju proizvoda")

	
}
