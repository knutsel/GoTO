/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

//main.go
//takes in username, password, and database arguments
//runs server that handles url table searches

package main

import (
	"./sqlParser"
	"./structFilter"
	"./structs"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	//	"net/url"
	"./urlParser"
	"os"
)

var (
	addr        = flag.Bool("addr", false, "find open address and print to final-port.txt")
	username    = os.Args[1]
	password    = os.Args[2]
	environment = os.Args[3]
	db          = sqlParser.ConnectToDatabase(username, password, environment)
)

//returns JSON of argument table name in database
func generateHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	request := urlParser.ParseURL(path)

	tableName := request.TableName
	fields := request.Fields

	//only valid names are existing tables in db
	if !structs.ValidStruct[tableName] {
		fmt.Printf("\"%s\" table not found.\n", tableName)

		http.NotFound(w, r)
	} else {
		fmt.Printf("\"%s\" table found.\n", tableName)

		rows := sqlParser.GetRows(tableName, fields)
		w.Header().Set("Content-Type", "application/json")

		if fields == "" {
			fmt.Printf("No fields\n")
			structs.MapTableToJson(tableName, rows, w)
		} else {
			fmt.Printf("%s\n", fields)
			fieldArray := strings.Split(fields, ",")
			structFilter.MapCustomTableToJson(tableName, rows, w, fieldArray)
		}
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/", generateHandler)

	if *addr {
		//runs on home
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(":8080", nil)
}
