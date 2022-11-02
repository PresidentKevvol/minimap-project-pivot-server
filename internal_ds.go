package main

import (
  "time"
  "fmt"
)

type EmptyContext struct {
}

//json payload from a beacon
type BeaconPayload struct {
  SourceName string
  Points []AccessPointInfo
}

type AccessPointInfo struct {
  SSID      string
  BSSID     string
  Channel   int32
  RSSI      float32
}

type BeaconRecord struct {
  RecordTime time.Time
  Points []AccessPointInfo
}

/*
type BeaconReadingsList struct {
  Capacity int
  Records []BeaconRecord
}

func (self *BeaconReadingsList) Push(br BeaconRecord) {
  // Push to the queue
  self.Records = append(self.Records, br)
  //fmt.Printf("records: %+v\n", self.Records)

  // Discard last element if capacity (of this beacon's list) reached
  if self.Capacity > 0 && len(self.Records) > self.Capacity {
    fmt.Printf("records removed.")
    self.Records = self.Records[1:]
  }
}
*/

type BeaconValuesDatabase struct {
  Capacity int
  //Bmap map[string]BeaconReadingsList
  Bmap map[string][]BeaconRecord
}

func (self *BeaconValuesDatabase) Push(name string, br BeaconRecord) {
  // create BeaconReadingsList object in the map if name not existed yet
  if _, ok := self.Bmap[name]; !ok {
    //self.Bmap[name] = BeaconReadingsList {Capacity: self.Capacity, Records: make([]BeaconRecord, 0)}
    self.Bmap[name] = make([]BeaconRecord, 0)
  }

  // Push to the respective list
  m, _ := self.Bmap[name]
  // m.Push(br)
  self.Bmap[name] = append(m, br)

  // Discard last element if capacity (of this beacon's list) reached
  m, _ = self.Bmap[name]
  if self.Capacity > 0 && len(m) > self.Capacity {
    fmt.Printf("records removed.")
    // self.Records = self.Records[1:]
    self.Bmap[name] = m[1:]
  }

  // fmt.Printf("self.Bmap[name]: %+v\n", m)
  fmt.Printf("Database: %+v\n", self)
}
