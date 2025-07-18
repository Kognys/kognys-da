package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MOSSV2/dimo-sdk-go/lib/log"
	"github.com/MOSSV2/dimo-sdk-go/lib/piece"
	"github.com/MOSSV2/dimo-sdk-go/lib/repo"
	"github.com/MOSSV2/dimo-sdk-go/lib/types"
	"github.com/MOSSV2/dimo-sdk-go/lib/utils"
	"github.com/MOSSV2/dimo-sdk-go/sdk"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

var logger = log.Logger("storage")

type StorageServer struct {
	Router *gin.Engine

	typ string

	rp repo.Repo

	ps types.IPieceStore

	sync.RWMutex
	storedData map[string][]byte // Simple in-memory storage for demo

	local common.Address

	auth types.Auth
}

func NewStorageServer(rp repo.Repo) (*http.Server, error) {
	log.SetLogLevel("DEBUG")

	gin.SetMode(gin.ReleaseMode)

	localAddr := rp.Key().Address()

	logger.Infof("storage node %s starting...", localAddr)

	router := gin.Default()

	auth, err := rp.Key().BuildAuth([]byte("storage"))
	if err != nil {
		return nil, err
	}

	s := &StorageServer{
		Router: router,

		typ:        types.StoreType,
		local:      localAddr,
		rp:         rp,
		ps:         piece.New(rp.MetaStore(), rp.DataStore()),
		auth:       auth,
		storedData: make(map[string][]byte),
	}

	err = s.register()
	if err != nil {
		logger.Warn("Failed to register with remote server:", err)
		// Continue anyway for local testing
	}

	s.registRoute()

	srv := &http.Server{
		Addr:    rp.Config().API.Endpoint,
		Handler: s.Router,
	}

	return srv, nil
}

func (s *StorageServer) registRoute() {
	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	
	s.Router.Use(cors.New(config))

	s.Router.Use(ginzap.Ginzap(log.Logger("gin").Desugar(), time.RFC3339, true))

	// Health check endpoint
	s.Router.GET("/health", s.healthHandler)

	// API routes
	r := s.Router.Group("/api")

	s.addInfo(r)
	s.addStore(r)
	s.addRetrieve(r)
	s.addUpload(r)
	s.addDownload(r)
}

func (s *StorageServer) register() error {
	auth, err := s.rp.Key().BuildAuth([]byte("register"))
	if err != nil {
		return err
	}

	mm := types.EdgeMeta{
		Type:      s.typ,
		Name:      auth.Addr,
		PublicKey: s.rp.Key().Public(),
		ExposeURL: s.rp.Config().API.Expose,
		Hardware:  utils.GetHardwareInfo(),
		ChainType: s.rp.Config().Chain.Type,
	}

	err = sdk.RegisterEdge(s.rp.Config().Remote.URL, auth, mm)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageServer) healthHandler(c *gin.Context) {
	response := map[string]interface{}{
		"status":    "healthy",
		"type":      "storage-node",
		"chain_type": s.rp.Config().Chain.Type,
		"address":   s.local.Hex(),
		"timestamp": time.Now(),
	}
	c.JSON(http.StatusOK, response)
}

func (s *StorageServer) addInfo(g *gin.RouterGroup) {
	g.GET("/info", func(c *gin.Context) {
		res := types.EdgeReceipt{
			EdgeMeta: types.EdgeMeta{
				Type: s.typ,
				Name: s.local,
				ExposeURL: s.rp.Config().API.Expose,
				ChainType: s.rp.Config().Chain.Type,
			},
		}

		c.JSON(http.StatusOK, res)
	})
}

func (s *StorageServer) addStore(g *gin.RouterGroup) {
	g.POST("/store", func(c *gin.Context) {
		var req struct {
			Key  string `json:"key" binding:"required"`
			Data string `json:"data" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s.Lock()
		s.storedData[req.Key] = []byte(req.Data)
		s.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"key":     req.Key,
			"size":    len(req.Data),
		})
	})
}

func (s *StorageServer) addRetrieve(g *gin.RouterGroup) {
	g.GET("/retrieve/:key", func(c *gin.Context) {
		key := c.Param("key")

		s.RLock()
		data, exists := s.storedData[key]
		s.RUnlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"key":  key,
			"data": string(data),
			"size": len(data),
		})
	})
}

// Compatible with the mock server endpoints
func (s *StorageServer) addUpload(g *gin.RouterGroup) {
	g.POST("/upload", func(c *gin.Context) {
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Store the entire JSON data
		id, ok := data["id"].(string)
		if !ok {
			id = "upload_" + time.Now().Format("20060102150405")
		}

		jsonData, _ := json.Marshal(data)
		
		s.Lock()
		s.storedData[id] = jsonData
		s.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Upload successful",
			"id":      id,
		})
	})
}

func (s *StorageServer) addDownload(g *gin.RouterGroup) {
	g.GET("/download", func(c *gin.Context) {
		name := c.Query("name")
		owner := c.Query("owner")

		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name parameter required"})
			return
		}

		s.RLock()
		data, exists := s.storedData[name]
		s.RUnlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		// Try to parse as JSON
		var jsonData map[string]interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			// If not JSON, return as string
			c.JSON(http.StatusOK, gin.H{
				"message": "Download successful",
				"name":    name,
				"owner":   owner,
				"data":    string(data),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Download successful",
			"name":    name,
			"owner":   owner,
			"data":    jsonData,
		})
	})
}