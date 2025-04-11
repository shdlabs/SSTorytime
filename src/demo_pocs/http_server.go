//
// Simple web server lookup
//

package main

import (
	"fmt"
	"net/http"
        SST "SSTorytime"
)

var CTX SST.PoSST

// *********************************************************************

func main() {

	CTX = SST.Open(true)	

	http.HandleFunc("/", Handler)
	fmt.Println("Listening at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// *********************************************************************

func Handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Credentials", "true")
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Add("Vary", "Origin")
	
	switch r.Method {

	case "GET":                // will be cached
		HandleGet(w,r)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func HandleGet(w http.ResponseWriter, r *http.Request) {

	// Get reply

	w.Header().Set("Content-Type", "application/json")

	nptrs := SST.GetDBNodePtrMatchingName(CTX,"","no")
	
	w.Write([]byte("{ \"matches\" : ["))
	
	for n := 0; n < len(nptrs); n++ {
		reply := []byte(SST.JSONNodeOrbit(CTX, nptrs[n]))
		w.Write(reply)
		if n != len(nptrs)-1 {
			w.Write([]byte(",\n"))
		}
	}

	w.Write([]byte("] }"))
}
























