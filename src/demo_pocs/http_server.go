//
// Simple web server lookup
//

package main

import (
	"fmt"
	"net/http"
	"strings"

        SST "SSTorytime"
)

var CTX SST.PoSST

// *********************************************************************

func main() {

	CTX = SST.Open(true)	

	http.HandleFunc("/Orbit", OrbitHandler)
	http.HandleFunc("/Cone", ConeHandler)
	http.HandleFunc("/Browse", SystematicHandler)
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

	chapter = strings.TrimSpace(chapter)

	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Matching OrbitNCC(",name,chapter,context,")")

	if name == "" {
		name = "fish"
	}
	
	nptrs := SST.GetDBNodePtrMatchingName(CTX,chapter,name)
	
	orbit := fmt.Sprintf("{ \"matches\" : [")
	
	for n := 0; n < len(nptrs); n++ {
		orbit += SST.JSONNodeOrbit(CTX, nptrs[n])
		if n != len(nptrs)-1 {
			orbit += ",\n"
		}
	}
	
	orbit += "] }"
	
	w.Write([]byte(orbit))
	fmt.Println(orbit)
	fmt.Println("Reply Orbit sent")
}

// *********************************************************************

func ConeHandler(w http.ResponseWriter, r *http.Request) {

	GenHeader(w,r)

	fmt.Println("NCC Cone response handler")
	
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

	chapter = strings.TrimSpace(chapter)
	name = strings.TrimSpace(name)

	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Matching...ConeNCC(",name,chapter,context,")")

	if name == "" {
		name = "lamb"
	}

	nptrs := SST.GetDBNodePtrMatchingName(CTX,chapter,name)
	cntxt := strings.Split(context," ")

	// Policy for ordering and search depth along each vector

	order    := []int{0,1,-1,2,-2,3,-3}
	maxdepth := []int{2,4,4,2,2,2,2}
	var count int

	// Encode

	multicone  := "{ \"paths\" : [\n"

	for n := 0; n < len(nptrs); n++ {

		thiscone := fmt.Sprintf(" { \"NPtr\" : \"%v\",\n",nptrs[n])
		thiscone += fmt.Sprintf("   \"Text\" : \"%s\",\n",name)
		empty := true

		for i := range order {
			sttype := order[i]
			cone,span := SST.GetFwdPathsAsLinks(CTX,nptrs[n],sttype,maxdepth[i])
			json := SST.JSONCone(CTX,cone,chapter,cntxt)

			if span > 0 {
				empty = false
			}

			thiscone += fmt.Sprintf("\"%s\" : %s ",SST.STTypeDBChannel(sttype),json)

			if i < len(order)-1 {
				thiscone += ",\n"
			} else {
				thiscone += "\n}"
			}
		}

		if !empty {
			if count > 0 {
				thiscone = "\n,"+thiscone
			}
			multicone += thiscone
			count++
		}
	}

	multicone += "]\n}\n"

	w.Write([]byte(multicone))
	fmt.Println(multicone)
	fmt.Println("Reply Cone sent")
}

//******************************************************************

func SystematicHandler(w http.ResponseWriter, r *http.Request) {

	GenHeader(w,r)

	fmt.Println("Browse response handler")
	var secnr int = 1

	switch r.Method {
	case "POST","GET":
		arrnames := r.FormValue("arrnames")
		chapter := r.FormValue("chapter")
		context := r.FormValue("context")
		pg := r.FormValue("pagenr")
		fmt.Sscanf(pg,"%d",&secnr)
		HandleSystematic(w,r,secnr,chapter,context,arrnames)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

//******************************************************************

func HandleSystematic(w http.ResponseWriter, r *http.Request,section int,chaptext string,cntstr,arrstr string) {

	chaptext = strings.TrimSpace(chaptext)

	arrnames,_ := Str2Array(arrstr)
	context,_ := Str2Array(cntstr)

	if section <= 0 {
		section = 1
	}

	fmt.Println("Matching...Browse(",section,chaptext,context,arrnames,")",len(arrnames))

	var arrows []SST.ArrowPtr

	for a := range arrnames {
		arr := SST.GetDBArrowByName(CTX,arrnames[a])
		arrows = append(arrows,arr)
	}

	qnodes := SST.GetDBNodeContextsMatchingArrow(CTX,chaptext,context,"",arrows)

	w.Header().Set("Content-Type", "application/json")

	EncodeBrowsing(w,r,qnodes,arrows,section,chaptext,context)

	fmt.Printf("Reply Systematic Browser page %d sent\n",section)
}

//**************************************************************

func EncodeBrowsing(w http.ResponseWriter, r *http.Request,qnodes []SST.QNodePtr,arrows []SST.ArrowPtr,section int,chapter string,context []string) {

	// Policy for ordering and search depth along each vector

	order    := []int{0,1,-1,2,-2,3,-3}
	maxdepth := []int{2,4,4,2,2,2,2}
	headerdone := false
	var secnr int
	var prev string = ""
	var multicone string
	var comma string

	// Encode

	fmt.Println("Looking for section",section)

	for q := range qnodes {

		if qnodes[q].Context != prev {
			prev = qnodes[q].Context
			secnr++
			
			if secnr > section {
				multicone += "]\n}\n"
				w.Write([]byte(multicone))
				fmt.Println(multicone)
				return
				}
		}
		
		if secnr == section {
			if !headerdone {
				multicone += fmt.Sprintf("{ \"section\" : \"%d\",\n",section)
				multicone += fmt.Sprintf("  \"chapter\" : \"%s\",\n",qnodes[q].Chapter)
				multicone += fmt.Sprintf("  \"context\" : \"%v\",\n",CleanText(qnodes[q].Context))
				multicone += fmt.Sprintf("  \"nptrs\" : [ ")
				headerdone = true
			}
			
			thiscone := fmt.Sprintf("%s\n { \"NPtr\" : \"%v\",\n",comma,qnodes[q].NPtr)
			comma = ","
			
			for i := range order {
				sttype := order[i]
				cone,_ := SST.GetFwdPathsAsLinks(CTX,qnodes[q].NPtr,sttype,maxdepth[i])
				json := SST.JSONCone(CTX,cone,chapter,context)
				thiscone += fmt.Sprintf("\"%s\" : %s ",SST.STTypeDBChannel(sttype),json)
				
				if i < len(order)-1 {
					thiscone += ",\n"
				} else {
					thiscone += "}"
				}
			}
			
			multicone += thiscone
		}		
	}
}

// *********************************************************************

func GenHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Add("Vary", "Origin")
}

// *********************************************************************

func CleanText(c string) string {

	c = strings.Replace(c,"{","",-1)
	c = strings.Replace(c,"}","",-1)
	c = strings.Replace(c,","," ",-1)
	c = strings.Replace(c,"\"","\\\"",-1)
	return c
}

// *********************************************************************

func Str2Array(s string) ([]string,int) {

	var non_zero int

	s = strings.Replace(s,"{","",-1)
	s = strings.Replace(s,"}","",-1)
	s = strings.Replace(s,"\"","",-1)

	arr := strings.Split(s,",")

	for a := 0; a < len(arr); a++ {
		arr[a] = strings.TrimSpace(arr[a])
		if len(arr[a]) > 0 {
			non_zero++
		}
	}

	return arr,non_zero
}





















