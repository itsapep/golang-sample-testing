package usecase

import (
	"errors"
	"testing"

	"enigmacamp.com/golatihanlagi/model"
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

type repoMock struct {
	mock.Mock
}

// Create implements repository.CustomerRepository
func (r *repoMock) Create(newCustomer model.Customer) error {
	args := r.Called(newCustomer)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return nil
}

// FindById implements repository.CustomerRepository
func (r *repoMock) FindById(id string) (model.Customer, error) {
	args := r.Called(id)
	if args.Get(1) != nil {
		return model.Customer{}, args.Error(1)
	}
	return args.Get(0).(model.Customer), nil
}

// RetrieveAll implements repository.CustomerRepository
func (r *repoMock) RetrieveAll() ([]model.Customer, error) {
	args := r.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Customer), nil
}

type CustomerUsecaseTestSuite struct {
	suite.Suite
	repoMock *repoMock
}

func (suite *CustomerUsecaseTestSuite) SetupTest() {
	suite.repoMock = new(repoMock)
}

func (suite *CustomerUsecaseTestSuite) TestCustomerFindById_Success() {
	// expected
	dummyCustomer := dummyCustomer[0]
	suite.repoMock.On("FindById", dummyCustomer.Id).Return(dummyCustomer, nil)
	// actual
	customerUseCaseTest := NewCustomerUseCase(suite.repoMock)
	customer, err := customerUseCaseTest.FindCustomerById(dummyCustomer.Id)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyCustomer.Id, customer.Id)

}
func (suite *CustomerUsecaseTestSuite) TestCustomerFindById_Failed() {
	dummyCustomer := dummyCustomer[0]
	suite.repoMock.On("FindById", dummyCustomer.Id).Return(model.Customer{}, errors.New("failed"))
	customerUseCaseTest := NewCustomerUseCase(suite.repoMock)
	customer, err := customerUseCaseTest.FindCustomerById(dummyCustomer.Id)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "failed", err.Error())
	assert.Equal(suite.T(), "", customer.Id)
}

func (suite *CustomerUsecaseTestSuite) TestCustomerRetrieveAll_Success() {
	suite.repoMock.On("RetrieveAll").Return(dummyCustomer, nil)
	customerUseCaseTest := NewCustomerUseCase(suite.repoMock)
	customer, err := customerUseCaseTest.GetAllCustomer()
	assert.NotEmpty(suite.T(), customer) // bisa pilih salah satu aja
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyCustomer, customer)

}
func (suite *CustomerUsecaseTestSuite) TestCustomerRetrieveAll_Failed() {
	suite.repoMock.On("RetrieveAll").Return(nil, errors.New("failed"))
	customerUseCaseTest := NewCustomerUseCase(suite.repoMock)
	customer, err := customerUseCaseTest.GetAllCustomer()
	assert.Empty(suite.T(), customer)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), []model.Customer(nil), customer)
}

func (suite *CustomerUsecaseTestSuite) TestCustomerCreate_Success() {
	dummyCustomer := dummyCustomer[0]
	suite.repoMock.On("Create", dummyCustomer.Id).Return(dummyCustomer, nil)
	customerUseCaseTest := NewCustomerUseCase(suite.repoMock)
	customer, err := customerUseCaseTest.FindCustomerById(dummyCustomer.Id)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyCustomer.Id, customer.Id)

}
func (suite *CustomerUsecaseTestSuite) TestCustomerCreate_Failed() {
	dummyCustomer := dummyCustomer[0]
	suite.repoMock.On("Create", dummyCustomer.Id).Return(model.Customer{}, errors.New("failed"))
	customerUseCaseTest := NewCustomerUseCase(suite.repoMock)
	customer, err := customerUseCaseTest.FindCustomerById(dummyCustomer.Id)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "failed", err.Error())
	assert.Equal(suite.T(), "", customer.Id)
}
func TestCustomerUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerUsecaseTestSuite))
}
