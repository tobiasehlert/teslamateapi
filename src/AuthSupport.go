package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// initAuthToken func
func initAuthToken() {
	// get token from environment variable API_TOKEN
	envToken = getEnv("API_TOKEN", "")
	if envToken == "" {
		log.Println("[warning] initAuthToken - environment variable API_TOKEN not set or is empty.")
	} else if len(envToken) < 32 {
		log.Println("[warning] initAuthToken - environment variable API_TOKEN too short.. should be 32 or longer.")
	} else {
		log.Println("[info] initAuthToken - environment variable API_TOKEN is set and good.")
	}
}

// validateAuthToken func
func validateAuthToken(c *gin.Context) (bool, string) {

	// if API_TOKEN_DISABLE is true, skip token validation.
	if getEnvAsBool("API_TOKEN_DISABLE", false) == true {
		log.Println("[debug] validateAuthToken - header authorization bearer token disabled.")
		return true, ""
	}

	// trying with http header - Authorization: Bearer <token>
	reqHeaderToken := c.Request.Header.Get("Authorization")

	// if length of reqHeaderToken is more than zero
	if len(reqHeaderToken) > 0 {

		// removing Bearer part from header to get token out of header
		splitToken := strings.Split(reqHeaderToken, "Bearer")

		if len(splitToken) != 2 {
			// bearer token is not proper formatted.. returning bad request
			log.Println("[info] validateAuthToken - header authorization bearer token is not proper formatted.. returning 401")
			return false, "header authorization bearer token is not proper formatted"

		} else if strings.TrimSpace(splitToken[1]) == "" {
			// bearer token is empty string.. we'll return unauthorized
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
