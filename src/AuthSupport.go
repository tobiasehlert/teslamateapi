package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// defining envToken that contains API_TOKEN value
	envToken string
	// header name to read API token from (default "Authorization")
	apiTokenHeaderName string
)

// initAuthToken func
func initAuthToken() {
	// get token from environment variable API_TOKEN
	envToken = getEnv("API_TOKEN", "")
	// get header name from environment variable API_TOKEN_HEADER (default: Authorization)
	apiTokenHeaderName = getEnv("API_TOKEN_HEADER", "Authorization")
	if apiTokenHeaderName == "" {
		apiTokenHeaderName = "Authorization"
	}

	if envToken == "" {
		log.Println("[warning] initAuthToken - environment variable API_TOKEN not set or is empty.")
	} else if len(envToken) < 32 {
		log.Println("[warning] initAuthToken - environment variable API_TOKEN too short.. should be 32 or longer.")
	} else {
		log.Println("[info] initAuthToken - environment variable API_TOKEN is set and good.")
	}

	log.Printf("[info] initAuthToken - using header for API token: %s", apiTokenHeaderName)
}

// helper to extract token from configured header
func extractTokenFromHeader(headerName, headerValue string) (string, string) {
	if strings.EqualFold(headerName, "Authorization") {
		// expecting format: Authorization: Bearer <token>
		splitToken := strings.Split(headerValue, "Bearer")
		// bearer token is not proper formatted.. returning bad request
		if len(splitToken) != 2 {
			return "", "header authorization bearer token is not proper formatted"
		}
		trimmed := strings.TrimSpace(splitToken[1])
		if trimmed == "" {
			return "", "header authorization bearer token is empty"
		}
		return trimmed, ""
	}

	// custom header - treat value as raw token (no Bearer prefix)
	trimmed := strings.TrimSpace(headerValue)
	if trimmed == "" {
		return "", "header token empty"
	}
	return trimmed, ""
}

// validateAuthToken func
func validateAuthToken(c *gin.Context) (bool, string) {
	// if API_TOKEN_DISABLE is true, skip token validation.
	if getEnvAsBool("API_TOKEN_DISABLE", false) {
		log.Println("[debug] validateAuthToken - header authorization bearer token disabled.")
		return true, ""
	}

	// try configured header first
	headerVal := c.Request.Header.Get(apiTokenHeaderName)
	if headerVal != "" {
		token, errMsg := extractTokenFromHeader(apiTokenHeaderName, headerVal)
		if errMsg != "" {
			log.Println("[info] validateAuthToken -", errMsg+".. returning 401")
			return false, errMsg
		}
		if checkAuthToken(token) {
			log.Println("[debug] validateAuthToken - token from header valid.")
			return true, ""
		}
		// invalid token from header
		if strings.EqualFold(apiTokenHeaderName, "Authorization") {
			log.Println("[info] validateAuthToken - header authorization bearer token invalid.. returning 401")
			return false, "header authorization bearer token invalid"
		}
		log.Println("[info] validateAuthToken - header token invalid.. returning 401")
		return false, "header token invalid"
	}

	// fallback to query parameter ?token=<token>
	tokenParamsValue := c.DefaultQuery("token", "")
	if tokenParamsValue != "" {
		if checkAuthToken(tokenParamsValue) {
			log.Println("[debug] validateAuthToken - param token valid.")
			return true, ""
		}
		log.Println("[info] validateAuthToken - param token invalid.. returning 401")
		return false, "param token invalid"
	}

	// unauthozie all calls!
	return false, "failed validation"
}

// checkAuthToken func
func checkAuthToken(token string) bool {
	// checking if it's valid or not
	if token == envToken {
		// check that envToken was longer than zero
		if len(envToken) == 0 {
			log.Println("[warning] checkAuthToken - returning false (API_TOKEN is not set or empty)")
			return false
		}
		log.Println("[info] checkAuthToken - returning true")
		return true
	}

	// failing check what so ever..
	log.Println("[info] checkAuthToken - returning false (other reason)")
	return false
}
