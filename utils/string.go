package utils

import (
	"math/big"
	"strings"
	"time"

	random "crypto/rand"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

func GenerateRandomString(length int) string {
	// Set karakter yang bisa digunakan
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	rand.Seed(uint64(time.Now().UnixNano())) // Seed untuk random generator berdasarkan waktu

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset)) // Pilih index acak dari charset
		result.WriteByte(charset[randomIndex]) // Tambahkan karakter ke hasil
	}

	return result.String() // Kembalikan string acak
}

func GenerateRandomNumeric(length int) (string, error) {
	const digits = "0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		// ambil angka acak dari 0 sampai 9
		n, err := random.Int(random.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		result[i] = digits[n.Int64()]
	}

	return string(result), nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
