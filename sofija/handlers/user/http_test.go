package user

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Bloxico/exchange-gateway/sofija/app"
	"github.com/Bloxico/exchange-gateway/sofija/core/domain"
	"github.com/Bloxico/exchange-gateway/sofija/core/usecases"
	"github.com/Bloxico/exchange-gateway/sofija/repo"
	"github.com/Bloxico/exchange-gateway/sofija/testutil"
	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

var testApp *app.App

type HttpSuite struct {
	suite.Suite
	userHttpSvc EgwUserHttpHandler
	wsContainer *restful.Container
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

	realUserRep := repo.NewEgwUserRepository(testApp.DB)
	realUserSvc := usecases.NewEgwUserService(realUserRep)
	suite.userHttpSvc = *NewEgwUserHandler(realUserSvc, suite.wsContainer)

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(HttpSuite))
}

func (suite *HttpSuite) TestRegisterUser() {
	// prepare registration data
	postData := RegisterRequestData{
		Email:   "testy@email.com",
		Name:    "First name",
		Surname: "Last name",
	}
	// make request
	responseRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/user/register", postData, nil)

	// validate response
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	var returnedUser RegisterResponseData
	err := json.Unmarshal(responseRec.Body.Bytes(), &returnedUser)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling user profile to json: %s", err)
	}

	assert.NotNil(suite.T(), returnedUser.AuthToken)
	assert.NotNil(suite.T(), returnedUser.User)
	assert.Equal(suite.T(), returnedUser.User.Email, postData.Email)
	assert.Equal(suite.T(), returnedUser.User.Surname, postData.Surname)
	assert.Equal(suite.T(), returnedUser.User.Name, postData.Name)
}

func (suite *HttpSuite) TestLoginUser() {
	// register the user before sending login request
	userEmail := "testy@email.com"
	// register the user
	userPass := "password123"
	passHash, _ := bcrypt.GenerateFromPassword([]byte(userPass), 10)
	suite.userHttpSvc.userSvc.RegisterUser(context.TODO(), &domain.EgwUser{
		Email:        userEmail,
		PasswordHash: string(passHash),
	})

	// prepare login data
	postData := LoginRequestData{Email: userEmail, Password: userPass}
	responseRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/user/login", postData, nil)

	// validate response
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	var returnedUser LoginResponseData
	err := json.Unmarshal(responseRec.Body.Bytes(), &returnedUser)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling user profile to json: %s", err)
	}

	assert.NotNil(suite.T(), returnedUser.AuthToken)
	assert.NotNil(suite.T(), returnedUser.User)
	assert.Equal(suite.T(), returnedUser.User.Email, postData.Email)
}

func (suite *HttpSuite) TestInvalidLogin() {
	// register the user before sending login request
	userEmail := "testy@email.com"
	// register the user
	userPass := "password123"
	passHash, _ := bcrypt.GenerateFromPassword([]byte(userPass), 10)
	suite.userHttpSvc.userSvc.RegisterUser(context.TODO(), &domain.EgwUser{
		Email:        userEmail,
		PasswordHash: string(passHash),
	})

	// prepare login data
	postData := LoginRequestData{Email: userEmail, Password: "invalid password 123"}
	responseRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/user/login", postData, nil)

	// validate response
	assert.Equal(suite.T(), http.StatusForbidden, responseRec.Code)
}

func (suite *HttpSuite) TestUpdateUser() {
	// prepare registration data
	postData := RegisterRequestData{
		Email:   "testy@email.com",
		Name:    "First name",
		Surname: "Last name",
	}
	// make request
	responseRec := testutil.MakeRequest(*suite.wsContainer, "POST", "/user/register", postData, nil)

	// validate response
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code, "Error registering")
	var returnedUser RegisterResponseData
	err := json.Unmarshal(responseRec.Body.Bytes(), &returnedUser)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling user profile to json: %s", err)
	}
	token := returnedUser.AuthToken

	_ = token

	// make request to update a user's handle
	updateData := UpdateRequestData{
		Name:    "New name",
		Surname: "New surname",
	}
	responseRec2 := testutil.MakeRequest(*suite.wsContainer, "PUT", "/user", updateData, &token)
	assert.Equal(suite.T(), http.StatusOK, responseRec2.Code)

	var updatedUser EgwUserModel
	err = json.Unmarshal(responseRec2.Body.Bytes(), &updatedUser)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling user profile to json: %s", err)
	}

	assert.Equal(suite.T(), updateData.Name, updatedUser.Name)
	assert.Equal(suite.T(), updateData.Surname, updatedUser.Surname)
}
