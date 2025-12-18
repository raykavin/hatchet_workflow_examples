package main

import (
	"context"
	"fmt"
	"workflow-test/pkg/tasks"

	hatchet "github.com/hatchet-dev/hatchet/sdks/go"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// Initialize Hatchet client
	client, err := hatchet.NewClient()
	if err != nil {
		panic(fmt.Errorf("failed to create hatchet client: %w", err))
	}

	// Register workflow
	jsonDescriptionsWorkflow := registerJSONDescriptionsWorkflow(client)

	// Create worker
	worker, err := client.NewWorker(
		"json-descriptions-worker",
		hatchet.WithWorkflows(jsonDescriptionsWorkflow),
		hatchet.WithSlots(1),
	)
	if err != nil {
		panic(fmt.Errorf("failed creating worker: %w", err))
	}

	fmt.Printf("Input file:  %s\n", tasks.InputFileName)
	fmt.Printf("Output file: %s\n", tasks.OutputFileName)
	fmt.Printf("Chunk size: %d routes\n", tasks.ChunkSize)
	fmt.Printf("Max concurrency: %d goroutines\n", tasks.MaxConcurrency)

	// Start worker
	if err := worker.StartBlocking(context.Background()); err != nil {
		panic(fmt.Errorf("error starting worker: %w", err))
	}
}

// registerJSONDescriptionsWorkflow registers the JSON descriptions extraction workflow
func registerJSONDescriptionsWorkflow(client *hatchet.Client) *hatchet.Workflow {
	workflow := client.NewWorkflow(
		"extract-json-descriptions",
		hatchet.WithWorkflowDescription("Extracts all description fields from routes-large.json to descriptions.txt"),
		hatchet.WithWorkflowVersion("1.0.0"),
	)

	workflow.NewTask("find-json-descriptions", tasks.ProcessJSONFileTask)

	return workflow
}
