package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type CarRegionAPI string

const (
	ChinaAPI  CarRegionAPI = "China"
	GlobalAPI              = "Global"
)

// decryptAccessToken funct to decrypt tokens from database
func decryptAccessToken(data string, encryptionKey string) string {

	/*
	   From Adrian....
	   I had a look at how to decode the binary input without additional libraries. Below is sample code for Elixir. An important detail is that  "Additional Authenticated Data (AAD) " is required to decrypt the tokens. The AAD is a fixed string, in this case "AES256GCM”.
	   << _type::bytes-1, length::integer, _tag::bytes-size(length), iv::bytes-12, ciphertag::bytes-16, ciphertext::bytes >> = input
	   key = :crypto.hash(:sha256, key)
	   aad = "AES256GCM"
	   plaintext = :crypto.crypto_one_time_aead(:aes_256_gcm, key, iv, ciphertext, aad, ciphertag, false)

	   How the encrypted content looks like....
	   +----------------------------------------------------------+----------------------+
	   |                          HEADER                          |         BODY         |
	   +-------------------+---------------+----------------------+----------------------+
	   | Key Tag (n bytes) | IV (n bytes)  | Ciphertag (16 bytes) | Ciphertext (n bytes) |
	   +-------------------+---------------+----------------------+----------------------+
	   |                   |_________________________________
	   |                                                     |
	   +---------------+-----------------+-------------------+
	   | Type (1 byte) | Length (1 byte) | Key Tag (n bytes) |
	   +---------------+-----------------+-------------------+
	*/

	h := sha256.New()
	h.Write([]byte(encryptionKey))
	if gin.IsDebugging() {
		log.Printf("[debug] decryptAccessToken - Key: %x \n", h.Sum(nil))
	}

	key := h.Sum(nil)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// first byte
	keyType := int([]rune(data)[0])
	// second byte
	keyLen := int([]rune(data)[1])
	keyTag := data[2 : 2+keyLen]
	if gin.IsDebugging() {
		log.Printf("[debug] decryptAccessToken - Type: %d \n", keyType)
		log.Printf("[debug] decryptAccessToken - Length: %d \n", keyLen)
		log.Printf("[debug] decryptAccessToken - Key Tag: %s \n", keyTag)
	}

	/*
	   With AES.GCM, 12-byte IV length is necessary for interoperability reasons.
	   See https://github.com/danielberkompas/cloak/issues/93
	   IV and nonce are often used interchangeably. Essentially though, an IV is a nonce with an additional requirement: it must be selected in a non-predictable way
	   https://medium.com/@fridakahsas/salt-nonces-and-ivs-whats-the-difference-d7a44724a447#:~:text=IV%20and%20nonce%20are%20often,an%20IV%20must%20be%20random.
	*/

	nonce := data[2+keyLen : 2+keyLen+12]
	if gin.IsDebugging() {
		log.Printf("[debug] decryptAccessToken - IV (hex): %x \n", nonce)

		ciphertag := data[2+keyLen+12 : 2+keyLen+12+16]
		log.Printf("[debug] decryptAccessToken - Ciphertag (hex): %x \n", ciphertag)
	}

	aesgcm, err := cipher.NewGCMWithTagSize(block, 16)
	if err != nil {
		panic(err.Error())
	}

	// https://stackoverflow.com/a/68353192
	// golang aes expects cipertag to append ciphertext....
	ciphertextTag := data[2+keyLen+12+16:] + data[2+keyLen+12:2+keyLen+12+16]

	// AES256GCM -- Additional Authenticated Data (AAD)
	plaintext, err := aesgcm.Open(nil, []byte(nonce), []byte(ciphertextTag), []byte("AES256GCM"))
	if err != nil {
		panic(err.Error())
	}

	if gin.IsDebugging() {
		// fmt.Printf("[debug] decryptAccessToken - Decrypted: %s\n", plaintext)
	}

	return string(plaintext)
}

// getCarRegionAPI function to get URL from iis in accessToken
func getCarRegionAPI(accessToken string) CarRegionAPI {
	payload := strings.Split(accessToken, ".")
	if len(payload) != 3 {
		return GlobalAPI
	}
	decodedStr, err := base64.RawStdEncoding.DecodeString(payload[1])
	if err != nil {
		return GlobalAPI
	}
	var result map[string]interface{}
	if err = json.Unmarshal(decodedStr, &result); err != nil {
		return GlobalAPI
	}
	issUrl, err := url.Parse(result["iss"].(string))
	if err != nil {
		return GlobalAPI
	}
	if strings.HasSuffix(issUrl.Host, ".cn") {
		return ChinaAPI
	}
	return GlobalAPI
}
