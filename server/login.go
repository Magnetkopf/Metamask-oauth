package server

import (
	"Metamask-oauth/common"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func HandleLogin(c *gin.Context) {
	var json struct {
		Address   string `json:"address" binding:"required"`
		Signature string `json:"signature" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No empty"})
		return
	}

	address := strings.ToLower(json.Address)

	common.Mu.Lock()
	req, ok := common.LoginRequests[address]
	if !ok {
		common.Mu.Unlock()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid/Expired login request"})
		return
	}
	//delete after login
	delete(common.LoginRequests, address)
	common.Mu.Unlock()

	fullMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(req.Nonce), req.Nonce)
	hash := crypto.Keccak256Hash([]byte(fullMessage))
	sig, err := hexutil.Decode(json.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的签名格式"})
		return
	}

	//标准化
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}

	//get public key
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法恢复公钥"})
		return
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey).Hex()

	if strings.ToLower(recoveredAddr) == address {
		fmt.Printf("[success]Address: %s\n", address)
		common.Logger.Printf("[success]Address: %s\n", address)
		c.JSON(http.StatusOK, gin.H{"status": "登录成功", "address": address})
	} else {
		common.Logger.Printf("[error]尝试地址: %s, 恢复地址: %s\n", address, recoveredAddr)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "签名验证失败"})
	}
}
