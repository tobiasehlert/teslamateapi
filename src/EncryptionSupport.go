package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

// decryptAccessToken funct to decrypt tokens from database
func decryptAccessToken(data string, encryption_key string) string {

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
	h.Write([]byte(encryption_key))
	fmt.Printf("key: %x \n", h.Sum(nil))

	key := h.Sum(nil)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// first byte
	key_type, _ := strconv.ParseInt(data[0:2], 16, 64)
	// second byte
	key_len, _ := strconv.ParseInt(data[2:4], 16, 64)
	key_tag, _ := hex.DecodeString(data[4 : 4+(key_len*2)])
	fmt.Printf("Type: %d \n", key_type)
	fmt.Printf("Length: %d \n", key_len)
	fmt.Printf("Key Tag: %s \n", key_tag)

	/*
	   With AES.GCM, 12-byte IV length is necessary for interoperability reasons.
	   See https://github.com/danielberkompas/cloak/issues/93
	   IV and nonce are often used interchangeably. Essentially though, an IV is a nonce with an additional requirement: it must be selected in a non-predictable way
	   https://medium.com/@fridakahsas/salt-nonces-and-ivs-whats-the-difference-d7a44724a447#:~:text=IV%20and%20nonce%20are%20often,an%20IV%20must%20be%20random.
	*/

	nonce, _ := hex.DecodeString(data[4+(key_len*2) : 4+(key_len*2)+24])
	fmt.Printf("IV: %x \n", nonce)

	ciphertag, _ := hex.DecodeString(data[4+(key_len*2)+24 : 4+(key_len*2)+24+32])
	fmt.Printf("Ciphertag: %x \n", ciphertag)

	aesgcm, err := cipher.NewGCMWithTagSize(block, 16)
	if err != nil {
		panic(err.Error())
	}

	// https://stackoverflow.com/a/68353192
	// golang aes expects cipertag to append ciphertext....
	ciphertext_tag, _ := hex.DecodeString(data[4+(key_len*2)+24+32:] + data[4+(key_len*2)+24:4+(key_len*2)+24+32])

	// AES256GCM -- Additional Authenticated Data (AAD)
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext_tag, []byte("AES256GCM"))
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Decrypted: %s\n", plaintext)

	return string(plaintext)
}
