package usecases

import (
	"testing"

	"github.com/Bloxico/exchange-gateway/sofija/repo"
	"github.com/Bloxico/exchange-gateway/sofija/testutil"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type EgwProductSuite struct {
	suite.Suite
	productRep *repo.EgwProductRepository
	productSvc *EgwProductService
}

func (suite *EgwProductSuite) SetupTest() {
}

func (suite *EgwProductSuite) TearDownTest() {
	// todo: clear the DB
}

func (suite *EgwProductSuite) SetupSuite() {

	app := testutil.InitTestApp()
	suite.productRep = repo.NewEgwProductRepository(app.DB)
	suite.productSvc = NewEgwProductService(suite.productRep)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(EgwProductSuite))
}
