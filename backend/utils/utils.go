package utils

import (
	"errors"
	"fmt"
	"math"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateNewToken(issuer int, secret []byte) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": issuer,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(8 * time.Hour).Unix(),
	}).SignedString(secret)
}

func GetClaimsFromToken(header string, secret []byte) (map[string]interface{}, error) {
	s := strings.Split(header, " ")
	if len(s) != 2 {
		return nil, errors.New("malformed Authorization header.")
	}

	token, err := jwt.Parse(s[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method.")
		}

		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("malformed token claims.")
	}
}

func VerifyCredentials(hash []byte, password string) bool {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return false
	}

	return true
}

func SaltPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func ValidateUsername(name string) bool {
	return len(name) < 18 && regexp.MustCompile(`^.*[a-zA-Z0-9].*$`).MatchString(name)
}

func ValidatePassword(password string) bool {
	var regex = []*regexp.Regexp{
		regexp.MustCompile(`^.*[a-z].*[a-z].*$`),
		regexp.MustCompile(`^.*[A-Z].*[A-Z].*$`),
		regexp.MustCompile(`^.*[0-9].*[0-9].*$`),
		regexp.MustCompile(`^.*[!@#$%&].*$`),
		regexp.MustCompile(`^.{8,}$`),
	}

	for _, v := range regex {
		if !v.MatchString(password) {
			return false
		}
	}

	return true
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func ValidateName(name string) bool {
	if regexp.MustCompile(`^.*[ ].*$`).MatchString(name) {
		return false
	}

	return true
}

func FormatQuery(query string) string {
	return strings.ToLower(fmt.Sprintf("%%%s%%", strings.ReplaceAll(query, "%20", "%")))
}

const alpha = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const MaxFileSize = 32 << 20 // ~30MB

func IntToBase62(n int) string {
	result := string(alpha[n%62])
	q := int(math.Floor(float64(n / 62)))

	for q > 0 {
		result = string(alpha[q%62]) + result
		q = int(math.Floor(float64(q / 62)))
	}

	return result
}

func GenerateUploadName(filename string, number int) string {
	split := strings.SplitN(filename, ".", 2)

	if len(split) > 1 {
		return split[0] + "-" + IntToBase62(number+1) + "." + split[1]
	}

	return split[0] + "-" + IntToBase62(number+1)
}
