package routes

import (
	"context"
	"github.com/Dissociable/Couploan/config"
	"github.com/Dissociable/Couploan/ent"
	"github.com/Dissociable/Couploan/pkg/services"
	"github.com/Dissociable/Couploan/pkg/tests"
	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var (
	Container *services.Container
	TestUser  *ent.User
)

func TestMain(m *testing.M) {
	// Set the environment to test
	config.SwitchEnvironment(config.EnvTest)

	// Initialize the container
	Container = services.NewContainer()

	var err error
	TestUser, err = tests.CreateUser(Container.ORM)
	if err != nil {
		panic(errors.Wrap(err, "failed to create test user"))
	}
	BuildRouter(Container)

	// Run the tests
	m.Run()
}

func TestApiCheckFraud(t *testing.T) {
	// User with no balance request
	var err error
	// Create a new fiber instance
	key := TestUser.Key
	headers := http.Header{
		"Authorization": {"Bearer " + *key},
	}
	beforeRequestBalance := TestUser.Balance
	resp, err := tests.NewContextTestWithHeaders(
		Container.Web, "/api/v1/check_fraud?ip=1.1.1.1", headers, func(ctx fiber.Ctx) {

		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)

	TestUser, err = Container.ORM.User.Get(context.Background(), TestUser.ID)
	require.NoError(t, err)

	afterRequestBalance := TestUser.Balance
	assert.Equal(t, beforeRequestBalance, afterRequestBalance)

	TestUser, err = TestUser.Update().SetBalance(100).Save(context.Background())
	require.NoError(t, err)
	require.NotNil(t, TestUser)

	// // User with balance request
	// beforeRequestBalance = TestUser.Balance
	// resp, err = tests.NewContextTestWithHeaders(
	// 	Container.Web, "/api/v1/check_fraud?ip=1.1.1.1", headers, func(ctx fiber.Ctx) {
	//
	// 	},
	// )
	// require.NoError(t, err)
	// require.NotNil(t, resp)
	//
	// TestUser, err = Container.ORM.User.Get(context.Background(), TestUser.ID)
	// require.NoError(t, err)
	//
	// afterRequestBalance = TestUser.Balance
	// assert.Equal(t, resp.StatusCode, http.StatusOK)
	// assert.Equal(t, beforeRequestBalance-Container.Config.Pricing.PricePerCheck, afterRequestBalance)
}
