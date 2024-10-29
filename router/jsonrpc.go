package router

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/satshub/go-bitcoind/go-electrum/electrum"
)

type addressUtxoRequest struct {
	Hex    bool `json:"hex,omitempty"`
	Amount int  `json:"amount,omitempty"`
}

type addressUtxoResponse struct {
	Height uint32 `json:"height"`
	Value  uint64 `json:"value"`
	TxId   string `json:"txid"`
	Vout   uint32 `json:"vout"`
}

func GetAddressUtxo(c *gin.Context) {
	address := c.Param("address")

	var request addressUtxoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client, err := electrum.NewClientTCP(context.Background(), "node.sathub.io:60601")

	if err != nil {
		log.Fatal(err)
	}

	scriptHash, err := electrum.AddressToElectrumScriptHash(address)
	if err != nil {
		log.Fatal("address to script hash error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utxos, err := client.ListUnspent(context.Background(), scriptHash)
	if err != nil {
		log.Fatal("list unspent error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"list unspent error": err.Error()})
		return
	}

	response := make([]addressUtxoResponse, 0)
	var spentAmount uint64

	uintValue, err := strconv.ParseUint(c.Query("amount"), 10, 0)
	if err != nil {
		log.Fatal("convert amount error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"convert amount error": err.Error()})
		return

	}

	for _, utxo := range utxos {
		if spentAmount >= uint64(uintValue) {
			break
		}

		log.Printf("utxo: %+v", utxo)
		response = append(response, addressUtxoResponse{
			Height: utxo.Height,
			Value:  utxo.Value,
			TxId:   utxo.Hash,
			Vout:   utxo.Position,
		})
		spentAmount += utxo.Value
	}

	c.JSON(http.StatusOK, response)
}
