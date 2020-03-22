package helpers

import (
	htpasswd "github.com/foomo/htpasswd"
	"golang.org/x/crypto/bcrypt"
)

//HashPassword Creates password hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPasswordHash checks password hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetMD5Hash(user_name, password string) (map[string]string, error) {
	passwords := htpasswd.HashedPasswords(map[string]string{})
	err := passwords.SetPassword(user_name, password, htpasswd.HashAPR1)
	return passwords, err
}
