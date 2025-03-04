package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func generateOTP() string {
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}

func generateJWT(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)
	otp := generateOTP()
	fmt.Println("OTP Generated:", otp)

	_, err = db.Exec("INSERT INTO users (email, password, role, verified, otp) VALUES ($1, $2, $3, $4, $5)", creds.Email, creds.Password, "user", false, otp)
	if err != nil {
		fmt.Println("Database Error:", err)
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User registered successfully. Please verify your email using OTP: %s", otp)
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	json.NewDecoder(r.Body).Decode(&creds)
	result, err := db.Exec("UPDATE users SET verified = true WHERE email = $1 AND otp = $2", creds.Email, creds.OTP)
	if err != nil {
		fmt.Println("Database Error:", err)
		http.Error(w, "Failed to verify email", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Invalid OTP or Email", http.StatusBadRequest)
		return
	}
	token, err := generateJWT(creds.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Email verified successfully! Your new Token: %s", token)
}

func forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&creds)
	otp := generateOTP()
	fmt.Println("OTP Generated for Reset Password:", otp)
	result, err := db.Exec("UPDATE users SET otp = $1 WHERE email = $2", otp, creds.Email)
	if err != nil {
		fmt.Println("Database Error:", err)
		http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Email not found", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "OTP for reset password has been sent to your email: %s", creds.Email)
}

func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)
	result, err := db.Exec("UPDATE users SET password = $1 WHERE email = $2 AND otp = $3", creds.NewPassword, creds.Email, creds.OTP)
	if err != nil {
		fmt.Println("Database Error:", err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Invalid OTP or Email", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Password successfully reset for %s", creds.Email)
}

func main() {
	db, err = sql.Open("postgres", "user=postgres password=postgres dbname=testdb sslmode=disable")
	if err != nil {
		fmt.Println("Database Connection Error:", err)
		return
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/verify-email", verifyEmailHandler).Methods("POST")
	router.HandleFunc("/forgot-password", forgotPasswordHandler).Methods("POST")
	router.HandleFunc("/reset-password", resetPasswordHandler).Methods("POST")

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", router)
}
