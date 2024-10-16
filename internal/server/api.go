package server

import (
	"context"
	"fmt"
	"github.com/LeeZXin/remote-sqlite/internal/sqlite"
	"github.com/LeeZXin/remote-sqlite/reqvo"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	secret string
)

func ListenAndServe() {
	initViper()
	initDataPath()
	secret = vp.GetString("secret")
	port := vp.GetInt("http.port")
	if port <= 0 {
		port = 15899
	}
	//gin mode
	gin.SetMode(gin.ReleaseMode)
	//create gin
	engine := gin.New()
	engine.UseH2C = true
	engine.ContextWithFallback = true
	router(engine)
	serv := &http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		ReadTimeout: 60 * time.Second,
		IdleTimeout: 60 * time.Second,
		Handler:     engine.Handler(),
		ErrorLog:    log.New(io.Discard, "", 0),
	}
	go func() {
		err := serv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("start http server failed with err: %v", err)
		}
	}()
	log.Printf("start http server port: %v", port)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	serv.Shutdown(context.Background())
}

func auth(c *gin.Context) {
	if secret != c.GetHeader("Rs-Secret") {
		c.String(http.StatusUnauthorized, "")
		c.Abort()
	} else {
		c.Next()
	}
}

func router(e *gin.Engine) {
	group := e.Group("/api/v1", auth)
	{
		group.POST("/newNamespace", newNamespace)
		group.POST("/deleteNamespace", deleteNamespace)
		group.POST("/showNamespace", showNamespace)
		group.POST("/createDB", createDB)
		group.POST("/executeCommand", executeCommand)
		group.POST("/queryCommand", queryCommand)
		group.POST("/dropDB", dropDB)
		group.POST("/getDBSize", getDBSize)
	}
}

func newNamespace(c *gin.Context) {
	var req reqvo.NewNamespaceReqVO
	if shouldBindJSON(&req, c) {
		err := sqlite.NewNamespace(dataPath, req.Namespace)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, "")
		}
	}
}

func deleteNamespace(c *gin.Context) {
	var req reqvo.DeleteNamespaceReqVO
	if shouldBindJSON(&req, c) {
		err := sqlite.DeleteNamespace(dataPath, req.Namespace)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, "")
		}
	}
}

func showNamespace(c *gin.Context) {
	var req reqvo.ShowNamespaceReqVO
	if shouldBindJSON(&req, c) {
		ret, err := sqlite.ShowNamespace(dataPath, req.Namespace)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, ret)
		}
	}
}

func createDB(c *gin.Context) {
	var req reqvo.CreateDBReqVO
	if shouldBindJSON(&req, c) {
		err := sqlite.CreateDB(dataPath, req.Namespace, req.DbName+".db")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, "")
		}
	}
}

func executeCommand(c *gin.Context) {
	var req reqvo.ExecuteCommandReqVO
	if shouldBindJSON(&req, c) {
		affectedRows, err := sqlite.ExecuteCommand(dataPath, req.Namespace, req.DbName+".db", req.Cmd)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"affectedRows": affectedRows,
			})
		}
	}
}

func queryCommand(c *gin.Context) {
	var req reqvo.QueryCommandReqVO
	if shouldBindJSON(&req, c) {
		result, err := sqlite.QueryCommand(dataPath, req.Namespace, req.DbName+".db", req.Cmd)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func dropDB(c *gin.Context) {
	var req reqvo.DropDBReqVO
	if shouldBindJSON(&req, c) {
		err := sqlite.DropDB(dataPath, req.Namespace, req.DbName+".db")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, "")
		}
	}
}

func getDBSize(c *gin.Context) {
	var req reqvo.GetDBSizeReqVO
	if shouldBindJSON(&req, c) {
		size, err := sqlite.GetDBSize(dataPath, req.Namespace, req.DbName+".db")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"size": size,
			})
		}
	}
}

func shouldBindJSON(obj reqvo.Validator, c *gin.Context) bool {
	err := c.ShouldBindJSON(obj)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return false
	}
	if !obj.IsValid() {
		c.String(http.StatusBadRequest, "")
		return false
	}
	return true
}
