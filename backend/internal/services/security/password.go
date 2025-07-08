package security

import "crypto/rand"

func GeneratePassword(length int) string {
	const characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%^&*"

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
