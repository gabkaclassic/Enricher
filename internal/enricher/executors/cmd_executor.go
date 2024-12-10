package executors

import (
	"encoding/json"
	"enricher/internal/enricher/cache"
	"enricher/internal/enricher/dto"
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
	
	cmd := exec.Command(path, enricherData.Data)

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