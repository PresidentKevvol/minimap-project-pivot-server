package main

import (
  "fmt"
  "os"
  //"strings"
  //"strconv"
  "log"
  //"math/rand"
  "time"
  "net/http"
  "path/filepath"
  "html/template"
  "encoding/json"
  "github.com/joho/godotenv"
	//"io/ioutil"
  //"github.com/gorilla/websocket"
)

//record the signal strengths value readings from the beacons
var beaconValues BeaconValuesDatabase = BeaconValuesDatabase {Capacity: 6, Bmap: make(map[string][]BeaconRecord)}

//check the beaconupdates sources' authenticity
var check_req_auth = true
// the beacon names ('usernames') and passwords
var auth_pw_map = make(map[string]string)

//template for the pages
var ex, exerr = os.Executable()
var workdir = filepath.Dir(ex)
var page_templates = template.Must(template.ParseFiles(
  workdir + "/views/index.html",
  workdir + "/views/updateinfo.html",
  workdir + "/views/lookup.html"))

//the handler for the index page
func handleIndex(w http.ResponseWriter, r *http.Request) {
  //render template
  err := page_templates.ExecuteTemplate(w, "index.html", EmptyContext {})
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// the handler for when a beacon sends an update
func handleBeaconUpdate(w http.ResponseWriter, r *http.Request) {
  if check_req_auth {
    // grab the http header
    http_header := r.Header
    beacon_name, name_exist := http_header["Beacon-Name"]
    beacon_pw, pw_exist := http_header["Beacon-Password"]
    if name_exist == false || pw_exist == false { //if there is an error
      http.Error(w, "", http.StatusInternalServerError)
      return
    }
    check_pw, chk_exist := auth_pw_map[beacon_name[0]]
    if chk_exist == false { //if beacon name not exist
      http.Error(w, "", http.StatusInternalServerError)
      return
    }
    if check_pw != beacon_pw[0] { //finally check the password is right
      http.Error(w, "", http.StatusInternalServerError)
      return
    }
  }

  //get the http POST parameter from the http request
  decoder := json.NewDecoder(r.Body)
  var content_json BeaconPayload
  err := decoder.Decode(&content_json)
  if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
  }
  // fmt.Printf("content_json.SourceName = %s\n", content_json.SourceName)
  // fmt.Printf("content_json.Numbers[1] = %f\n", content_json.Points[1])
  // address := r.FormValue("address")
	// fmt.Fprintf(w, "Address = %s\n", address)

  // TODO: record reported information of signal strengths
  dt := time.Now()   //record current time
  record := BeaconRecord {RecordTime: dt, Points: content_json.Points}
  beaconValues.Push(content_json.SourceName, record)

  //render template
  err = page_templates.ExecuteTemplate(w, "updateinfo.html", EmptyContext {})
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// the handler for displaying the debug page with lookup for the values
func handleUpdateLookup(w http.ResponseWriter, r *http.Request) {
  //render template
  err := page_templates.ExecuteTemplate(w, "lookup.html", beaconValues)
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func main() {
  //load the .env file
  err := godotenv.Load(workdir + "/.env")
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  //print the current working directory
  fmt.Println("directory of this executable: " + workdir)

  //define handler functions for pages
  http.HandleFunc("/", handleIndex)
  http.HandleFunc("/p/", handleBeaconUpdate)
  http.HandleFunc("/l/", handleUpdateLookup)

  //setup the auth passwords
  auth_pw_map["SBU-01"] = "cghj1A90tS3h7Msd"

  //the handler for ajax requests
  //http.HandleFunc("/ajax/createpost/", handleAjaxCreatePost)

  //for handling static css and js files
  // http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(workdir + "/static/css"))))
  // http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(workdir + "/static/js"))))
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
