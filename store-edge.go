package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "sync"
    "time"
)

type StorageNode struct {
    sync.RWMutex
    data map[string][]byte
    chainType string
    port string
}

func (s *StorageNode) healthHandler(w http.ResponseWriter, r *http.Request) {
    response := map[string]interface{}{
        "status": "healthy",
        "type": "storage-node",
        "chain_type": s.chainType,
        "port": s.port,
        "timestamp": time.Now(),
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (s *StorageNode) infoHandler(w http.ResponseWriter, r *http.Request) {
    response := map[string]interface{}{
        "type": "store",
        "name": "unibase-storage-node",
        "chainType": s.chainType,
        "exposeURL": os.Getenv("EXPOSE_URL"),
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (s *StorageNode) uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var data map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    id := fmt.Sprintf("%v", data["id"])
    if id == "<nil>" {
        id = fmt.Sprintf("file_%d", time.Now().Unix())
    }
    
    jsonData, _ := json.Marshal(data)
    s.Lock()
    s.data[id] = jsonData
    s.Unlock()
    
    response := map[string]interface{}{
        "success": true,
        "message": "Upload successful",
        "id": id,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (s *StorageNode) downloadHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        http.Error(w, "name parameter required", http.StatusBadRequest)
        return
    }
    
    s.RLock()
    data, exists := s.data[name]
    s.RUnlock()
    
    if !exists {
        http.Error(w, "file not found", http.StatusNotFound)
        return
    }
    
    var jsonData map[string]interface{}
    json.Unmarshal(data, &jsonData)
    
    response := map[string]interface{}{
        "message": "Download successful",
        "name": name,
        "owner": r.URL.Query().Get("owner"),
        "data": jsonData,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}

func main() {
    port := os.Getenv("EXPOSE_PORT")
    if port == "" {
        port = "8082"
    }
    
    chainType := os.Getenv("CHAIN_TYPE")
    if chainType == "" {
        chainType = "bnb-testnet"
    }
    
    node := &StorageNode{
        data: make(map[string][]byte),
        chainType: chainType,
        port: port,
    }
    
    http.HandleFunc("/health", enableCORS(node.healthHandler))
    http.HandleFunc("/api/info", enableCORS(node.infoHandler))
    http.HandleFunc("/api/upload", enableCORS(node.uploadHandler))
    http.HandleFunc("/api/download", enableCORS(node.downloadHandler))
    
    log.Printf("Unibase Storage Node starting on port %s", port)
    log.Printf("Chain Type: %s", chainType)
    log.Printf("External URL: %s", os.Getenv("EXPOSE_URL"))
    
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}