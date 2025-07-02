package security

import "crypto/rand"

func Generate2FA(length int) string {
	var characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	code := make([]byte, length)
	_, err := rand.Read(code)
	if err != nil {
		return ""
	}

	result := ""
	for _, char := range code {
		result += string(characters[char%byte(len(characters))])
	}

	return result
}
