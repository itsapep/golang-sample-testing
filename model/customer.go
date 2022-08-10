package model

import "fmt"

type Customer struct {
	Id      string
	Nama    string
	Address string
}

func (c Customer) String() {
	fmt.Println(c.Id, c.Nama, c.Address)
}

func NewCustomer(id, name, address string) Customer {
	return Customer{
		Id:      id,
		Nama:    name,
		Address: address,
	}

}
