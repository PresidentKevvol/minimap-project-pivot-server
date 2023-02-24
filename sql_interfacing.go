package main

import (
  "fmt"
  //"os"
  //"strings"
  //"strconv"
  //"log"
  //"math/rand"
  //"time"
  "context"
  "github.com/go-redis/redis/v9"
  "encoding/json"
	//"io/ioutil"
  //"github.com/gorilla/websocket"
)

/*
initialize the setup for using SQL database in this server application
*/
func sqlInit(addr string, pw string) {
  // TODO: add SQL server connection creation
}

/*
writes a BeaconRecord object to the relational database
*/
func storeBeaconRecord(name string, br BeaconRecord) {
  // TODO: add adding beacon reading data to SQL server
}

/*
writes a data collection request to the relational database
*/
func storeCollectRecord(pl FingerprintDataCollectPayload) {
  // TODO: add adding point collection data to SQL server
}
