//
// Simple web server lookup
//

package main

import (
	"fmt"
	"net/http"
	"strings"
	"os"
	"encoding/json"

        SST "SSTorytime"
)

var CTX SST.PoSST

// *********************************************************************

func main() {

	CTX = SST.Open(true)	

	http.HandleFunc("/",PageHandler)
	http.HandleFunc("/Orbit", OrbitHandler)
	http.HandleFunc("/NPtrOrbit", OrbitHandler)
	http.HandleFunc("/Cone", ConeHandler)
	http.HandleFunc("/Browse", SystematicHandler)
	http.HandleFunc("/TOC", TableOfContents)
	http.HandleFunc("/Sequence", SequenceHandler)

	fmt.Println("Listening at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// *********************************************************************

func PageHandler(w http.ResponseWriter, r *http.Request) {

	GenHeader(w,r)

	switch r.Method {
	case "GET":

		w.Header().Set("Content-Type", "text/html")
		page,err := os.ReadFile("./page.html")

		if err != nil {
			fmt.Println("Can't find ./page.html")
			os.Exit(-1)
		}

		w.Write(page)

	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func OrbitHandler(w http.ResponseWriter, r *http.Request) {

	GenHeader(w,r)

	fmt.Println("NCC Orbit response handler")
	
	switch r.Method {
	case "POST","GET":
		nclass := r.FormValue("nclass")
		ncptr := r.FormValue("ncptr")
		chapter := r.FormValue("chapter")
		context := r.FormValue("context")
		name := r.FormValue("name")

		if nclass == "" || ncptr == "" {
			if name == "" {
				name = "semantic"
			}
			fmt.Println("Matching Orbit by name(",name,chapter,context,")")
			nptrs := SST.GetDBNodePtrMatchingName(CTX,name,chapter)
			HandleOrbit(w,r,nptrs,chapter,context)
		} else {
			fmt.Println("Matching Orbit by NPtr(",nclass,ncptr,chapter,context,")")
			var nptrs []SST.NodePtr
			var nptr SST.NodePtr
			fmt.Sscanf(nclass,"%d",&nptr.Class)
			fmt.Sscanf(ncptr,"%d",&nptr.CPtr)
			nptrs = append(nptrs,nptr)
			HandleOrbit(w,r,nptrs,chapter,context)
		}

	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func HandleOrbit(w http.ResponseWriter, r *http.Request,nptrs []SST.NodePtr,chapter,context string) {

	chapter = strings.TrimSpace(chapter)

	w.Header().Set("Content-Type", "application/json")

	events := fmt.Sprintf("{ \"events\" : [")
	
	for n := 0; n < len(nptrs); n++ {
		events += SST.JSONNodeEvent(CTX, nptrs[n])
		if n != len(nptrs)-1 {
			events += ",\n"
		}
	}
	
	events += "] }"
	
	w.Write([]byte(events))
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
		arrnames := r.FormValue("arrnames")
		HandleEntireCone(w,r,name,chapter,context,arrnames)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func HandleEntireCone(w http.ResponseWriter, r *http.Request,name,chapter,cntstr,arrstr string) {

	chapter = strings.TrimSpace(chapter)
	name = strings.TrimSpace(name)

	arrnames,_ := SST.Str2Array(arrstr)
	cntxt,_ := SST.Str2Array(cntstr)

	var arrows []SST.ArrowPtr

	for a := range arrnames {
		if len(arrnames[a]) > 1 {
			arr := SST.GetDBArrowByName(CTX,arrnames[a])
			arrows = append(arrows,arr)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Matching...EntireCone(",name,chapter,cntxt,arrows,")")

	nptrs := SST.GetDBNodePtrMatching(CTX,name,chapter,cntxt,arrows)

	maxdepth := 20
	var count int

	// Encode

	multicone  := "{ \"paths\" : [\n"

	for n := 0; n < len(nptrs); n++ {

		thiscone := fmt.Sprintf(" { \"NClass\" : %d,\n",nptrs[n].Class)
		thiscone += fmt.Sprintf("   \"NCPtr\" : %d,\n",nptrs[n].CPtr)
		thiscone += fmt.Sprintf("   \"Title\" : \"%s\",\n",name)
		empty := true

		cone,span := SST.GetEntireConePathsAsLinks(CTX,"any",nptrs[n],maxdepth)
		
		json := SST.JSONCone(CTX,cone,chapter,cntxt)
		
		if span > 0 {
			empty = false
		}
		
		thiscone += fmt.Sprintf("\"Entire\" : %s ",json)		
		thiscone += "\n}"

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

	arrnames,_ := SST.Str2Array(arrstr)
	context,_ := SST.Str2Array(cntstr)

	if section <= 0 {
		section = 1
	}

	fmt.Println("Matching...Browse(",section,chaptext,context,arrnames,")",len(arrnames))

	var arrows []SST.ArrowPtr

	for a := range arrnames {
		arr := SST.GetDBArrowByName(CTX,arrnames[a])
		if arr != 0 {
			arrows = append(arrows,arr)
		}
	}

	qnodes := SST.GetDBNodeContextsMatchingArrow(CTX,"",chaptext,context,arrows,section)

	w.Header().Set("Content-Type", "application/json")

	EncodeBrowsing(w,r,qnodes,arrows,section,chaptext,context)

	fmt.Printf("Reply Systematic Browser page %d sent\n",section)
}

//**************************************************************

func EncodeBrowsing(w http.ResponseWriter, r *http.Request,qnodes []SST.QNodePtr,arrows []SST.ArrowPtr,section int,chapter string,context []string) {

	// Policy for ordering and search depth along each vector

	order    := []int{0,1,-1,2,-2,3,-3}
	maxdepth := []int{2,8, 3,2, 2,3, 2}
	headerdone := false
	var multicone string
	var comma string

	// Encode

	fmt.Println("Looking for section",section)

	for q := range qnodes {

		if !headerdone {
			multicone += fmt.Sprintf("{ \"section\" : \"%d\",\n",section)
			multicone += fmt.Sprintf("  \"chapter\" : \"%s\",\n",qnodes[q].Chapter)
			multicone += fmt.Sprintf("  \"context\" : \"%v\",\n",CleanText(qnodes[q].Context))
			multicone += fmt.Sprintf("  \"nptrs\" : [ ")
			headerdone = true
		}
		
		s := SST.GetDBNodeByNodePtr(CTX,qnodes[q].NPtr).S
		thiscone := fmt.Sprintf("%s\n { \"NClass\" : %d,\n",comma,qnodes[q].NPtr.Class)
		thiscone += fmt.Sprintf(" \"NCPtr\" :%d,\n",qnodes[q].NPtr.CPtr)
		title,_ := json.Marshal(s)
		thiscone += fmt.Sprintf(" \"Title\" : %s,\n",string(title))
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

	multicone += "]\n}\n"
	w.Write([]byte(multicone))
}

// *********************************************************************

func TableOfContents(w http.ResponseWriter, r *http.Request) {

	GenHeader(w,r)

	fmt.Println("TableOfContents handler")

	switch r.Method {
	case "POST","GET":
		chapter := r.FormValue("chapter")
		cntstr := r.FormValue("context")
		context,_ := SST.Str2Array(cntstr)

		toc := SST.JSON_TableOfContents(CTX,chapter,context)
		w.Write([]byte(toc))
		fmt.Println(toc)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func SequenceHandler(w http.ResponseWriter, r *http.Request) {

        // Find a sequence of arrows matching arrname/default "then" for which
        // something in the orbit matches the search strings

	GenHeader(w,r)

	fmt.Println("Sequence search response handler")

	switch r.Method {
	case "POST","GET":
		name := r.FormValue("name")
		chapter := r.FormValue("chapter")
		context := r.FormValue("context")
		arrnames := r.FormValue("arrnames")
		HandleSequence(w,r,name,chapter,context,arrnames)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func HandleSequence(w http.ResponseWriter, r *http.Request,searchtext,chaptext string,cntstr,arrow string) {

	chapter := strings.TrimSpace(chaptext)
	context,_ := SST.Str2Array(cntstr)
	searchtext = strings.TrimSpace(searchtext)

	stories := SST.GetSequenceContainers(CTX,arrow,searchtext,chapter,context)
	orbits,_ := json.Marshal(stories)

        // returns story in events.Axis, with any container/title first in Story type

	story := fmt.Sprintf("{ \"events\" : %s }",orbits)

	w.Write([]byte(story))
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















