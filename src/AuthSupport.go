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

// validateAuthToken func
func validateAuthToken(c *gin.Context) (bool, string) {

	// if API_TOKEN_DISABLE is true, skip token validation.
	if getEnvAsBool("API_TOKEN_DISABLE", false) {
		log.Println("[debug] validateAuthToken - header authorization bearer token disabled.")
		return true, ""
	}

	// trying with configured http header
	reqHeaderToken := c.Request.Header.Get(apiTokenHeaderName)

	// if header provided
	if len(reqHeaderToken) > 0 {

		if strings.EqualFold(apiTokenHeaderName, "Authorization") {
			// expecting format: Authorization: Bearer <token>
			splitToken := strings.Split(reqHeaderToken, "Bearer")

			// bearer token is not proper formatted.. returning bad request
			if len(splitToken) != 2 {
				log.Println("[info] validateAuthToken - header authorization bearer token is not proper formatted.. returning 401")
				return false, "header authorization bearer token is not proper formatted"

			} else if strings.TrimSpace(splitToken[1]) == "" {
				log.Println("[info] validateAuthToken - header authorization bearer token is empty.. returning 401")
				return false, "header authorization bearer token is empty"

			} else if checkAuthToken(strings.TrimSpace(splitToken[1])) {
				// the bearer token is valid!
				log.Println("[debug] validateAuthToken - header authorization bearer token valid.")
				return true, ""

			}
			// the check did fail.. bearer token is invalid
			log.Println("[info] validateAuthToken - header authorization bearer token invalid.. returning 401")
			return false, "header authorization bearer token invalid"

		} else {
			// custom header - treat value as raw token (no Bearer prefix)
			tokenValue := strings.TrimSpace(reqHeaderToken)
			if tokenValue == "" {
				log.Println("[info] validateAuthToken - configured header value is empty.. returning 401")
				return false, "header token empty"
			}
			if checkAuthToken(tokenValue) {
				log.Println("[debug] validateAuthToken - token from custom header valid.")
				return true, ""
			}
			log.Println("[info] validateAuthToken - token from custom header invalid.. returning 401")
			return false, "header token invalid"
		}
	}

	// trying with http parameter - ?token=<token>
	tokenParamsValue := c.DefaultQuery("token", "")

	// if validTokenParams is longer than zero
	if len(tokenParamsValue) > 0 {

		// checking if token is valid (since it's over zero length)
		if checkAuthToken(tokenParamsValue) {
			// the token is valid!
			log.Println("[debug] validateAuthToken - param token valid.")
			return true, ""

		}
		// the token is invalid.
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
