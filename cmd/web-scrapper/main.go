package main

import (
	"context"
	"flag"
	"fmt"
	"workflow-test/pkg/tasks"

	hatchet "github.com/hatchet-dev/hatchet/sdks/go"
	"github.com/joho/godotenv"
)

func main() {
	flag.Parse()

	// Load environment variables
	godotenv.Load()

	// Initialize Hatchet client
	client, err := hatchet.NewClient()
	if err != nil {
		panic(err)
	}

	// Register workflow
	siteTextScrapperWorkflow := registerSiteScrapperWorkflows(client)

	// Create worker
	worker, err := client.NewWorker(
		"scrapper-worker",
		hatchet.WithWorkflows(siteTextScrapperWorkflow),
		hatchet.WithSlots(1),
	)
	if err != nil {
		panic(fmt.Errorf("failed creating worker: %w", err))
	}

	// Start worker
	if err := worker.StartBlocking(context.Background()); err != nil {
		panic(fmt.Errorf("error starting worker: %w", err))
	}
}

func registerSiteScrapperWorkflows(client *hatchet.Client) *hatchet.Workflow {
	workflow := client.NewWorkflow(
		"website-scrapper",
		hatchet.WithWorkflowDescription("Examples of website text scrapper using hatchet workflow"),
		hatchet.WithWorkflowVersion("1.0.0"),
	)

	workflow.NewTask("text-scrapper-task", tasks.TextScrapperTask)

	return workflow
}
