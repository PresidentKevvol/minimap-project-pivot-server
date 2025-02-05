package main

import (
  "fmt"
  "os"
  //"strings"
  "strconv"
  //"log"
  //"math/rand"
  "time"
  "net/http"
  "path/filepath"
  "html/template"
  "encoding/json"
	"io/ioutil"
  //"github.com/gorilla/websocket"
)

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
  //do the security checking if authentication is required
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

  //now marshall the BeaconRecord to be fed to functions
  //the converting and converting back is crucial in validating the data format/object schema
  // b, merr := json.Marshal(content_json)
  // if merr != nil {
  //   fmt.Println(merr)
  //   http.Error(w, merr.Error(), http.StatusInternalServerError)
  // }
  // var br = string(b)

  raw_marshalled, merr2 := json.Marshal(record)
  if merr2 != nil {
    fmt.Println(merr2)
    http.Error(w, merr2.Error(), http.StatusInternalServerError)
  }
  var rc = string(raw_marshalled)

  //record to redis server too
  writeBeaconRecord(content_json.SourceName, rc)

  //store to sql server
  if sql_dataset_mode {
    storeBeaconRecord(content_json.SourceName, dt, content_json.Points)
  }

  //render template
  err = page_templates.ExecuteTemplate(w, "updateinfo.html", EmptyContext {})
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// the handler for displaying the debug page with lookup for the values
func handleUpdateLookup(w http.ResponseWriter, r *http.Request) {
  //get the values from the redis base
  var bValuesRedis = getAllBeaconRedingsRedis()
  //render template
  //err := page_templates.ExecuteTemplate(w, "lookup.html", beaconValues)
  err := page_templates.ExecuteTemplate(w, "lookup.html", bValuesRedis)

  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// the handler for data collection for fingerprinting
func handleFingerprintDataCollect(w http.ResponseWriter, r *http.Request) {
  //get the http POST parameter from the http request
  decoder := json.NewDecoder(r.Body)
  var content_json FingerprintDataCollectPayload
  err := decoder.Decode(&content_json)
  if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
  }

  dt := time.Now()   //record current time

  if file_dataset_mode {
    // p, merr := json.Marshal(content_json)
    // if merr != nil {
    //   fmt.Println(merr)
    //   http.Error(w, merr.Error(), http.StatusInternalServerError)
    // }
    // var pl = string(p)

    // retrieve latest reading of each beacon
    brs := make(map[string][]AccessPointInfo)
    for k, v := range beaconValues.Bmap {
      pts := v[0].Points
      brs[k] = pts
    }

    res := FingerprintDataPoint {
      SourceDeviceId    : content_json.SourceDeviceId,
      SourceReadings    : content_json.Points,
      SpatialId         : content_json.SpatialId,
      BeaconReadings    : brs,
    }

    filename := dt.Format("2006_01_02_15_04_05_000") + "-" + content_json.SourceDeviceId + ".json"

    // write the point to a file
    obj, err := json.Marshal(res)
    //obj, err := json.MarshalIndent(res, "", " ")
    if err != nil { //if there is an error
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    err = ioutil.WriteFile(fingerprint_data_storage + filename, obj, 0644)
    if err != nil { //if there is an error
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  }

  //store to sql server
  if sql_dataset_mode {
    storeCollectRecord(dt, content_json)
  }

  //render template
  err = page_templates.ExecuteTemplate(w, "updateinfo.html", EmptyContext {})
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// the handler for calculating the position with naive method
func handleCalculatePositionNaive(w http.ResponseWriter, r *http.Request) {
  //get the http POST parameter from the http request
  decoder := json.NewDecoder(r.Body)
  var content_json PositioningRequestPayload
  err := decoder.Decode(&content_json)
  if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
  }
}

func handleDataCollectionToggle(w http.ResponseWriter, r *http.Request)  {
  sql_dataset_mode = !sql_dataset_mode

  fmt.Fprintf(w, "sql_dataset_mode toggled.\ncurrent value:" + strconv.FormatBool(sql_dataset_mode))
}
