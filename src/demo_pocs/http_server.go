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

	http.HandleFunc("/Orbit", OrbitHandler)
	http.HandleFunc("/Cone", ConeHandler)
	fmt.Println("Listening at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// *********************************************************************

func OrbitHandler(w http.ResponseWriter, r *http.Request) {

	GenHeader(w,r)

	fmt.Println("NCC Orbit response handler")
	
	switch r.Method {
	case "POST","GET":
		name := r.FormValue("name")
		chapter := r.FormValue("chapter")
		context := r.FormValue("context")
		HandleOrbit(w,r,name,chapter,context)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func HandleOrbit(w http.ResponseWriter, r *http.Request,name,chapter,context string) {

	// Get reply

	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Matching...NCC(",name,chapter,context,")")

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
	fmt.Println("Reply sent")
}

// *********************************************************************

func ConeHandler(w http.ResponseWriter, r *http.Request) {

	GenHeader(w,r)

	fmt.Println("NCC Orbit response handler")
	
	switch r.Method {
	case "POST","GET":
		name := r.FormValue("name")
		chapter := r.FormValue("chapter")
		context := r.FormValue("context")
		HandleCone(w,r,name,chapter,context)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func HandleCone(w http.ResponseWriter, r *http.Request,name,chapter,context string) {

	// Get reply

	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Matching...NCC(",name,chapter,context,")")

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
	fmt.Println("Reply sent")
}

// *********************************************************************

func GenHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Add("Vary", "Origin")
}























