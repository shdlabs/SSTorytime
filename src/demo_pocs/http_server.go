//
// Simple web server lookup
//

package main

import (
	"fmt"
	"net/http"
	"strings"
	"os"

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
				name = "fish"
			}
			fmt.Println("Matching Orbit by name(",name,chapter,context,")")
			nptrs := SST.GetDBNodePtrMatchingName(CTX,chapter,name)
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

	arrnames,_ := Str2Array(arrstr)
	cntxt,_ := Str2Array(cntstr)

	var arrows []SST.ArrowPtr

	for a := range arrnames {
		if len(arrnames[a]) > 1 {
			arr := SST.GetDBArrowByName(CTX,arrnames[a])
			arrows = append(arrows,arr)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Matching...EntireCone(",name,chapter,cntxt,arrows,")")

	nptrs := SST.GetDBNodePtrMatching(CTX,chapter,name,cntxt,arrows)

	maxdepth := 20
	var count int

	// Encode

	multicone  := "{ \"paths\" : [\n"

	for n := 0; n < len(nptrs); n++ {

		thiscone := fmt.Sprintf(" { \"NPtr\" : \"%v\",\n",nptrs[n])
		thiscone += fmt.Sprintf("   \"Text\" : \"%s\",\n",name)
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
		if arr != 0 {
			arrows = append(arrows,arr)
		}
	}

	qnodes := SST.GetDBNodeContextsMatchingArrow(CTX,chaptext,context,"",arrows,section)

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

	multicone += "]\n}\n"
	w.Write([]byte(multicone))
	fmt.Println(multicone)
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





















