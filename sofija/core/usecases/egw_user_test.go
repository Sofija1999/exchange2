package usecases

import (
	"context"
	"testing"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"
	"github.com/Bloxico/exchange-gateway/sofija/repo"
	"github.com/Bloxico/exchange-gateway/sofija/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type EgwUserSuite struct {
	suite.Suite
	userRep *repo.EgwUserRepository
	userSvc *EgwUserService
}

func (suite *EgwUserSuite) SetupTest() {
}

func (suite *EgwUserSuite) TearDownTest() {
	// todo: clear the DB
}

func (suite *EgwUserSuite) SetupSuite() {

	app := testutil.InitTestApp()
	suite.userRep = repo.NewEgwUserRepository(app.DB)
	suite.userSvc = NewEgwUserService(suite.userRep)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(EgwUserSuite))
}

func (suite *EgwUserSuite) TestUserRegistration() {

	userEmail := "email1@provider.com"
	err := suite.userSvc.RegisterUser(context.TODO(), &domain.EgwUser{Email: userEmail})

	if err != nil {
		suite.T().Fatal(err)
	}

	user, err := suite.userSvc.FindByEmail(context.TODO(), userEmail)
	if err != nil {
		suite.T().Fatal(err)
	}

	assert.Equal(suite.T(), user.Email, userEmail)
}

func (suite *EgwUserSuite) TestCannotRegisterWithExistingEmail() {
	userEmail := "email2@provider.com"
	err := suite.userSvc.RegisterUser(context.TODO(), &domain.EgwUser{Email: userEmail})

	if err != nil {
		suite.T().Fatal(err)
	}

	err = suite.userSvc.RegisterUser(context.TODO(), &domain.EgwUser{Email: userEmail})

	assert.ErrorIs(suite.T(), err, repo.ErrDuplicateEmail)
}
