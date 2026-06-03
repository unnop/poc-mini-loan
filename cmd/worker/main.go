package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"go.temporal.io/sdk/worker"

	"poc-mini-loan/internal/config"
	"poc-mini-loan/internal/loan"
	"poc-mini-loan/internal/pkg"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DBConn)
	if err != nil {
		log.Fatalf("❌ Worker can not connect Postgres : %v", err)
	}
	defer db.Close()

	c := pkg.InitTemporalClient(cfg)
	defer c.Close()

	w := worker.New(c, cfg.TaskQueue, worker.Options{})

	w.RegisterWorkflow(loan.MerchantLoanWorkflow)

	loanActs := &loan.LoanActivities{DB: db}
	w.RegisterActivity(loanActs)

	log.Printf("Mini-Loan Worker starting process...")

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("❌ Worker stop!!: %v", err)
	}
}
