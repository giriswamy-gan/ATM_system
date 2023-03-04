package routes

import (
	"atm-system/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.Engine) {
	r.POST("/create", controllers.CreateAccount)
	r.POST("/deposit", controllers.Deposit)
	r.POST("/withdraw", controllers.Withdraw)
	r.POST("/transfer", controllers.Transfer)
	r.POST("/setpin", controllers.SetPin)
	r.POST("/bankstatement", controllers.BankStatement)
}
