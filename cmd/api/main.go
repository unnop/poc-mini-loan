package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"

	"poc-mini-loan/internal/config"
	"poc-mini-loan/internal/loan"
	"poc-mini-loan/internal/pkg"
	"poc-mini-loan/models"
)

func main() {

	cfg := config.LoadConfig()

	gin.SetMode(cfg.GinMode)

	c := pkg.InitTemporalClient(cfg)
	defer c.Close()

	r := gin.Default()

	r.POST("/api/v1/loans/apply", func(ctx *gin.Context) {
		var req models.LoanRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		workflowID := "merchant-loan-" + req.MerchantID

		workflowOptions := client.StartWorkflowOptions{
			ID:        workflowID,
			TaskQueue: cfg.TaskQueue,
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, loan.MerchantLoanWorkflow, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": " Workflow start failed : " + err.Error()})
			return
		}

		resp := models.LoanResponse{
			Message:    "Your loan submitted",
			WorkflowID: we.GetID(),
			RunID:      we.GetRunID(),
			Status:     "PENDING_PRE_SCREENING",
			CreatedAt:  time.Now().Format(time.RFC3339),
		}
		ctx.JSON(http.StatusAccepted, resp)
	})

	r.POST("/api/v1/loans/approve", func(ctx *gin.Context) {
		var req models.ApprovalRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		workflowID := "merchant-loan-" + req.MerchantID

		log.Printf("📡 [API] Sending signal [%s] to Workflow ID: %s", req.Action, workflowID)

		err := c.SignalWorkflow(context.Background(), workflowID, "", "loan-approval-signal", req.Action)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send signal to workflow: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":     "Signal sent successfully",
			"merchant_id": req.MerchantID,
			"action":      req.Action,
		})
	})

	log.Printf("🚀 API ready port :%s", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
