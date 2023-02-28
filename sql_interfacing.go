package main

import (
  //"fmt"
  //"os"
  //"strings"
  //"strconv"
  "log"
  //"math/rand"
  "time"
  //"context"
  //"github.com/go-redis/redis/v9"
  //"encoding/json"
	//"io/ioutil"
  //"github.com/gorilla/websocket"
  "database/sql"
	_ "github.com/lib/pq"
)

var sqldb *sql.DB

/*
initialize the setup for using SQL database in this server application
*/
func sqlInit(addr string, user string, pw string, dbName string) {
  // TODO: add SQL server connection creation
  sqldb, err := sql.Open("postgres", user + ":" + pw + "@tcp(" + addr + ")/" + dbName)
	if err != nil {
		log.Fatal(err)
	}

  //ping to create connection
  err = sqldb.Ping()
  if err != nil {
	   log.Fatal(err)
  }
}

/*
writes a BeaconRecord object to the relational database
*/
func storeBeaconRecord(BeaconName string, recordTime time.Time, br_marshalled string) {
  // TODO: add adding beacon reading data to SQL server
  //prepare the insert query
  var query = "INSERT INTO dataset_collection_beacon (beacon_name, record_Time, readings) VALUES (?, ?, ?);"
  stmt, prepareErr := sqldb.Prepare(query)
  if prepareErr != nil {
	   log.Fatal(prepareErr)
  }

  //execute the insert and note any errors
  _, execErr := stmt.Exec(BeaconName, recordTime, br_marshalled)
  if execErr != nil {
	   log.Fatal(execErr)
  }
}

/*
writes a data collection request to the relational database
*/
func storeCollectRecord(recordTime time.Time, pl_marshalled string) {
  // TODO: add adding point collection data to SQL server
  //prepare the insert query
  var query = "INSERT INTO dataset_collection_client (record_Time, readings) VALUES (?, ?);"
  stmt, prepareErr := sqldb.Prepare(query)
  if prepareErr != nil {
	   log.Fatal(prepareErr)
  }

  //execute the insert and note any errors
  _, execErr := stmt.Exec(recordTime, pl_marshalled)
  if execErr != nil {
	   log.Fatal(execErr)
  }
}
