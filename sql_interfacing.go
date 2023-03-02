package main

import (
  "fmt"
  //"os"
  //"strings"
  //"strconv"
  "log"
  //"math/rand"
  "time"
  //"context"
  //"github.com/go-redis/redis/v9"
  "encoding/json"
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
  //sqldb, err := sql.Open("postgres", user + ":" + pw + "@tcp(" + addr + ")/" + dbName)
  db, err := sql.Open("postgres", "user=" + user + " password=" + pw + " host=" + addr + " dbname=" + dbName)
	if err != nil {
		log.Fatal(err)
	}
  sqldb = db

  //ping to create connection
  err = sqldb.Ping()
  if err != nil {
	   log.Fatal(err)
  }
}

/*
writes a BeaconRecord object to the relational database
*/
func storeBeaconRecord(BeaconName string, recordTime time.Time, points []AccessPointInfo) {
  // TODO: add adding beacon reading data to SQL server
  p, merr := json.Marshal(points)
  if merr != nil {
    fmt.Println(merr)
  }
  var points_marshalled = string(p)

  //prepare the insert query
  var query = `INSERT INTO dataset_collection_beacon (beacon_name, record_time, points) VALUES ($1, $2, $3)`
  //execute with exec function
  _, e := sqldb.Exec(query, BeaconName, recordTime, points_marshalled)
  if e != nil {
	   fmt.Println(e)
  }
}

/*
writes a data collection request to the relational database
*/
func storeCollectRecord(recordTime time.Time, pl FingerprintDataCollectPayload) {
  // TODO: add adding point collection data to SQL server
  p, merr := json.Marshal(pl.Points)
  if merr != nil {
    fmt.Println(merr)
  }
  var pts_marshalled = string(p)

  //prepare the insert query
  var query = `INSERT INTO dataset_collection_client (source_id, record_time, points, spatial_id, note) VALUES ($1, $2, $3, $4, $5);`
  //execute with exec function
  _, e := sqldb.Exec(query, pl.SourceDeviceId, recordTime, pts_marshalled, pl.SpatialId, pl.Note)
  if e != nil {
     fmt.Println(e)
  }
}
