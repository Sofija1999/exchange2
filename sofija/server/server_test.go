package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bloxico/exchange-gateway/sofija/config"
	"github.com/Bloxico/exchange-gateway/sofija/server/auth"
	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	server Server
}

func (suite *ServerSuite) SetupTest() {
}

func (suite *ServerSuite) TearDownTest() {
}

func (suite *ServerSuite) SetupSuite() {

	cfg := config.ServerConfig{
		Port:   3001,
		Logger: nil,
	}

	suite.server = *NewServer(cfg, nil)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

// Validate the ping route is working, using basic JWT auth
func (suite *ServerSuite) TestPingRoute() {
	userEmail := "testy@email.com"
	jwtToken, err := auth.CreateJWT(userEmail, "id-123")
	if err != nil {
		suite.T().Fatal(err)
	}

	// then request their profile
	httpRequest, _ := http.NewRequest("GET", "/ping", nil)
	httpRequest.Header.Set("Authorization", "Bearer "+jwtToken)
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)
	responseRec := httptest.NewRecorder()

	suite.server.wsCont.ServeHTTP(responseRec, httpRequest)

	assert.Equal(suite.T(), 200, responseRec.Result().StatusCode)
}
