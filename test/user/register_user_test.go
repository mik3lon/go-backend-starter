package user_test

import (
	"bytes"
	"go-boilerplate/test"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go-boilerplate/internal/pkg/infrastructure/kernel"
)

// Define the test suite struct, embedding suite.Suite
type UserTestSuite struct {
	suite.Suite
	k *kernel.Kernel
}

// SetupSuite runs before the suite starts
func (suite *UserTestSuite) SetupSuite() {
}

// TeardownSuite runs after the suite finishes
func (suite *UserTestSuite) TearDownSuite() {
	// Any suite-level teardown logic
}

// SetupTest runs before each test
func (suite *UserTestSuite) SetupTest() {
}

// Test for successful user registration
func (suite *UserTestSuite) TestRegisterUserSuccess() {
	r := test.DoHttpRequest(http.MethodPost,
		"/users",
		bytes.NewBufferString(`{
		"id": "01F8MECHZX3TBDSZ7XRADM79XV",
		"username": "john_doe",
		"email": "john@example.com",
		"password": "password123",
		"name": "John",
		"surname": "Doe"
	}`),
	)

	// Use assert from testify
	assert.Equal(suite.T(), http.StatusNoContent, r.Code)
	assert.Equal(suite.T(), "", r.Body.String())
}

// This runs the test suite
func TestUserTestSuite(t *testing.T) {
	test.LoadTestEnv()
	suite.Run(t, new(UserTestSuite))
}
