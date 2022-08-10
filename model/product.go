package model

import "fmt"

type Product struct {
	Nama string
}
func (c *Product) GetCustomerInfo() string {
	return fmt.Sprintf("Nama : %s, Alaamat: %s", c.Nama)
}
