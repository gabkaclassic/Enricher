package executors

import (
	"encoding/json"
	"enricher/internal/enricher/cache"
	"enricher/internal/enricher/dto"
	"fmt"
	"log"
	"os/exec"
)

type EnricherCmdExecutorService EnricherExecutorService

func NewEnricherCmdExecutorService(cacheClient *cache.CacheClient, enrichers map[dto.EnricherArgType][]dto.Enricher) *EnricherCmdExecutorService {
	execService := &EnricherCmdExecutorService{
		processor: CmdExecute,
		enrichers: enrichers,
		cacheClient: *cacheClient,
	}

	return execService
}

func CmdExecute(enricher dto.Enricher, enricherData dto.EnricherInputData) (dto.EnricherResult, error) {
	log.Printf("Execute enricher: %s", enricher.Name)
	path := enricher.ExecutablePath
	args := make([]string, 0, len(enricherData.Data))
	for key, value := range enricherData.Data {
		args = append(args, fmt.Sprintf("%s=%v", key, value))
	}

	cmd := exec.Command(path, args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Execute cmd error: %v, %s", err, output)
		return dto.EnricherResult{}, err
	}

	var result dto.EnricherResult
	err = json.Unmarshal(output, &result)

	if err != nil {
		log.Printf("Unmarshal error: %v", err)
		return dto.EnricherResult{}, err
	}

	return result, nil
}