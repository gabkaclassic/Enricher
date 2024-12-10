package main

import (
	"enricher/configs"
	"enricher/internal/enricher"
	"enricher/internal/enricher/cache"
	"enricher/internal/enricher/dto"
	"enricher/internal/enricher/executors"
	"enricher/internal/server/handlers"
	"enricher/internal/server/middlewares"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Enricher service starting...")

	configPath := flag.String("config", "configs/application.yml", "Path to config file")
	flag.Parse()

	err := run(*configPath)

	if err != nil {
		log.Println("Enricher service emergency stopped")
		panic(err)
	}

	log.Println("Enricher service stopped successfully")
}

func run(configPath string) error {
	appConfig, err := getConfigs(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configs: %w", err)
	}

	enrichers, err := getEnrichers(*appConfig.Enrichers)
	if err != nil {
		return fmt.Errorf("failed to load enrichers: %w", err)
	}

	cacheClient := getCacheClient(*appConfig.Cache)

	enricherExecutorService := getExecutorService(cacheClient, enrichers)

	if err := startServer(*appConfig.Server, *appConfig.API, enricherExecutorService); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func getConfigs(configPath string) (configs.Config, error) {

	log.Println("Loading config...")
	configManager := configs.ConfigManager{}
	appConfig, err := configManager.GetConfig(configPath)

	if err != nil {
		log.Printf("Error while config load: %v", err)
		return configs.Config{}, err
	}
	log.Println("Configs loaded successfully")

	return *appConfig, nil
}

func getEnrichers(config configs.EnrichersConfig) (map[dto.EnricherArgType][]dto.Enricher, error) {
	log.Println("Loading enrichers...")
	enricherManager := enricher.NewEnricherManager()

	enrichers, err := enricherManager.GetEnrichers(config.Path)

	if err != nil {
		log.Printf("Error while enrichers load: %v", err)
		return nil, err
	}
	log.Println("Configs loaded successfully")

	return enrichers, nil
}

func startServer(config configs.ServerConfig, apiConfig configs.APIConfig, executorService executors.EnricherExecutorService) error {
	log.Println("Configuring server...")

	enrichmentHandler := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		handlers.Enrichment(response, request, executorService)
	})

	authEnrichmentHandler := middlewares.AuthMiddleware(enrichmentHandler, apiConfig)

	http.Handle("/enrichment", authEnrichmentHandler)

	serverHost := fmt.Sprintf("%s:%d", config.Host, config.Port)
	fmt.Printf("Starting server %s...\n", serverHost)
	err := http.ListenAndServe(serverHost, nil)

	return err
}

func getCacheClient(config configs.CacheConfig) cache.CacheClient {
	log.Println("Creating cache client...")

	if config.Address == "" {
		return cache.NewInMemoryCacheClient()
	}
	return cache.NewRedisCacheClient(
		config.Address,
		config.Password,
		config.Db,
	)
}

func getExecutorService(cacheClient cache.CacheClient, enrichers map[dto.EnricherArgType][]dto.Enricher) executors.EnricherExecutorService {
	log.Println("Creating enrichers executor service...")

	return executors.EnricherExecutorService(*executors.NewEnricherCmdExecutorService(
		&cacheClient,
		enrichers,
	))
}
