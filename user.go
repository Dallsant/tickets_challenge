package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string   `json:"username" gorm:"type:varchar(300);"`
	Email    string   `json:"email" gorm:"type:varchar(700);"`
	Password string   `json:"password" gorm:"type:varchar(900);"`
	Tickets  []Ticket `gorm:"foreignKey:UserId;references:id"`
}

type BlacklistedToken struct {
	gorm.Model
	Token string `json:"token"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenerateJWT(user string) (string, error) {
	// Token expires in one year
	expirationTime := time.Now().Add(525600 * time.Minute)
	claims := &Claims{
		Username: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SigningKey)
	if err != nil {
		logger(err.Error(), "Token Generation")
		return "", err
	}
	return tokenString, nil
}

func extractClaims(tokenString string) (jwt.MapClaims, bool) {

	hmacSecretString := SigningKey
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {

		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		logger(err.Error(), "Hashing Password")
	}
	return string(hash)
}

// ENDPOINTS

func login(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var credentials Credentials
	json.Unmarshal(reqBody, &credentials)
	bytesPwd := []byte(credentials.Password)
	var user User
	db.Where("username = ?", credentials.Username).Find(&user)
	if err != nil {
		logger(err.Error(), "Login")
		return
	}
	jwt, err := GenerateJWT(credentials.Username)

	if err != nil {
		logger(err.Error(), "Token generation")
	}
	if ComparePasswords(user.Password, bytesPwd) {
		t := make(map[string]string)
		t["token"] = jwt

		json.NewEncoder(w).Encode(t)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	var alreadyExists []User
	json.Unmarshal(reqBody, &user)
	db.Where("Username = ?", user.Username).Find(&alreadyExists)

	if len(alreadyExists) > 0 {
		http.Error(w, "User already exists", http.StatusUnprocessableEntity)

	}
	bytesPwd := []byte(user.Password)
	hash := HashAndSalt(bytesPwd)
	user.Password = hash
	db.Create(&user)
	if err != nil {
		logger(err.Error(), "Register")
	}
	json.NewEncoder(w).Encode(user)
}
