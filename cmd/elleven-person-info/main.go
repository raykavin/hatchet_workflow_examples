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
	ellevenPersonInfoWorkflow := registerEllevenPersonInfoWorkflow(client)

	// Create worker
	worker, err := client.NewWorker(
		"elleven-person-info-worker",
		hatchet.WithWorkflows(ellevenPersonInfoWorkflow),
		hatchet.WithSlots(1),
	)
	if err != nil {
		panic(err)
	}

	// Start worker
	if err := worker.StartBlocking(context.Background()); err != nil {
		panic(fmt.Errorf("error starting worker: %w", err))
	}

}

func registerEllevenPersonInfoWorkflow(client *hatchet.Client) *hatchet.Workflow {
	workflow := client.NewWorkflow(
		"elleven-person-info-workflow",
		hatchet.WithWorkflowDescription("Get elleven person info workflow"),
		hatchet.WithWorkflowVersion("1.0.0"),
	)

	workflow.NewTask("get-elleven-person-info-task", tasks.EllevenPersonInfo)

	return workflow
}
