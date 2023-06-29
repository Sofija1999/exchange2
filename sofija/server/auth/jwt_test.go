package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bloxico/exchange-gateway/sofija/server/params"
	"github.com/emicklei/go-restful/v3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	suite.Suite
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

func (suite *AuthSuite) TestCreateToken() {
	userEmail := "testy@email.com"
	userId := "abcd-1234"
	// create the token
	jwtToken, err := CreateJWT(userEmail, userId)
	if err != nil {
		suite.T().Fatal(err)
	}

	// wrap and decode it
	header := WrapJWTHeader(jwtToken)
	isValid := IsJWTHeaderValid(header)

	assert.True(suite.T(), isValid)

	res, err := GetJWTClaims(jwtToken)
	if err != nil {
		suite.T().Fatal(err)
	}

	// check the email is embedded into the JWT
	assert.Equal(suite.T(), res["email"], userEmail)
	assert.Equal(suite.T(), res["id"], userId)
}

// Test that the JWT package attaches desired fields to the request context
func (suite *AuthSuite) TestJWTAttachesContext() {

	var routeParam string

	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON)
	restful.Add(ws)

	// create a route used in this test only
	ws.Route(ws.GET("/jwt/test").Filter(AuthJWT).To(func(r1 *restful.Request, r2 *restful.Response) {
		// pull the value from the request context (should be set by JWT package)
		routeParam, _ = params.StringFrom(r1.Request, USER_EMAIL_CTX_KEY)
	}))

	userEmail := "testy@email.com"
	jwtToken, err := CreateJWT(userEmail, "abcd-123")
	if err != nil {
		suite.T().Fatal(err)
	}

	// make a request to the specified route
	httpRequest, _ := http.NewRequest("GET", "/jwt/test", nil)
	httpRequest.Header.Set("Authorization", "Bearer "+jwtToken)
	responseRec := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(responseRec, httpRequest)

	assert.Equal(suite.T(), routeParam, userEmail)
}
