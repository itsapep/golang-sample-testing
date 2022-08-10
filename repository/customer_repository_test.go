package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"

	"enigmacamp.com/golatihanlagi/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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

type CustomerRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock // github.com/DATA-DOG/go-sqlmock
}

func (suite *CustomerRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("Error when opening data stub connection", err)
	}
	suite.mockDb = mockDb
	suite.mockSql = mockSql
}

func (suite *CustomerRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func (suite *CustomerRepositoryTestSuite) TestCustomerRetrieveAll_Success() {
	// expected
	rows := sqlmock.NewRows([]string{"id", "nama", "address"})
	for _, value := range dummyCustomer {
		rows.AddRow(value.Id, value.Nama, value.Address)
	}
	suite.mockSql.ExpectQuery("select \\* from customer").WillReturnRows(rows)
	// actual
	repo := NewCustomerDbRepository(suite.mockDb)
	actual, err := repo.RetrieveAll()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, len(actual))
	assert.Equal(suite.T(), "C001", actual[0].Id)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerRetrieveAll_Failed() {
	// expected
	rows := sqlmock.NewRows([]string{"id", "nama", "address"})
	for _, value := range dummyCustomer {
		rows.AddRow(value.Id, value.Nama, value.Address)
	}
	suite.mockSql.ExpectQuery("select \\* from customer").WillReturnError(errors.New("failed"))
	// actual
	repo := NewCustomerDbRepository(suite.mockDb)
	actual, err := repo.RetrieveAll()
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), 0, len(actual))
	// assert.Equal(suite.T(), "", actual[0].Id)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerFindById_Success() {
	// expected
	dummyCustomer := dummyCustomer[0]
	rows := sqlmock.NewRows([]string{"id", "nama", "address"})
	rows.AddRow(dummyCustomer.Id, dummyCustomer.Nama, dummyCustomer.Address)
	suite.mockSql.ExpectQuery("select \\* from customer where id").WillReturnRows(rows)
	// actual
	repo := NewCustomerDbRepository(suite.mockDb)
	actual, err := repo.FindById(dummyCustomer.Id)
	fmt.Println(actual)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), actual)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerFindById_Failed() {
	// expected
	dummyCustomer := dummyCustomer[0]
	rows := sqlmock.NewRows([]string{"id", "nama", "address"})
	rows.AddRow(dummyCustomer.Id, dummyCustomer.Nama, dummyCustomer.Address)
	suite.mockSql.ExpectQuery("select \\* from customer where id").WillReturnError(errors.New("failed"))
	// actual
	repo := NewCustomerDbRepository(suite.mockDb)
	actual, err := repo.FindById(dummyCustomer.Id)
	fmt.Println(actual)
	func() {
		defer func() {
			if r := recover(); r == nil {
				assert.Error(suite.T(), err)
			}
		}()
		repo.FindById(dummyCustomer.Id)
	}()
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Customer{}, actual)
	assert.NotEqual(suite.T(), dummyCustomer, actual)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerCreate_Success() {
	// expected
	// rows := sqlmock.NewRows([]string{"id", "nama", "address"})
	// rows.AddRow("C003","Apep","Surabaya")
	// dummyCustomer := dummyCustomer[0]
	suite.mockSql.ExpectExec("insert into customer values").WithArgs(dummyCustomer[0].Id, dummyCustomer[0].Nama, dummyCustomer[0].Address).WillReturnResult(sqlmock.NewResult(1, 1))
	// actual
	repo := NewCustomerDbRepository(suite.mockDb)
	err := repo.Create(dummyCustomer[0])
	assert.Nil(suite.T(), err)
	// assert.NotNil(suite.T(), actual)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerCreate_Failed() {
	// expected
	dummyCustomer := dummyCustomer[0]
	suite.mockSql.ExpectExec("insert into customer values").WillReturnError(errors.New("failed"))
	// actual
	repo := NewCustomerDbRepository(suite.mockDb)
	err := repo.Create(dummyCustomer)
	assert.Error(suite.T(), err)
	// assert.NotNil(suite.T(), actual)
}

func TestCustomerUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()
}
