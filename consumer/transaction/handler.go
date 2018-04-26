package transaction

import (
	"encoding/json"
	"fmt"
	"log"

	nsq "github.com/bitly/go-nsq"
	"github.com/o3labs/neo-utils/neoutils/neorpc"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
	"github.com/o3labs/transaction-log/models"
	"github.com/o3labs/transaction-log/services"
)

type MessageHandler struct {
	Service      *services.TransactionLogService
	NEORPCClient *neorpc.NEORPCClient
}

func (m *MessageHandler) OnFinish(message *nsq.Message) {

}

func (h *MessageHandler) HandleMessage(message *nsq.Message) error {
	tx := models.TransactionLog{}
	json.Unmarshal(message.Body, &tx)

	//get transaction detail or block detail
	if tx.TXType == "block" {
		//fetch block
		result := h.NEORPCClient.GetBlock(tx.TXID)
		tx.Data, _ = models.JSONBFromObject(result.Result)
	} else if tx.TXType == "tx" {
		result := h.NEORPCClient.GetRawTransaction(tx.TXID)
		tx.Data, _ = models.JSONBFromObject(result.Result)
		if result.Result.Type == "InvocationTransaction" {
			log.Printf("txid = %v", tx.TXID)
			parser := smartcontract.NewParserWithScript(result.Result.Script)
			scripts, err := parser.GetListOfScriptHashes()
			if err == nil {
				log.Printf("scripthash = %+v", scripts)
			}
		}
	}
	_, err := h.Service.SaveTransaction(tx)

	if err != nil {
		if message.Attempts > 10 {
			message.Finish()
			return fmt.Errorf("too many attempt")
		}
		log.Printf("error %+v", err)
		message.Requeue(1)
		return err
	}
	message.Finish()
	return nil
}
