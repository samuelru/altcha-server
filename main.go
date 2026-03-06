package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/altcha-org/altcha-lib-go"
)

var secret = ""
var challengeTTL = 1 * time.Hour
var allowedOrigins []string
var isDev = false

func enableCORS(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin != "" {
		allowed := false
		if isDev {
			allowed = true
		} else {
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func challengeHandler(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	options := altcha.ChallengeOptions{
		HMACKey: secret,
		Expires: new(time.Now().Add(challengeTTL)),
	}

	challenge, err := altcha.CreateChallenge(options)
	if err != nil {
		http.Error(w, "Failed to create challenge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenge)

	// Log requested challenges
	f, err := os.OpenFile("challenges.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		logEntry := fmt.Sprintf("Time: %d, Challenge: %+v\n", challenge.Salt, challenge)
		f.WriteString(logEntry)
	}
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Payload string `json:"payload"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check for reuse
	if exists, _ := isPayloadReused(payload.Payload); exists {
		logVerification("REUSED", payload.Payload)
		http.Error(w, "Payload already used", http.StatusForbidden)
		return
	}

	valid, err := altcha.VerifySolution(payload.Payload, secret, true)
	if err != nil {
		logVerification("ERROR", payload.Payload)
		http.Error(w, "Verification failed", http.StatusInternalServerError)
		return
	}

	if !valid {
		logVerification("INVALID", payload.Payload)
		http.Error(w, "Invalid Altcha payload", http.StatusForbidden)
		return
	}

	// Record the solution to prevent reuse
	saveVerifiedSolution(payload.Payload)
	logVerification("SUCCESS", payload.Payload)

	fmt.Fprintln(w, "Verification successful")
}

func isPayloadReused(payload string) (bool, error) {
	f, err := os.Open("solutions.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer f.Close()

	// Using bufio for memory efficiency and correctness (line by line matching)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == payload {
			return true, nil
		}
	}

	return false, scanner.Err()
}

func saveVerifiedSolution(payload string) {
	f, err := os.OpenFile("solutions.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(payload + "\n")
}

func logVerification(status, payload string) {
	f, err := os.OpenFile("verifications.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("Status: %s, Payload: %s\n", status, payload))
}

func main() {
	secret = os.Getenv("ALTCHA_SECRET")
	if secret == "" {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			log.Fatalf("Failed to generate secret key: %v", err)
		}
		secret = hex.EncodeToString(b)
		log.Printf("ALTCHA_SECRET not set, use a random secret key")
	}

	if ttlStr := os.Getenv("ALTCHA_TTL"); ttlStr != "" {
		if d, err := time.ParseDuration(ttlStr); err == nil {
			challengeTTL = d
			log.Printf("ALTCHA_TTL set to %s", challengeTTL)
		} else {
			log.Printf("Invalid ALTCHA_TTL value %q, using default 1h: %v", ttlStr, err)
		}
	}

	if originStr := os.Getenv("ALTCHA_CORS_ORIGIN"); originStr != "" {
		allowedOrigins = strings.Split(originStr, ",")
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
		log.Printf("ALTCHA_CORS_ORIGIN set to %v", allowedOrigins)
	} else {
		allowedOrigins = []string{"*"}
	}

	if os.Getenv("IS_DEV") == "true" {
		isDev = true
		log.Printf("IS_DEV is enabled, allowing all CORS origins")
	}

	http.HandleFunc("/challenge", challengeHandler)
	http.HandleFunc("/verify", verifyHandler)

	fmt.Println("Server starting on :3947...")
	log.Fatal(http.ListenAndServe(":3947", nil))
}
