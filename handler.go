package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TrainLora(c *gin.Context) {
	jobID, err := StartTrainingContainer([]string{"", ""})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"job_id": jobID})
}

func GetLoraLogs(c *gin.Context) {
	jobID := c.Param("job_id")
	logs, err := FetchContainerLogs(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, logs)
}

func CheckLoraStatus(c *gin.Context) {
	jobID := c.Param("job_id")
	status, err := GetContainerStatus(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func CancelLoraJob(c *gin.Context) {
	jobID := c.Param("job_id")
	err := StopContainer(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}
