package middlewares

import (
	"enricher/configs"
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler, apiConfig configs.APIConfig) http.Handler {

	keys := make(map[string]string)

	for _, key := range apiConfig.Keys {
		keys[key.Key] = key.Name
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		key := request.Header.Get("Authorization")
		keyName, exists := keys[key]

		if !exists {
			log.Printf("Invalid API key: %s", key)
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("API key '%s' have access to %s", keyName, request.URL.Path)

		next.ServeHTTP(writer, request)
	})
}
