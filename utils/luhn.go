package utils

import (
	"math/rand"
	"strconv"
	"strings"
)

func GenerateCardNumber() string {
	prefix := "4" // Visa
	length := 16

	number := prefix
	for i := 0; i < length-2; i++ {
		number += strconv.Itoa(rand.Intn(10))
	}

	checkDigit := calculateLuhnChecksum(number)
	number += strconv.Itoa(checkDigit)
	return number
}

func ValidateCardNumber(number string) bool {
	if len(number) < 13 || len(number) > 19 {
		return false
	}

	sum := 0
	isSecond := false

	for i := len(number) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}

		if isSecond {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isSecond = !isSecond
	}

	return sum%10 == 0
}

func calculateLuhnChecksum(number string) int {
	sum := 0
	isSecond := true

	for i := len(number) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(number[i]))

		if isSecond {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isSecond = !isSecond
	}

	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit
}

func MaskCardNumber(number string) string {
	if len(number) < 8 {
		return strings.Repeat("*", len(number))
	}
	return number[:6] + strings.Repeat("*", len(number)-10) + number[len(number)-4:]
}
