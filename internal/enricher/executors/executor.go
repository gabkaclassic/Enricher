package executors

import (
	"encoding/json"
	"enricher/internal/common"
	"enricher/internal/enricher/cache"
	"enricher/internal/enricher/dto"
	"errors"
	"fmt"
	"log"
)

type EnricherExecutor interface {
	ExecuteEnrichers(enricherData dto.EnricherInputData) ([]dto.EnricherResult, error)
}

type EnricherExecutorService struct {
	enrichers   map[dto.EnricherArgType][]dto.Enricher
	cacheClient cache.CacheClient
	processor   common.Processor[dto.Enricher, dto.EnricherInputData, dto.EnricherResult]
}

func (executor EnricherExecutorService) getEnrichmentResultCacheKey(enricher dto.Enricher, data dto.EnricherInputData) string {
	return fmt.Sprintf("%s-%v-%s", enricher.Name, data.Data, data.DataType)
}

func (executor EnricherExecutorService) encodeCacheValue(data interface{}) ([]byte, error) {

	value, err := json.Marshal(data)

	return value, err
}

func (executor EnricherExecutorService) decodeEnrichmentResult(data []byte) (dto.EnricherResult, error) {

	var result dto.EnricherResult
	err := json.Unmarshal(data, &result)

	return result, err
}

func (executor EnricherExecutorService) putEnrichmentResultToCache(enricher dto.Enricher, data dto.EnricherInputData, result dto.EnricherResult) error {

	cacheKey := executor.getEnrichmentResultCacheKey(enricher, data)
	cacheValue, err := executor.encodeCacheValue(result.Report)

	if err != nil {
		log.Printf("Encode enrichment result error: %v", err)
		return err
	}

	err = executor.cacheClient.SetWithTTL(
		cacheKey,
		cacheValue,
		int(enricher.Timeout),
	)

	return err
}

func (executor EnricherExecutorService) getEnrichmentResultFromCache(enricher dto.Enricher, data dto.EnricherInputData) (dto.EnricherResult, error) {

	cacheKey := executor.getEnrichmentResultCacheKey(enricher, data)

	value, err := executor.cacheClient.Get(cacheKey)

	if err != nil {
		log.Printf("Get enrichemnt result from cache error: %v", err)
		return dto.EnricherResult{}, err
	}

	result, err := executor.decodeEnrichmentResult(value)

	if err != nil {
		log.Printf("Decode enrichment result error: %v", err)
		return dto.EnricherResult{}, err
	}

	return result, err
}

func (executor EnricherExecutorService) ExecuteEnrichers(enricherData dto.EnricherInputData) ([]dto.EnricherResult, error) {
	allowedEnrichersList, exists := executor.enrichers[enricherData.DataType]
	if !exists {
		errorMessage := fmt.Sprintf("Enricher for data type %s not found", string(enricherData.DataType))
		return []dto.EnricherResult{}, errors.New(errorMessage)
	}

	enabledEnrichers := getEnabledEnrichers(allowedEnrichersList)

	results := make([]dto.EnricherResult, 0, len(enabledEnrichers))
	var errors []error
	var enrichersToExecute []dto.Enricher

	for _, enricher := range enabledEnrichers {
		cachedResult, err := executor.getEnrichmentResultFromCache(enricher, enricherData)
		if err == nil {
			results = append(results, cachedResult)
		} else {
			enrichersToExecute = append(enrichersToExecute, enricher)
		}
	}

	if len(enrichersToExecute) > 0 {
		parallelResults, err := common.ParallelExecute(enricherData, enrichersToExecute, executor.processor)
		if err != nil {
			errors = append(errors, err)
		} else {
			results = append(results, parallelResults...)
			for i, enricher := range enrichersToExecute {
				executor.putEnrichmentResultToCache(enricher, enricherData, parallelResults[i])
			}
		}
	}

	return results, common.MergeErrors(errors)
}

func getEnabledEnrichers(enricher []dto.Enricher) []dto.Enricher {
	var enabledEnrichers []dto.Enricher
	for _, enricher := range enricher {
		if enricher.Enabled {
			enabledEnrichers = append(enabledEnrichers, enricher)
		}
	}
	return enabledEnrichers
}
