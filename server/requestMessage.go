package server

import (
	"Metamask-oauth/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func HandleRequestMessage(c *gin.Context) {
	var json struct {
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No empty"})
		return
	}

	address := strings.ToLower(json.Address)
	nonce, err := common.GenerateNonce(16) //random string with 32 char
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成随机内容"})
		return
	}

	messageToSign := fmt.Sprintf("MetaMask oauth\n签名此消息以证明您是地址的所有者\nSign thie message to show your ownership of this address\n\n随机码: %s", nonce)

	common.Mu.Lock()
	common.LoginRequests[address] = common.LoginRequest{
		Address:   address,
		Nonce:     messageToSign,
		Timestamp: time.Now().Unix(),
	}
	common.Mu.Unlock()

	common.Logger.Printf("Address %s->需要对 %s 签名\n", address, nonce)

	c.JSON(http.StatusOK, gin.H{"messageToSign": messageToSign})
}
