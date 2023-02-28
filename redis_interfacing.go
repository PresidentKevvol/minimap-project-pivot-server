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

var ctx = context.Background()
var rdb *redis.Client

/*
initialize the setup for using redis in this server application
*/
func redisInit(addr string, pw string, reset bool) {
  rdb = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: pw, // empty string means no password set
        DB:       0,  // use default DB
    })
}

/*
below are for all the CRUD functionalities with the redis db
the redis should have the following entries

BeaconNames: redis set with all the beacon's names

BeaconRecord-<name>: redis list for the records received from a beacon, one for each beacon
*/

/*
writes a BeaconRecord object to the redis base
*/
func writeBeaconRecord(name string, br_marshalled string) {
  //check if name exist
  nameExist, err := rdb.SIsMember(ctx, "BeaconNames", name).Result()
  if err != nil {
    fmt.Println(err)
    return
  }
  //if not, add to the set of names
  if !nameExist {
    _, err = rdb.SAdd(ctx, "BeaconNames", name).Result()
    if err != nil {
      fmt.Println(err)
      return
    }
  }

  _, err = rdb.LPush(ctx, "BeaconRecord-" + name, br_marshalled).Result()
  if err != nil {
    fmt.Println(err)
    return
  }
  //trim list as needed
  _, err = rdb.LTrim(ctx, "BeaconRecord-" + name, 0, int64(BeaconValuesDBCapacity)).Result()
  if err != nil {
    fmt.Println(err)
    return
  }
}

func getAllBeaconRedingsRedis() BeaconValuesDatabase {
  var res = BeaconValuesDatabase{Capacity: BeaconValuesDBCapacity, Bmap: make(map[string][]BeaconRecord)}
  // run the scan command to get all beacon name strings
  keys, _, sscanErr := rdb.SScan(ctx, "BeaconNames", 0, "*", 99).Result()
  if sscanErr != nil {
    fmt.Println(sscanErr)
  }

  // get item from each key
  for _, key := range keys {
    cur_readings, lrangeErr := rdb.LRange(ctx, "BeaconRecord-" + key, 0, int64(BeaconValuesDBCapacity)+1).Result()
    if lrangeErr != nil {
      fmt.Println(lrangeErr)
    }

    var cur_readings_objs []BeaconRecord
    for _, reading := range cur_readings {
      //decode the json string stored in the redis base
      var content_json BeaconRecord
      decodeErr := json.Unmarshal([]byte(reading), &content_json)
      if decodeErr != nil {
        fmt.Println(decodeErr)
      }

      //add each reading to slice
      cur_readings_objs = append(cur_readings_objs, content_json)
    }
    res.Bmap[key] = cur_readings_objs
  }
  return res
}
