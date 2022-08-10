package controller

import (
	"net/http"

	"enigmacamp.com/golatihanlagi/model"
	"enigmacamp.com/golatihanlagi/usecase"
	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	router  *gin.Engine
	usecase usecase.CustomerUsecase
}

func (cc *CustomerController) getAllCustomer(ctx *gin.Context) {
	customers, err := cc.usecase.GetAllCustomer()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, customers)
}
func (cc *CustomerController) getCustomerById(ctx *gin.Context) {
	id := ctx.Param("id")
	customers, err := cc.usecase.FindCustomerById(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, customers)
}

func (cc *CustomerController) registerCustomer(ctx *gin.Context) {
	var customer model.Customer
	err := ctx.ShouldBindJSON(&customer)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	err = cc.usecase.RegisterCustomer(customer)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func NewCustomerController(r *gin.Engine, usecase usecase.CustomerUsecase) *CustomerController {
	controller := CustomerController{
		router:  r,
		usecase: usecase,
	}
	r.GET("/customer", controller.getAllCustomer)
	r.GET("/customer/:id", controller.getCustomerById)
	r.POST("/customer", controller.registerCustomer)
	return &controller
}
