package router

import (
	"context"
	"electrumx-proxy-go/common/log"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/satshub/go-bitcoind/go-electrum/electrum"
)

type addressUtxoRequest struct {
	Hex    bool `json:"hex,omitempty"`
	Amount int  `json:"amount,omitempty"`
}

type UtxoStatus struct {
	BlockHeight uint32 `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   uint64 `json:"block_time"`
	Confirmed   bool   `json:"confirmed"`
}
type addressUtxoResponse struct {
	Height uint32     `json:"height"`
	Value  uint64     `json:"value"`
	TxId   string     `json:"txid"`
	Vout   uint32     `json:"vout"`
	Status UtxoStatus `json:"status"`
}

type mempoolAddressUtxoResponse struct {
	//Height uint32     `json:"height"`
	Value  uint64     `json:"value"`
	TxId   string     `json:"txid"`
	Vout   uint32     `json:"vout"`
	Status UtxoStatus `json:"status"`
}

/*
var request addressUtxoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
*/
func GetAddressUtxoWithoutAmount(c *gin.Context) {
	address := c.Param("address")

	client, err := electrum.NewClientTCP(context.Background(), "node.sathub.io:60601")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	scriptHash, err := electrum.AddressToElectrumScriptHash(address)
	if err != nil {
		log.Error("address to script hash error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utxos, err := client.ListUnspent(context.Background(), scriptHash)
	if err != nil {
		log.Error("list unspent error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]mempoolAddressUtxoResponse, 0)

	for _, utxo := range utxos {
		log.Infof("utxo: %+v", utxo)
		tx, err := client.GetTransaction(context.Background(), utxo.Hash)
		if err != nil {
			continue
		}
		if tx.Confirmations == 0 {
			continue
		}
		response = append(response, mempoolAddressUtxoResponse{
			Value: utxo.Value,
			TxId:  utxo.Hash,
			Vout:  utxo.Position,
			Status: UtxoStatus{
				BlockHeight: utxo.Height,
				BlockHash:   tx.Blockhash,
				BlockTime:   tx.Blocktime,
				Confirmed:   true}})
	}

	c.JSON(http.StatusOK, response)
}

func GetAddressUtxoWithAmount(c *gin.Context) {
	address := c.Param("address")

	client, err := electrum.NewClientTCP(context.Background(), "node.sathub.io:60601")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	scriptHash, err := electrum.AddressToElectrumScriptHash(address)
	if err != nil {
		log.Error("address to script hash error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utxos, err := client.ListUnspent(context.Background(), scriptHash)
	if err != nil {
		log.Error("list unspent error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]addressUtxoResponse, 0)
	var spentAmount uint64

	uintValue, err := strconv.ParseUint(c.Query("amount"), 10, 0)
	if err != nil {
		log.Error("convert amount error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	for _, utxo := range utxos {
		if spentAmount >= uint64(uintValue) {
			break
		}

		log.Infof("utxo: %+v", utxo)
		tx, err := client.GetTransaction(context.Background(), utxo.Hash)
		if err != nil {
			continue
		}
		if tx.Confirmations == 0 {
			continue
		}
		response = append(response, addressUtxoResponse{
			Height: utxo.Height,
			Value:  utxo.Value,
			TxId:   utxo.Hash,
			Vout:   utxo.Position,
			Status: UtxoStatus{
				BlockHeight: utxo.Height,
				BlockHash:   tx.Blockhash,
				BlockTime:   tx.Blocktime,
				Confirmed:   true}})

		spentAmount += utxo.Value
	}

	c.JSON(http.StatusOK, response)
}

func GetAddressUtxo(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
		return
	}
	amount := c.Query("amount")
	if amount == "" {
		GetAddressUtxoWithoutAmount(c)
		return
	}
	GetAddressUtxoWithAmount(c)
}

func GetRawTransaction(c *gin.Context) {
	txid := c.Param("txid")
	if txid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "txid is required"})
		return
	}

	client, err := electrum.NewClientTCP(context.Background(), "node.sathub.io:60601")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	txHex, err := client.GetRawTransaction(context.Background(), txid)
	if err != nil {
		log.Error("get raw tx error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() + "...get raw tx error"})
		return
	}

	// 将 Hex 字符串转换为 []byte
	rawTx, err := hex.DecodeString(txHex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode transaction"})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\"raw\"")

	// 返回二进制数据
	c.Data(http.StatusOK, "application/octet-stream", rawTx)
}
