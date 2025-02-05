package main

import (
  "fmt"
  "os"
  //"strings"
  "strconv"
  "log"
  //"math/rand"
  //"time"
  "net/http"
  //"path/filepath"
  //"html/template"
  //"encoding/json"
  "github.com/joho/godotenv"
	//"io/ioutil"
  //"github.com/gorilla/websocket"
)

//record the signal strengths value readings from the beacons
var beaconValues BeaconValuesDatabase
var BeaconValuesDBCapacity int

//check the beaconupdates sources' authenticity
var check_req_auth = false
// the beacon names ('usernames') and passwords
var auth_pw_map = make(map[string]string)

// where to store the collected fingerprint data
var fingerprint_data_storage string

//SQL mode to write collected data (beacon + point client) to SQL server as dataset
var sql_dataset_mode = false
var file_dataset_mode = false

func main() {
  //load the .env file
  err := godotenv.Load(workdir + "/.env")
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  fingerprint_data_storage = os.Getenv("FINGERPRINT_DATA_DESTINATION")

  //print the current working directory
  fmt.Println("directory of this executable: " + workdir)

  //define handler functions for pages
  http.HandleFunc("/", handleIndex)
  http.HandleFunc("/p/", handleBeaconUpdate)
  http.HandleFunc("/f/", handleFingerprintDataCollect)
  http.HandleFunc("/c/", handleCalculatePositionNaive)
  http.HandleFunc("/l/", handleUpdateLookup)
  http.HandleFunc("/ds/t/", handleDataCollectionToggle)

  //setup the auth passwords
  auth_pw_map["SBU-01"] = "cghj1A90tS3h7Msd"

  //setup the redis connection
  redisInit(os.Getenv("REDIS_IP"), os.Getenv("REDIS_PW"), false)
  fmt.Println("redis connection has been setup.")
  //setup the sql connection
  if len(os.Getenv("SQL_IP")) > 0 {
    sqlInit(os.Getenv("SQL_IP"), os.Getenv("SQL_USER"), os.Getenv("SQL_PW"), os.Getenv("SQL_DBNAME"))
    fmt.Println("postgresql connection has been setup.")
  }
  // setup the array's capacity value
  bvc := os.Getenv("BEACON_VALUES_CAPACITY")
  if len(bvc) == 0 {
    BeaconValuesDBCapacity = 6
  } else {
    i, atoi_err := strconv.Atoi(bvc)
    if atoi_err != nil {
      BeaconValuesDBCapacity = 6
    } else {
      BeaconValuesDBCapacity = i
    }
  }
  beaconValues = BeaconValuesDatabase {Capacity: BeaconValuesDBCapacity, Bmap: make(map[string][]BeaconRecord)}

  //the handler for ajax requests
  //http.HandleFunc("/ajax/createpost/", handleAjaxCreatePost)

  //for handling static css and js files
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(workdir + "/static/css"))))
  http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(workdir + "/static/js"))))
  // http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir(workdir + "/static/fonts"))))
  // http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(workdir + "/static/images"))))
  //setup the hostname env variable
  hostname := os.Getenv("host_name")
  if hostname == "" {
    hostname = ":8884"
  }
  //see if there is a ssl cert to be used
  ssl_cert := os.Getenv("ssl_cert")
  ssl_key := os.Getenv("ssl_key")
  if ssl_cert != "" && ssl_key != "" {
    //start the server program with ssl
    log.Fatal(http.ListenAndServeTLS(hostname, ssl_cert, ssl_key, nil))
  } else {
    //start the server program
    log.Fatal(http.ListenAndServe(hostname, nil))
  }
}
