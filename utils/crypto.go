package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func CreateHmac(data, key string) string {
	// Buat HMAC menggunakan SHA-256
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data)) // Update dengan data

	// Konversi hasilnya ke string heksadesimal
	return hex.EncodeToString(h.Sum(nil))
}

func CreateHmac512(data, key string) string {
	h := hmac.New(sha512.New, []byte(key))
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}

func GenerateMD5Hash(merchantSecretKey, signature, timestamp string) string {
	// Gabungkan string sesuai urutan
	data := merchantSecretKey + signature + timestamp

	// Buat hash MD5
	hash := md5.Sum([]byte(data))

	// Konversi hasil hash ke string heksadesimal
	return hex.EncodeToString(hash[:])
}

func HashPassword(p string) string {
	password := []byte(p)
	hash, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return string(hash)
}
func ComparePass(h, p []byte) bool {
	hash, pass := []byte(h), []byte(p)
	err := bcrypt.CompareHashAndPassword(hash, pass)
	return err == nil
}
