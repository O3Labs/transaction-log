// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/o3labs/neo-transaction-watcher/neotx"
	"github.com/o3labs/neo-transaction-watcher/neotx/network"
	"github.com/o3labs/transaction-log/config"
	"github.com/o3labs/transaction-log/models"
	"github.com/o3labs/transaction-log/q"
)

var currentConfig config.Configuration

//this is NEO part
func startConnectToSeed(c config.Configuration) {
	first := c.SeedList[0]
	host, port, err := net.SplitHostPort(first)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}
	var neoNodeConfig = neotx.Config{
		Network:   network.NEONetworkMagic(c.Magic),
		Port:      uint16(portInt),
		IPAddress: host,
	}
	client := neotx.NewClient(neoNodeConfig)
	handler := &NEOConnectionHandler{}
	handler.config = c
	handler.q = q.Shared()
	handler.Connected = false

	client.SetDelegate(handler)

	fmt.Printf("connecting to %v:%v...\n", neoNodeConfig.IPAddress, neoNodeConfig.Port)
	err = client.Start()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}
}

type NEOConnectionHandler struct {
	config    config.Configuration
	Connected bool
	version   network.Version
	q         *q.QService
}

type TransactionMessage struct {
	Type string      `json:"type"`
	TXID string      `json:"txID"`
	Data interface{} `json:"data,omitempty"`
}

//implement the message protocol
func (h *NEOConnectionHandler) OnReceive(tx neotx.TX) {

	if tx.Type == network.InventotyTypeTX {
		m := TransactionMessage{
			Type: tx.Type.String(),
			TXID: tx.ID,
		}
		sendMessage(h, m)
		return
	}
	if tx.Type == network.InventotyTypeConsensus {
		m := TransactionMessage{
			Type: tx.Type.String(),
			TXID: tx.ID,
		}
		sendMessage(h, m)
		return
	}
	//another type of INV. consensus and block
	m := TransactionMessage{
		Type: tx.Type.String(),
		TXID: tx.ID,
	}
	sendMessage(h, m)
}

func sendMessage(h *NEOConnectionHandler, m TransactionMessage) {
	log.Printf("%+v", m)

	source, err := models.JSONBFromObject(h.version)
	if err != nil {
		return
	}
	tx := models.TransactionLog{
		TXType: m.Type,
		Source: source,
		TXID:   m.TXID,
	}
	h.q.PushTransaction(tx)
}

func (h *NEOConnectionHandler) OnConnected(c network.Version) {
	fmt.Printf("connected %+v\n", c)
	h.Connected = true
	h.version = c
}

func (h *NEOConnectionHandler) OnError(e error) {
	if h.Connected == true {
		h.Connected = false
		fmt.Printf("Disconnected from host. will try to connect in 15 seconds...")
		for {
			time.Sleep(15 * time.Second)
			//we need to implement backoff and retry to reconnect here
			//if the error is EOF then we try to reconnect
			go startConnectToSeed(currentConfig)
		}
	}
}

func main() {

	mode := flag.String("network", "", "Network to connect to. main | test | private")
	flag.Parse()

	if *mode == "" {
		//default mode is private
		defaultEnv := "private"
		mode = &defaultEnv
	}

	file := "config.privatenet.json"
	if *mode == "main" {
		file = "config.json"
	} else if *mode == "test" {
		file = "config.testnet.json"
	}

	fmt.Printf("Loading config file:%v\n", file)

	c, err := config.LoadConfigurationFile(file)

	if err != nil {
		fmt.Printf("Error loading config file: %v", err)
		return
	}
	//assign the current configuration to global
	currentConfig = c

	var sync sync.WaitGroup
	sync.Add(1)
	go startConnectToSeed(currentConfig)
	sync.Wait()
}
