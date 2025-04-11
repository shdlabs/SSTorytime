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
	case "POST":
		name := r.FormValue("name")
		chapter := r.FormValue("chapter")
		context := r.FormValue("context")
		HandleGet(w,r,name,chapter,context)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func HandleGet(w http.ResponseWriter, r *http.Request,name,chapter,context string) {

	// Get reply

	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Matching...",name,chapter,context)

	if name == "" {
		name = "fish"
	}

	nptrs := SST.GetDBNodePtrMatchingName(CTX,chapter,name)
	
	w.Write([]byte("{ \"matches\" : ["))
	
	for n := 0; n < len(nptrs); n++ {
		reply := []byte(SST.JSONNodeOrbit(CTX, nptrs[n]))
		w.Write(reply)
		if n != len(nptrs)-1 {
			w.Write([]byte(",\n"))
		}
	}

	w.Write([]byte("] }"))
	fmt.Println("Replace sent")
}
























