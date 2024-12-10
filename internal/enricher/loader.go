package enricher

import (
	"encoding/json"
	"enricher/internal/enricher/dto"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func loadEnrichers(enrichersPath string) (map[dto.EnricherArgType][]dto.Enricher, error) {
	var enrichers = make(map[dto.EnricherArgType][]dto.Enricher)

	entries, err := os.ReadDir(enrichersPath)
	if err != nil {
		log.Printf("Error reading Enrichers from path %s: %v", enrichersPath, err)
		return nil, err
	}

	var wg sync.WaitGroup
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		wg.Add(1)
		go func(entry os.DirEntry) {
			defer wg.Done()
			subDirPath := filepath.Join(enrichersPath, entry.Name())
			settingsPath := filepath.Join(subDirPath, "settings.json")

			info, err := os.Stat(settingsPath)
			if err != nil || info.IsDir() {
				return
			}
			file, err := os.ReadFile(settingsPath)
			if err != nil {
				log.Printf("Read enricher config file error: %v", err)
				return
			}

			var enricher dto.Enricher
			err = json.Unmarshal(file, &enricher)

			if err != nil {
				log.Printf("Unmarshal enricher error: %v", err)
				return
			}

			executablePath := filepath.Join(subDirPath, enricher.ExecutablePath)

			absoluteExecutablePath, err := filepath.Abs(executablePath)
			if err != nil {
				log.Printf("Error converting to absolute path: %v", err)
				return
			}
			enricher.ExecutablePath = absoluteExecutablePath

			err = validateEnricher(enricher)

			if err != nil {
				log.Printf("Validation enricher %s failed: %v", enricher.Name, err)
				return
			}

			for _, allowedType := range enricher.AllowedTypes {
				enrichers[allowedType] = append(enrichers[allowedType], enricher)
			}
		}(entry)
	}
	wg.Wait()

	if len(enrichers) == 0 {
		log.Println("Enrichers not found")
	}

	return enrichers, err
}

type enricherManager struct {
	enrichers map[dto.EnricherArgType][]dto.Enricher
	mu        sync.RWMutex
}

func NewEnricherManager() *enricherManager {
	return &enricherManager{
		enrichers: nil,
	}
}

func (manager *enricherManager) LoadEnrichers(enrichersPath string) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	loadedEnrichers, err := loadEnrichers(enrichersPath)
	if err != nil {
		return err
	}

	manager.enrichers = loadedEnrichers

	return nil
}

func (manager *enricherManager) GetEnrichers(enrichersPath string) (map[dto.EnricherArgType][]dto.Enricher, error) {

	if manager.enrichers == nil || len(manager.enrichers) == 0 {
		err := manager.LoadEnrichers(enrichersPath)
		if err != nil {
			return nil, err
		}
	}

	return manager.enrichers, nil
}
