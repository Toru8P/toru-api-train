package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/lora-trainer/train", TrainLora)
	r.GET("/lora-trainer/logs/:job_id", GetLoraLogs)
	r.GET("/lora-trainer/status/:job_id", CheckLoraStatus)
	r.POST("/lora-trainer/cancel/:job_id", CancelLoraJob)

	r.Run(":8000") // expose port
}
