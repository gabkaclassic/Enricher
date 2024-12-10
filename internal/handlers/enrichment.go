package handlers

import (
	"bytes"
	"encoding/json"
	"enricher/internal/common"
	"enricher/internal/enricher/dto"
	"enricher/internal/enricher/executors"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func SendEnrichmentResult(enrichmentResult dto.EnricherResult, url string) (bool, error) {

	resultJson, err := json.Marshal(enrichmentResult)

	if err != nil {
		log.Printf("Error marshalling enrichment result: %v", err)
		return false, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(resultJson))

	if err != nil {
		log.Printf("Error sending enriched result: %v", err)
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Error sending enriched result: %v", resp.Status)
		log.Printf(errorMessage)
		return false, errors.New(errorMessage)
	}

	return true, nil
}

func Enrichment(response http.ResponseWriter, request *http.Request, executor executors.EnricherExecutor) {

	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)

	if err != nil {
		http.Error(response, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error reading request body: %v", err)
		return
	}
	defer request.Body.Close()

	var inputEnricher dto.EnricherInputData

	if err := json.Unmarshal(body, &inputEnricher); err != nil {
		http.Error(response, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error unmarshalling request body: %v", err)
		return
	}

	go func() {
		results, err := executor.ExecuteEnrichers(inputEnricher)
		if err != nil {
			log.Printf("Error executing enricher: %v", err)
			return
		}

		_, err = common.ParallelExecute(inputEnricher.WebhookUri, results, SendEnrichmentResult)
		if err != nil {
			log.Printf("Error sending enriched results: %v", err)
			return
		}
	}()

	response.WriteHeader(http.StatusAccepted)
	response.Write([]byte("Enrichment process started"))
}
