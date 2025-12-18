package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	hatchet "github.com/hatchet-dev/hatchet/sdks/go"
)

const (
	InputFileName  = "routes-large.json"
	OutputFileName = "descriptions.txt"
	ChunkSize      = 50
	MaxConcurrency = 10
)

// ProcessJSONFileOutput represents the final output of the workflow
type ProcessJSONFileOutput struct {
	InputFile          string  `json:"input_file"`
	OutputFile         string  `json:"output_file"`
	TotalRoutes        int     `json:"total_routes"`
	TotalDescriptions  int     `json:"total_descriptions"`
	TotalChunks        int     `json:"total_chunks"`
	ProcessingTime     float64 `json:"processing_time_seconds"`
	DescriptionsPerSec float64 `json:"descriptions_per_second"`
	Success            bool    `json:"success"`
}

// ProcessJSONFileTask reads routes-large.json, extracts all descriptions, and writes to descriptions.txt
// This task does everything in one go using concurrent processing
func ProcessJSONFileTask(
	ctx hatchet.Context,
	input struct{},
) (ProcessJSONFileOutput, error) {
	startTime := time.Now()

	log.Printf("Reading file: %s", InputFileName)

	// Read JSON file
	data, err := os.ReadFile(InputFileName)
	if err != nil {
		return ProcessJSONFileOutput{}, fmt.Errorf("failed to read file %s: %w", InputFileName, err)
	}

	// Parse JSON as array
	var routes []any
	if err := json.Unmarshal(data, &routes); err != nil {
		return ProcessJSONFileOutput{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	log.Printf("Parsed %d routes from JSON", len(routes))

	// Divide into chunks
	chunks := divideIntoChunks(routes, ChunkSize)
	log.Printf("Divided into %d chunks of ~%d routes each", len(chunks), ChunkSize)

	//  Process chunks concurrently
	log.Printf("Processing chunks with %d concurrent goroutines...", MaxConcurrency)
	allDescriptions := processChunksConcurrently(chunks)

	log.Printf("Extracted %d total descriptions", len(allDescriptions))

	// Write to output file
	if len(allDescriptions) == 0 {
		return ProcessJSONFileOutput{}, fmt.Errorf("no descriptions found in JSON file")
	}

	content := strings.Join(allDescriptions, "\n")
	if err := os.WriteFile(OutputFileName, []byte(content), 0644); err != nil {
		return ProcessJSONFileOutput{}, fmt.Errorf("failed to write output file: %w", err)
	}

	processingTime := time.Since(startTime).Seconds()
	descriptionsPerSec := float64(len(allDescriptions)) / processingTime

	output := ProcessJSONFileOutput{
		InputFile:          InputFileName,
		OutputFile:         OutputFileName,
		TotalRoutes:        len(routes),
		TotalDescriptions:  len(allDescriptions),
		TotalChunks:        len(chunks),
		ProcessingTime:     processingTime,
		DescriptionsPerSec: descriptionsPerSec,
		Success:            true,
	}

	log.Printf("Output file: %s", OutputFileName)
	log.Printf("Total descriptions: %d", len(allDescriptions))
	log.Printf("Processing time: %.4f seconds", processingTime)
	log.Printf("Throughput: %.2f descriptions/sec", descriptionsPerSec)

	return output, nil
}

// divideIntoChunks divides routes into smaller chunks for parallel processing
func divideIntoChunks(routes []any, chunkSize int) [][]any {
	chunks := make([][]any, 0)

	for i := 0; i < len(routes); i += chunkSize {
		end := i + chunkSize
		if end > len(routes) {
			end = len(routes)
		}
		chunks = append(chunks, routes[i:end])
	}

	return chunks
}

// processChunksConcurrently processes all chunks using goroutines with concurrency limit
func processChunksConcurrently(chunks [][]any) []string {
	results := make([][]string, len(chunks))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, MaxConcurrency)

	// Process each chunk concurrently
	for i, chunk := range chunks {
		wg.Add(1)
		go func(chunkID int, routes []any) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Extract descriptions from this chunk
			descriptions := []string{}
			for _, route := range routes {
				extractDescriptions(route, &descriptions)
			}

			results[chunkID] = descriptions
			log.Printf("  Chunk #%d: found %d descriptions", chunkID, len(descriptions))
		}(i, chunk)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Flatten results
	allDescriptions := []string{}
	for _, chunkDescriptions := range results {
		allDescriptions = append(allDescriptions, chunkDescriptions...)
	}

	return allDescriptions
}

// extractDescriptions recursively traverses JSON structure and extracts all "description" fields
func extractDescriptions(data any, descriptions *[]string) {
	switch v := data.(type) {
	case map[string]any:
		// Check for "description" key in map
		for key, value := range v {
			if key == "description" {
				if strValue, ok := value.(string); ok && strValue != "" {
					*descriptions = append(*descriptions, strValue)
				}
			} else {
				// Recursively process nested structures
				extractDescriptions(value, descriptions)
			}
		}
	case []any:
		// Process each element in array
		for _, item := range v {
			extractDescriptions(item, descriptions)
		}
	}
}
