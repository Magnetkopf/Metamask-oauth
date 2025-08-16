package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"Metamask-oauth/common"
	"Metamask-oauth/server"
	"github.com/gin-gonic/gin"
)

func initLogger() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	common.Logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func cleanupExpiredRequests() {
	for range time.Tick(1 * time.Minute) {
		common.Mu.Lock()
		now := time.Now().Unix()
		//1min
		expirationTime := int64(60)
		for addr, req := range common.LoginRequests {
			if now-req.Timestamp > expirationTime {
				delete(common.LoginRequests, addr)
				common.Logger.Printf("[cleanup]Address: %s 超时\n", addr)
			}
		}
		common.Mu.Unlock()
	}
}

func main() {
	initLogger()

	go cleanupExpiredRequests()

	r := gin.Default()
	r.LoadHTMLGlob("view/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	r.POST("/request-message", func(c *gin.Context) { server.HandleRequestMessage(c) })
	r.POST("/login", func(c *gin.Context) { server.HandleLogin(c) })

	fmt.Println("Started at :9090")
	common.Logger.Println("Started")
	err := r.Run(`:9090`)
	if err != nil {
		return
	}
}
