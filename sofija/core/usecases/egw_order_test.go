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
type EgwOrderSuite struct {
	suite.Suite
	orderRep *repo.EgwOrderRepository
	orderSvc *EgwOrderService
}

func (suite *EgwOrderSuite) SetupTest() {
}

func (suite *EgwOrderSuite) TearDownTest() {
	// todo: clear the DB
}

func (suite *EgwOrderSuite) SetupSuite() {

	app := testutil.InitTestApp()
	suite.orderRep = repo.NewEgwOrderRepository(app.DB)
	suite.orderSvc = NewEgwOrderService(suite.orderRep)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(EgwOrderSuite))
}
