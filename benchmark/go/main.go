package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "regexp"
    "strings"
    "time"

    "github.com/google/uuid"
    _ "github.com/lib/pq"
)

type User struct {
    Name      string   `json:"name"`
    Email     string   `json:"email"`
    Age       int      `json:"age"`
    Bio       string   `json:"bio"`
    Interests []string `json:"interests"`
}

var db *sql.DB

func validateUser(user *User) error {
    if len(user.Name) < 2 || len(user.Name) > 50 {
        return httpError("Invalid name length", 400)
    }
    if !strings.HasSuffix(user.Email, "@example.com") {
        return httpError("Email must end with @example.com", 400)
    }
    if user.Age < 0 || user.Age > 120 {
        return httpError("Invalid age", 400)
    }
    if len(user.Bio) < 10 || !regexp.MustCompile(`(?i)(engineer|developer|programmer)`).MatchString(user.Bio) {
        return httpError("Bio must be longer and contain job keyword", 400)
    }
    if len(user.Interests) == 0 {
        return httpError("At least one interest required", 400)
    }
    return nil
}

func httpError(msg string, code int) error {
    return &httpErrorType{msg, code}
}

type httpErrorType struct {
    msg  string
    code int
}

func (e *httpErrorType) Error() string { return e.msg }

func handler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if err := validateUser(&user); err != nil {
        if herr, ok := err.(*httpErrorType); ok {
            http.Error(w, herr.msg, herr.code)
        } else {
            http.Error(w, "Validation error", 400)
        }
        return
    }

    id := uuid.New().String()
    _, err := db.Exec("INSERT INTO users (id, name, email) VALUES ($1, $2, $3)", id, user.Name, user.Email)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func main() {
    var err error
    db, err = sql.Open("postgres", "postgres://benchmark:benchmark@postgres:5432/benchmark?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }

    // 커넥션 풀 설정
    db.SetMaxOpenConns(50)               // 최대 동시 연결 수
    db.SetMaxIdleConns(25)               // 유휴 커넥션 수
    db.SetConnMaxLifetime(5 * time.Minute) // 연결 생명 주기

    // 연결 확인
    if err := db.Ping(); err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }

    http.HandleFunc("/data", handler)
    log.Println("Go server running on port 8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}