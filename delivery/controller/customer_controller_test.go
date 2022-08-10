package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"enigmacamp.com/golatihanlagi/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyCustomer = []model.Customer{
	{
		Id:      "C001",
		Nama:    "Dummy Name 1",
		Address: "Dummy Address 1",
	},
	{
		Id:      "C002",
		Nama:    "Dummy Name 2",
		Address: "Dummy Address 2",
	},
}

type CustomerUsecaseMock struct {
	mock.Mock
}

type CustomerControllerTestSuite struct {
	suite.Suite
	mockRouter  *gin.Engine
	mockUsecase *CustomerUsecaseMock
}

func (suite *CustomerControllerTestSuite) SetupTest() {
	suite.mockRouter = gin.Default()
	suite.mockUsecase = new(CustomerUsecaseMock)
}

func (r *CustomerUsecaseMock) RegisterCustomer(customer model.Customer) error {
	args := r.Called(customer)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return nil
}
func (r *CustomerUsecaseMock) FindCustomerById(id string) (model.Customer, error) {
	args := r.Called(id)
	if args.Get(1) != nil {
		return model.Customer{}, args.Error(1)
	}
	return args.Get(0).(model.Customer), nil
}
func (r *CustomerUsecaseMock) GetAllCustomer() ([]model.Customer, error) {
	args := r.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Customer), nil
}

func (suite *CustomerControllerTestSuite) TestGetAllCustomerApi_Success() {
	suite.mockUsecase.On("GetAllCustomer").Return(dummyCustomer, nil)
	NewCustomerController(suite.mockRouter, suite.mockUsecase)
	r := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/customer", nil)
	suite.mockRouter.ServeHTTP(r, request)

	var actualCustomers []model.Customer
	response := r.Body.String()
	json.Unmarshal([]byte(response), &actualCustomers)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	assert.Equal(suite.T(), 2, len(actualCustomers))
	assert.Equal(suite.T(), dummyCustomer[0].Nama, actualCustomers[0].Nama)
	assert.Nil(suite.T(), err)
}

func (suite *CustomerControllerTestSuite) TestGetAllCustomerApi_Failed() {
	suite.mockUsecase.On("GetAllCustomer").Return(nil, errors.New("failed"))
	NewCustomerController(suite.mockRouter, suite.mockUsecase)
	r := httptest.NewRecorder()

	request, _ := http.NewRequest(http.MethodGet, "/customer", nil)
	suite.mockRouter.ServeHTTP(r, request)

	var errorResponse struct{ Err string }
	response := r.Body.String()
	json.Unmarshal([]byte(response), &errorResponse)
	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
	assert.Equal(suite.T(), "failed", errorResponse.Err)
}

func (suite *CustomerControllerTestSuite) TestRegisterCustomerApi_Success() {
	dummyCustomer := dummyCustomer[0]
	suite.mockUsecase.On("RegisterCustomer", dummyCustomer).Return(nil)
	NewCustomerController(suite.mockRouter, suite.mockUsecase)
	r := httptest.NewRecorder()

	requestBody, _ := json.Marshal(dummyCustomer)
	request, _ := http.NewRequest(http.MethodPost, "/customer", bytes.NewBuffer(requestBody))
	suite.mockRouter.ServeHTTP(r, request)

	var actualCustomers model.Customer
	response := r.Body.String()
	json.Unmarshal([]byte(response), &actualCustomers)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	assert.Equal(suite.T(), dummyCustomer.Nama, actualCustomers.Nama)
}

func (suite *CustomerControllerTestSuite) TestRegisterCustomerApi_FailedBinding() {
	// suite.mockUsecase.On("RegisterCustomer").Return(dummyCustomer, nil)
	NewCustomerController(suite.mockRouter, suite.mockUsecase)
	r := httptest.NewRecorder()

	request, _ := http.NewRequest(http.MethodPost, "/customer", nil)
	suite.mockRouter.ServeHTTP(r, request)

	var errorResponse struct{ Err string }
	response := r.Body.String()
	json.Unmarshal([]byte(response), &errorResponse)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.NotEmpty(suite.T(), errorResponse.Err)
}

func (suite *CustomerControllerTestSuite) TestRegisterCustomerApi_FailedUsecase() {
	dummyCustomer := dummyCustomer[0]
	suite.mockUsecase.On("RegisterCustomer", dummyCustomer).Return(errors.New("failed"))
	NewCustomerController(suite.mockRouter, suite.mockUsecase)
	r := httptest.NewRecorder()

	requestBody, _ := json.Marshal(dummyCustomer)
	request, _ := http.NewRequest(http.MethodPost, "/customer", bytes.NewBuffer(requestBody))
	suite.mockRouter.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)

	var errorResponse struct{ Err string }
	response := r.Body.String()
	json.Unmarshal([]byte(response), &errorResponse)
	assert.Equal(suite.T(), "failed", errorResponse.Err)
}

func TestCustomerUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerControllerTestSuite))
}
