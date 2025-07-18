package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

type HealthResponse struct {
    Status    string    `json:"status"`
    ChainType string    `json:"chain_type"`
    Port      string    `json:"port"`
    Timestamp time.Time `json:"timestamp"`
}

type UploadResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    ID      string `json:"id"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    response := HealthResponse{
        Status:    "healthy",
        ChainType: os.Getenv("CHAIN_TYPE"),
        Port:      os.Getenv("EXPOSE_PORT"),
        Timestamp: time.Now(),
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var data map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    response := UploadResponse{
        Success: true,
        Message: "Mock upload successful",
        ID:      fmt.Sprintf("%v", data["id"]),
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    owner := r.URL.Query().Get("owner")
    
    response := map[string]string{
        "message": "Mock download",
        "name":    name,
        "owner":   owner,
        "data":    "This is mock data from Unibase DA storage node",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    port := os.Getenv("EXPOSE_PORT")
    if port == "" {
        port = "8082"
    }
    
    http.HandleFunc("/health", healthHandler)
    http.HandleFunc("/api/upload", uploadHandler)
    http.HandleFunc("/api/download", downloadHandler)
    
    log.Printf("Mock Unibase DA Storage Node starting on port %s", port)
    log.Printf("Chain Type: %s", os.Getenv("CHAIN_TYPE"))
    log.Printf("This is a MOCK server for testing deployment")
    
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}