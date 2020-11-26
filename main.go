/*
Copyright (C) 2020  Zach Strauss

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, version 3.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type designWithAddress struct {
	Design   Design `json:"design"`
	Hostname string `json:"hostname"`
}

func main() {
	// Cores force HTTP with a default selfsigned cert, so we have to ignore verification issues connecting to it
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	address := flag.String("address", "", "Network address of the core")
	port := flag.Int64("port", 8080, "Port to start HTTP server on")

	flag.Parse()

	if *address == "" {
		flag.Usage()
		os.Exit(1)
	}

	http.HandleFunc("/designData.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		design, readErr := GetDesign(*address)
		if readErr != nil {
			http.Error(w, readErr.Error(), 500)
			log.Printf("Could not get design from %s: %s", *address, readErr)
			return
		}
		taggedDesign := designWithAddress{
			Design:   *design,
			Hostname: *address,
		}
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		jsonErr := encoder.Encode(taggedDesign)
		if jsonErr != nil {
			log.Printf("Error while encoding JSON response: %s", jsonErr)
		}
		log.Printf("Proxied request for design '%s' (%s)", taggedDesign.Design.DesignName, taggedDesign.Design.CompileGUID)
	})

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/", fs)

	log.Printf("Starting HTTP listener on port %d\n", *port)
	log.Printf("Proxying requests to Core at %s", *address)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
