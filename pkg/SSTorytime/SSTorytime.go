//**************************************************************
//
// An interface for postgres for graph analytics and semantics
//
//**************************************************************

package SSTorytime

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	//"sort"

	_ "github.com/lib/pq"

)

//**************************************************************
// Errors
//**************************************************************

const (
	ERR_ST_OUT_OF_BOUNDS="Link STtype is out of bounds (must be -3 to +3)"
	ERR_ILLEGAL_LINK_CLASS="ILLEGAL LINK CLASS"
	ERR_NO_SUCH_ARROW = "No such arrow has been declared in the configuration: "
	ERR_MEMORY_DB_ARROW_MISMATCH = "Arrows in database are not in synch (shouldn't happen)"

	NEAR = 0
	LEADSTO = 1   // +/-
	CONTAINS = 2  // +/-
	EXPRESS = 3   // +/-

	// And shifted indices for array indicesin Go

	ST_ZERO = EXPRESS
	ST_TOP = ST_ZERO + EXPRESS + 1

	// For the SQL table, as 2d arrays not good

	I_MEXPR = "Im3"
	I_MCONT = "Im2"
	I_MLEAD = "Im1"
	I_NEAR  = "In0"
	I_PLEAD = "Il1"
	I_PCONT = "Ic2"
	I_PEXPR = "Ie3"

	// For separating text types

	N1GRAM = 1
	N2GRAM = 2
	N3GRAM = 3
	LT128 = 4
	LT1024 = 5
	GT1024 = 6
)

//**************************************************************

type Node struct {

	L         int     // length of text string
	S         string  // text string itself

	Chap      string  // section/chapter name in which this was added
	NPtr      NodePtr // Pointer to self index

	I [ST_TOP][]Link  // link incidence list, by arrow type
  	                  // NOTE: carefully how offsets represent negative SSTtypes
}

//**************************************************************

type NodeArrowNode struct {

	NFrom NodePtr
	STType int
	Arr ArrowPtr
	Wgt float64
	Ctx []string
	NTo NodePtr
  	                  // NOTE: carefully how offsets represent negative SSTtypes
}

//**************************************************************

type STTypeMatroid struct {

	NFrom NodePtr
	STType int
	NTo []NodePtr
}

//**************************************************************

type ArrowMatroid struct {

	NFrom NodePtr
	Arr ArrowPtr
	NTo []NodePtr
}

//**************************************************************

type Link struct {  // A link is a type of arrow, with context
                    // and maybe with a weightfor package math
	Arr ArrowPtr         // type of arrow, presorted
	Wgt float64          // numerical weight of this link
	Ctx []string         // context for this pathway
	Dst NodePtr          // adjacent event/item/node
}

const NODEPTR_TYPE = "CREATE TYPE NodePtr AS  " +
	"(                    " +
	"Chan     int,        " +
	"CPtr     int         " +
	")"

const LINK_TYPE = "CREATE TYPE Link AS  " +
	"(                    " +
	"Arr      int,        " +
	"Wgt      real,       " +
	"Ctx      text,       " +
	"Dst      NodePtr     " +
	")"

const NODE_TABLE = "CREATE TABLE IF NOT EXISTS Node " +
	"( " +
	"NPtr      NodePtr,        " +
	"L         int,            " +
	"S         text,           " +
	"Chap      text,           " +
	I_MEXPR+"  Link[],         " + // Im3
	I_MCONT+"  Link[],         " + // Im2
	I_MLEAD+"  Link[],         " + // Im1
	I_NEAR +"  Link[],         " + // In0
	I_PLEAD+"  Link[],         " + // Il1
	I_PCONT+"  Link[],         " + // Ic2
	I_PEXPR+"  Link[]          " + // Ie3
	")"

const LINK_TABLE = "CREATE TABLE IF NOT EXISTS NodeArrowNode " +
	"( " +
	"NFrom    NodePtr, " +
	"STtype   int,     " +
	"Arr      int,     " +
	"Wgt      int,     " +
	"Ctx      text[],  " +
	"NTo      NodePtr  " +
	")"

//**************************************************************

type NodeDirectory struct {

	// Power law n-gram frequencies

	N1grams map[string]ClassedNodePtr
	N1directory []Node
	N1_top ClassedNodePtr

	N2grams map[string]ClassedNodePtr
	N2directory []Node
	N2_top ClassedNodePtr

	N3grams map[string]ClassedNodePtr
	N3directory []Node
	N3_top ClassedNodePtr

	// Use linear search on these exp fewer long strings

	LT128 []Node
	LT128_top ClassedNodePtr
	LT1024 []Node
	LT1024_top ClassedNodePtr
	GT1024 []Node
	GT1024_top ClassedNodePtr
}

//**************************************************************

var NODE_CACHE = make(map[NodePtr]NodePtr)

//**************************************************************

type NodePtr struct {

	Class int            // Text size-class
	CPtr  ClassedNodePtr // index of within name class lane
}

type ClassedNodePtr int  // Internal pointer type of size-classified text

//**************************************************************

type ArrowDirectory struct {

	STAindex  int
	Long    string
	Short   string
	Ptr     ArrowPtr
}

type ArrowPtr int // ArrowDirectory index

const ARROW_DIRECTORY_TABLE = "CREATE TABLE IF NOT EXISTS ArrowDirectory " +
	"(    " +
	"STAindex int,             " +
	"Long text,              " +
	"Short text,             " +
	"ArrPtr int primary key  " +
	")"

const ARROW_INVERSES_TABLE = "CREATE TABLE IF NOT EXISTS ArrowInverses " +
	"(    " +
	"Plus int,  " +
	"Minus int,  " +
	"Primary Key(Plus,Minus)" +
	")"

//**************************************************************
// Lookup tables
//**************************************************************

var ( 
	ARROW_DIRECTORY []ArrowDirectory
	ARROW_SHORT_DIR = make(map[string]ArrowPtr) // Look up short name int referene
	ARROW_LONG_DIR = make(map[string]ArrowPtr)  // Look up long name int referene
	ARROW_DIRECTORY_TOP ArrowPtr = 0
	INVERSE_ARROWS = make(map[ArrowPtr]ArrowPtr)

	NODE_DIRECTORY NodeDirectory  // Internal histo-representations

	NO_NODE_PTR NodePtr // see Init()

	WIPE_DB bool = false

)

//******************************************************************
// LIBRARY
//******************************************************************

type PoSST struct {

   DB *sql.DB
}

//******************************************************************

func Open(load_arrows bool) PoSST {

	var ctx PoSST
	var err error

	const (
		host     = "localhost"
		port     = 5432
		user     = "sstoryline"
		password = "sst_1234"
		dbname   = "newdb"
	)

        connStr := "user="+user+" dbname="+dbname+" password="+password+" sslmode=disable"

        ctx.DB, err = sql.Open("postgres", connStr)

	if err != nil {
	   	fmt.Println("Error connecting to the database: ", err)
		os.Exit(-1)
	}
	
	err = ctx.DB.Ping()
	
	if err != nil {
		fmt.Println("Error pinging the database: ", err)
		os.Exit(-1)
	}

	MemoryInit()
	Configure(ctx,load_arrows)

	NO_NODE_PTR.Class = 0
	NO_NODE_PTR.CPtr =  -1

	return ctx
}

// **************************************************************************

func MemoryInit() {

	if NODE_DIRECTORY.N1grams == nil {
		NODE_DIRECTORY.N1grams = make(map[string]ClassedNodePtr)
	}

	if NODE_DIRECTORY.N2grams == nil {
		NODE_DIRECTORY.N2grams = make(map[string]ClassedNodePtr)
	}

	if NODE_DIRECTORY.N3grams == nil {
		NODE_DIRECTORY.N3grams = make(map[string]ClassedNodePtr)
	}
}

// **************************************************************************

func Configure(ctx PoSST,load_arrows bool) {

	// Tmp reset

	if WIPE_DB {

		fmt.Println("***********************")
		fmt.Println("* WIPING DB")
		fmt.Println("***********************")
		
		ctx.DB.QueryRow("drop function fwdconeaslinks")
		ctx.DB.QueryRow("drop function fwdconeasnodes")
		ctx.DB.QueryRow("drop function fwdpathsaslinks")
		ctx.DB.QueryRow("drop function getfwdlinks")
		ctx.DB.QueryRow("drop function getfwdnodes")
		ctx.DB.QueryRow("drop function getneighboursbytype")
		ctx.DB.QueryRow("drop function getsingletonaslink")
		ctx.DB.QueryRow("drop function getsingletonaslinkarray")
		ctx.DB.QueryRow("drop function idempinsertnode")
		ctx.DB.QueryRow("drop function sumfwdpaths")
		ctx.DB.QueryRow("drop table Node")
		ctx.DB.QueryRow("drop table NodeArrowNode")
		ctx.DB.QueryRow("drop type NodePtr")
		ctx.DB.QueryRow("drop type Link")

		ctx.DB.QueryRow("drop table ArrowDirectory")
		ctx.DB.QueryRow("drop table ArrowInverses")
	}

	if !CreateType(ctx,NODEPTR_TYPE) {
		fmt.Println("Unable to create type as, ",NODEPTR_TYPE)
		os.Exit(-1)
	}

	if !CreateType(ctx,LINK_TYPE) {
		fmt.Println("Unable to create type as, ",LINK_TYPE)
		os.Exit(-1)
	}

	if !CreateTable(ctx,NODE_TABLE) {
		fmt.Println("Unable to create table as, ",NODE_TABLE)
		os.Exit(-1)
	}

	if !CreateTable(ctx,LINK_TABLE) {
		fmt.Println("Unable to create table as, ",LINK_TABLE)
		os.Exit(-1)
	}

	if !CreateTable(ctx,ARROW_INVERSES_TABLE) {
		fmt.Println("Unable to create table as, ",ARROW_INVERSES_TABLE)
		os.Exit(-1)
	}
	if !CreateTable(ctx,ARROW_DIRECTORY_TABLE) {
		fmt.Println("Unable to create table as, ",ARROW_DIRECTORY_TABLE)
		os.Exit(-1)
	}

	DefineStoredFunctions(ctx)

	if load_arrows {
		DownloadArrowsFromDB(ctx)
	}
}

// **************************************************************************

func Close(ctx PoSST) {
	ctx.DB.Close()
}

// **************************************************************************
// In memory representation structures
// **************************************************************************

func GetNodeTxtFromPtr(frptr NodePtr) string {

	class := frptr.Class
	index := frptr.CPtr

	var node Node

	switch class {
	case N1GRAM:
		node = NODE_DIRECTORY.N1directory[index]
	case N2GRAM:
		node = NODE_DIRECTORY.N2directory[index]
	case N3GRAM:
		node = NODE_DIRECTORY.N3directory[index]
	case LT128:
		node = NODE_DIRECTORY.LT128[index]
	case LT1024:
		node = NODE_DIRECTORY.LT1024[index]
	case GT1024:
		node = NODE_DIRECTORY.GT1024[index]
	}

	return node.S
}

// **************************************************************************

func GetNodeFromPtr(frptr NodePtr) Node {

	class := frptr.Class
	index := frptr.CPtr

	var node Node

	switch class {
	case N1GRAM:
		node = NODE_DIRECTORY.N1directory[index]
	case N2GRAM:
		node = NODE_DIRECTORY.N2directory[index]
	case N3GRAM:
		node = NODE_DIRECTORY.N3directory[index]
	case LT128:
		node = NODE_DIRECTORY.LT128[index]
	case LT1024:
		node = NODE_DIRECTORY.LT1024[index]
	case GT1024:
		node = NODE_DIRECTORY.GT1024[index]
	}

	return node
}

//**************************************************************

func AppendTextToDirectory(event Node) NodePtr {

	var cnode_slot ClassedNodePtr = -1
	var ok bool = false
	var node_alloc_ptr NodePtr

	switch event.NPtr.Class {
	case N1GRAM:
		cnode_slot,ok = NODE_DIRECTORY.N1grams[event.S]
	case N2GRAM:
		cnode_slot,ok = NODE_DIRECTORY.N2grams[event.S]
	case N3GRAM:
		cnode_slot,ok = NODE_DIRECTORY.N3grams[event.S]
	case LT128:
		cnode_slot,ok = LinearFindText(NODE_DIRECTORY.LT128,event)
	case LT1024:
		cnode_slot,ok = LinearFindText(NODE_DIRECTORY.LT1024,event)
	case GT1024:
		cnode_slot,ok = LinearFindText(NODE_DIRECTORY.GT1024,event)
	}

	node_alloc_ptr.Class = event.NPtr.Class

	if ok {
		node_alloc_ptr.CPtr = cnode_slot
		return node_alloc_ptr
	}

	switch event.NPtr.Class {
	case N1GRAM:
		cnode_slot = NODE_DIRECTORY.N1_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.N1directory = append(NODE_DIRECTORY.N1directory,event)
		NODE_DIRECTORY.N1grams[event.S] = cnode_slot
		NODE_DIRECTORY.N1_top++ 
		return node_alloc_ptr
	case N2GRAM:
		cnode_slot = NODE_DIRECTORY.N2_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.N2directory = append(NODE_DIRECTORY.N2directory,event)
		NODE_DIRECTORY.N2grams[event.S] = cnode_slot
		NODE_DIRECTORY.N2_top++
		return node_alloc_ptr
	case N3GRAM:
		cnode_slot = NODE_DIRECTORY.N3_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.N3directory = append(NODE_DIRECTORY.N3directory,event)
		NODE_DIRECTORY.N3grams[event.S] = cnode_slot
		NODE_DIRECTORY.N3_top++
		return node_alloc_ptr
	case LT128:
		cnode_slot = NODE_DIRECTORY.LT128_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.LT128 = append(NODE_DIRECTORY.LT128,event)
		NODE_DIRECTORY.LT128_top++
		return node_alloc_ptr
	case LT1024:
		cnode_slot = NODE_DIRECTORY.LT1024_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.LT1024 = append(NODE_DIRECTORY.LT1024,event)
		NODE_DIRECTORY.LT1024_top++
		return node_alloc_ptr
	case GT1024:
		cnode_slot = NODE_DIRECTORY.GT1024_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.GT1024 = append(NODE_DIRECTORY.GT1024,event)
		NODE_DIRECTORY.GT1024_top++
		return node_alloc_ptr
	}

	return NO_NODE_PTR
}

//**************************************************************

func AppendLinkToNode(frptr NodePtr,link Link,toptr NodePtr) {

	frclass := frptr.Class
	frm := frptr.CPtr
	stindex := ARROW_DIRECTORY[link.Arr].STAindex

	link.Dst = toptr // fill in the last part of the reference

	switch frclass {

	case N1GRAM:
		NODE_DIRECTORY.N1directory[frm].I[stindex] = append(NODE_DIRECTORY.N1directory[frm].I[stindex],link)
	case N2GRAM:
		NODE_DIRECTORY.N2directory[frm].I[stindex] = append(NODE_DIRECTORY.N2directory[frm].I[stindex],link)
	case N3GRAM:
		NODE_DIRECTORY.N3directory[frm].I[stindex] = append(NODE_DIRECTORY.N3directory[frm].I[stindex],link)
	case LT128:
		NODE_DIRECTORY.LT128[frm].I[stindex] = append(NODE_DIRECTORY.LT128[frm].I[stindex],link)
	case LT1024:
		NODE_DIRECTORY.LT1024[frm].I[stindex] = append(NODE_DIRECTORY.LT1024[frm].I[stindex],link)
	case GT1024:
		NODE_DIRECTORY.GT1024[frm].I[stindex] = append(NODE_DIRECTORY.GT1024[frm].I[stindex],link)
	}
}

//**************************************************************

func LinearFindText(in []Node,event Node) (ClassedNodePtr,bool) {

	for i := 0; i < len(in); i++ {

		if event.L != in[i].L {
			continue
		}

		if in[i].S == event.S {
			return ClassedNodePtr(i),true
		}
	}

	return -1,false
}

//**************************************************************

func GetSTIndexByName(stname,pm string) int {

	var encoding  int
	var sign int

	switch pm {
	case "+":
		sign = 1
	case "-":
		sign = -1
	}

	switch stname {

	case "leadsto":
		encoding = ST_ZERO + LEADSTO * sign
	case "contains":
		encoding = ST_ZERO + CONTAINS * sign
	case "properties":
		encoding = ST_ZERO + EXPRESS * sign
	case "similarity":
		encoding = ST_ZERO + NEAR
	}

	return encoding

}

//**************************************************************

func PrintSTAIndex(stindex int) string {

	sttype := stindex - ST_ZERO
	var ty string

	switch sttype {
	case -EXPRESS:
		ty = "-(expressed by)"
	case -CONTAINS:
		ty = "-(part of)"
	case -LEADSTO:
		ty = "-(arriving from)"
	case NEAR:
		ty = "(close to)"
	case LEADSTO:
		ty = "+(leading to)"
	case CONTAINS:
		ty = "+(containing)"
	case EXPRESS:
		ty = "+(expressing)"
	default:
		ty = "unknown relation!"
	}

	const green = "\x1b[36m"
	const endgreen = "\x1b[0m"

	return green + ty + endgreen
}

//**************************************************************

func InsertArrowDirectory(stname,alias,name,pm string) ArrowPtr {

	// Insert an arrow into the forward/backward indices

	var newarrow ArrowDirectory

	newarrow.STAindex = GetSTIndexByName(stname,pm)
	newarrow.Long = name
	newarrow.Short = alias
	newarrow.Ptr = ARROW_DIRECTORY_TOP

	ARROW_DIRECTORY = append(ARROW_DIRECTORY,newarrow)
	ARROW_SHORT_DIR[alias] = ARROW_DIRECTORY_TOP
	ARROW_LONG_DIR[name] = ARROW_DIRECTORY_TOP
	ARROW_DIRECTORY_TOP++

	return ARROW_DIRECTORY_TOP-1
}

//**************************************************************

func InsertInverseArrowDirectory(fwd,bwd ArrowPtr) {

	// Lookup inverse by long name, only need this in search presentation

	INVERSE_ARROWS[fwd] = bwd
	INVERSE_ARROWS[bwd] = fwd
}

//**************************************************************
// Write to database
//**************************************************************

func GraphToDB(ctx PoSST) {

	fmt.Println("Storing nodes...")

	for class := N1GRAM; class <= GT1024; class++ {
		switch class {
		case N1GRAM:
			for n := range NODE_DIRECTORY.N1directory {
				org := NODE_DIRECTORY.N1directory[n]
				UploadNodeToDB(ctx,org)
			}
		case N2GRAM:
			for n := range NODE_DIRECTORY.N2directory {
				org := NODE_DIRECTORY.N2directory[n]
				UploadNodeToDB(ctx,org)
			}
		case N3GRAM:
			for n := range NODE_DIRECTORY.N3directory {
				org := NODE_DIRECTORY.N3directory[n]
				UploadNodeToDB(ctx,org)
			}
		case LT128:
			for n := range NODE_DIRECTORY.LT128 {
				org := NODE_DIRECTORY.LT128[n]
				UploadNodeToDB(ctx,org)
			}
		case LT1024:
			for n := range NODE_DIRECTORY.LT1024 {
				org := NODE_DIRECTORY.LT1024[n]
				UploadNodeToDB(ctx,org)
			}

		case GT1024:
			for n := range NODE_DIRECTORY.GT1024 {
				org := NODE_DIRECTORY.GT1024[n]
				UploadNodeToDB(ctx,org)
			}
		}
	}


	fmt.Println("Storing Arrows...")

	for arrow := range ARROW_DIRECTORY {

		UploadArrowToDB(ctx,ArrowPtr(arrow))
	}

	fmt.Println("Storing inverse Arrows...")

	for arrow := range INVERSE_ARROWS {

		UploadInverseArrowToDB(ctx,ArrowPtr(arrow))
	}

	// CREATE INDICES

	fmt.Println("Indexing ....")

	ctx.DB.QueryRow("CREATE INDEX on NodeArrowNode (Arr,STType)")
	ctx.DB.QueryRow("CREATE INDEX on Node (((NPtr).Chan),L,S)")
}

// **************************************************************************
// Postgres
// **************************************************************************

func CreateType(ctx PoSST, defn string) bool {

	row,err := ctx.DB.Query(defn)

	if err != nil {
		s := fmt.Sprintln("Failed to create datatype PGLink ",err)
		
		if strings.Contains(s,"already exists") {
			return true
		} else {
			return false
		}
	}

	row.Close();
	return true
}

// **************************************************************************

func CreateTable(ctx PoSST,defn string) bool {

	row,err := ctx.DB.Query(defn)
	
	if err != nil {
		s := fmt.Sprintln("Failed to create a table %.10 ...",defn,err)
		
		if strings.Contains(s,"already exists") {
			return true
		} else {
			return false
		}
	}

	row.Close()
	return true
}

// **************************************************************************
// Store
// **************************************************************************

func CreateDBNode(ctx PoSST, n Node) Node {

	var qstr string

	// No need to trust the values

        n.L,n.NPtr.Class = StorageClass(n.S)
	
	cptr := n.NPtr.CPtr
	es := EscapeString(n.S)

	qstr = fmt.Sprintf("SELECT IdempInsertNode(%d,%d,%d,'%s','%s')",n.L,n.NPtr.Class,cptr,es,n.Chap)

	row,err := ctx.DB.Query(qstr)
	
	if err != nil {
		s := fmt.Sprint("Failed to insert",err)
		
		if strings.Contains(s,"duplicate key") {
		} else {
			fmt.Println(s,"FAILED \n",qstr,err)
		}
		return n
	}

	var whole string
	var cl,ch int

	for row.Next() {		
		err = row.Scan(&whole)
		fmt.Sscanf(whole,"(%d,%d)",&cl,&ch)
	}

	n.NPtr.Class = cl
	n.NPtr.CPtr = ClassedNodePtr(ch)

	row.Close()

	return n
}

// **************************************************************************

func UploadNodeToDB(ctx PoSST, org Node) {

	CreateDBNode(ctx, org)

	for stindex := range org.I {

		for lnk := range org.I[stindex] {

			dstlnk := org.I[stindex][lnk]
			sttype := STIndexToSTType(stindex)

			AppendDBLinkToNode(ctx,org.NPtr,dstlnk,sttype)
			CreateDBNodeArrowNode(ctx,org.NPtr,dstlnk,sttype)
			fmt.Print(".")
		}
	}
}

// **************************************************************************

func UploadArrowToDB(ctx PoSST,arrow ArrowPtr) {

	staidx := ARROW_DIRECTORY[arrow].STAindex
	long := ARROW_DIRECTORY[arrow].Long
	short := ARROW_DIRECTORY[arrow].Short

	qstr := fmt.Sprintf("INSERT INTO ArrowDirectory (STAindex,Long,Short,ArrPtr) VALUES (%d,'%s','%s',%d)",staidx,long,short,arrow)

	row,err := ctx.DB.Query(qstr)
	
	if err != nil {
		s := fmt.Sprint("Failed to insert",err)
		
		if strings.Contains(s,"duplicate key") {
		} else {
			fmt.Println(s,"FAILED \n",qstr,err)
		}
		return
	}

	row.Close()
}

// **************************************************************************

func UploadInverseArrowToDB(ctx PoSST,arrow ArrowPtr) {

	plus := arrow
	minus := INVERSE_ARROWS[arrow]

	qstr := fmt.Sprintf("INSERT INTO ArrowInverses (Plus,Minus) VALUES (%d,%d)",plus,minus)

	row,err := ctx.DB.Query(qstr)
	
	if err != nil {
		s := fmt.Sprint("Failed to insert",err)
		
		if strings.Contains(s,"duplicate key") {
		} else {
			fmt.Println(s,"FAILED \n",qstr,err)
		}
		return
	}

	row.Close()
}

// **************************************************************************

func AppendDBLinkToNode(ctx PoSST, n1ptr NodePtr, lnk Link, sttype int) bool {

	// Want to make this idempotent, because SQL is not (and not clause)

	if sttype < -EXPRESS || sttype > EXPRESS {
		fmt.Println(ERR_ST_OUT_OF_BOUNDS,sttype)
		os.Exit(-1)
	}

	//                       Arr,Wgt,Ctx,  Dst
	linkval := fmt.Sprintf("(%d, %f, %s, (%d,%d)::NodePtr)",lnk.Arr,lnk.Wgt,FormatSQLStringArray(lnk.Ctx),lnk.Dst.Class,lnk.Dst.CPtr)

	literal := fmt.Sprintf("%s::Link",linkval)

	link_table := STTypeDBChannel(sttype)

	qstr := fmt.Sprintf("UPDATE NODE set %s=array_append(%s,%s) where (NPtr).CPtr = '%d' and (NPtr).Chan = '%d' and (%s is null or not %s = ANY(%s))",
		link_table,
		link_table,
		literal,
		n1ptr.CPtr,
		n1ptr.Class,
		link_table,
		literal,
		link_table)

	row,err := ctx.DB.Query(qstr)

	if err != nil {
		fmt.Println("Failed to append",err,qstr)
	       return false
	}

	row.Close()
	return true
}

// **************************************************************************

func CreateDBNodeArrowNode(ctx PoSST, org NodePtr, dst Link, sttype int) bool {

	qstr := fmt.Sprintf("SELECT IdempInsertNodeArrowNode(" +
		"%d," + //infromptr
		"%d," + //infromchan
		"%d," + //isttype
		"%d," + //iarr
		"%.2f," + //iwgt
		"%s," + //ictx
		"%d," + //intoptr
		"%d " + //intochan,
		")",
		org.CPtr,
		org.Class,
		sttype,
		dst.Arr,
		dst.Wgt,
		FormatSQLStringArray(dst.Ctx),
		dst.Dst.CPtr,
		dst.Dst.Class)

	row,err := ctx.DB.Query(qstr)

	if err != nil {
		fmt.Println("Failed to make node-arrow-node",err,qstr)
	       return false
	}

	row.Close()
	return true
}

// **************************************************************************

func DefineStoredFunctions(ctx PoSST) {

	// NB! these functions are in "plpgsql" language, NOT SQL. They look similar but they are DIFFERENT!
	
	// Insert a node structure, also an anchor for and containing link arrays
	
	qstr := "CREATE OR REPLACE FUNCTION IdempInsertNode(iLi INT, iszchani INT, icptri INT, iSi TEXT, ichapi TEXT)\n" +
		"RETURNS TABLE (    \n" +
		"    ret_cptr INTEGER," +
		"    ret_channel INTEGER" +
		") AS $fn$ " +
		"DECLARE \n" +
		"BEGIN\n" +
		"  IF NOT EXISTS (SELECT (NPtr).Chan,(NPtr).CPtr FROM Node WHERE s = iSi) THEN\n" +
		"     INSERT INTO Node (Nptr.Chan,Nptr.Cptr,L,S,chap) VALUES (iszchani,icptri,iLi,iSi,ichapi);" +
		"  end if;\n" +
		"      return query SELECT (NPtr).Chan,(NPtr).CPtr FROM Node WHERE s = iSi;\n" +
		"END ;\n" +
		"$fn$ LANGUAGE plpgsql;";

	row,err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	row.Close()

	// For lookup by arrow

	qstr = "CREATE OR REPLACE FUNCTION IdempInsertNodeArrowNode\n" +
		"(\n" +
		"infromptr  int,   \n" +
		"infromchan int,   \n" +
		"isttype    int,   \n" +
		"iarr       int,   \n" +
		"iwgt       real,  \n" +
		"ictx       text[],\n" +
		"intoptr    int,   \n" +
		"intochan   int    \n" +
		")\n" +

		"RETURNS real AS $fn$ " +

		"DECLARE \n" +
		"  ret_wgt real;\n" +
		"BEGIN\n" +

		"  IF NOT EXISTS (SELECT Wgt FROM NodeArrowNode WHERE (NFrom).Cptr=infromptr AND Arr=iarr AND (NTo).Cptr=intoptr) THEN\n" +

		"     INSERT INTO NodeArrowNode (nfrom.Cptr,nfrom.Chan,sttype,arr,wgt,ctx,nto.Cptr,nto.Chan) \n" +
		"       VALUES (infromptr,infromchan,isttype,iarr,iwgt,ictx,intoptr,intochan);" +

		"  END IF;\n" +
		"  SELECT Wgt into ret_wgt FROM NodeArrowNode WHERE (NFrom).Cptr=infromptr AND Arr=iarr AND (NTo).Cptr=intoptr;\n" +
		"  RETURN ret_wgt;" +
		"END ;\n" +
		"$fn$ LANGUAGE plpgsql;";

	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	row.Close()

	// Construct an empty link pointing nowhere as a starting node

	qstr = "CREATE OR REPLACE FUNCTION GetSingletonAsLinkArray(start NodePtr)\n"+
		"RETURNS Link[] AS $fn$\n"+
		"DECLARE \n"+
		"    level Link[] := Array[] :: Link[];\n"+
		"    lnk Link := (0,1.0,Array[]::text[],(0,0));\n"+
		"BEGIN\n"+
		" lnk.Dst = start;\n"+
		" level = array_append(level,lnk);\n"+
		"RETURN level; \n"+
		"END ;\n"+
		"$fn$ LANGUAGE plpgsql;"

	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	row.Close()

	// Construct an empty link pointing nowhere as a starting node

	qstr = "CREATE OR REPLACE FUNCTION GetSingletonAsLink(start NodePtr)\n"+
		"RETURNS Link AS $fn$\n"+
		"DECLARE \n"+
		"    lnk Link := (0,1.0,Array[]::text[],(0,0));\n"+
		"BEGIN\n"+
		" lnk.Dst = start;\n"+
		"RETURN lnk; \n"+
		"END ;\n"+
		"$fn$ LANGUAGE plpgsql;"

	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	row.Close()

	// Construct search by sttype. since table names are static we need a case statement

	qstr = "CREATE OR REPLACE FUNCTION GetNeighboursByType(start NodePtr, sttype int)\n"+
		"RETURNS Link[] AS $fn$\n"+
		"DECLARE \n"+
		"    fwdlinks Link[] := Array[] :: Link[];\n"+
		"    lnk Link := (0,1.0,Array[]::text[],(0,0));\n"+
		"BEGIN\n"+
		"   CASE sttype \n"
	
	for st := -EXPRESS; st <= EXPRESS; st++ {
		qstr += fmt.Sprintf("WHEN %d THEN\n"+
			"     SELECT %s INTO fwdlinks FROM Node WHERE Nptr=start;\n",st,STTypeDBChannel(st));
	}
	qstr += "ELSE RAISE EXCEPTION 'No such sttype %', sttype;\n" +
		"END CASE;\n" +
		"    RETURN fwdlinks; \n" +
		"END ;\n" +
		"$fn$ LANGUAGE plpgsql;\n"

	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	row.Close()

	// Get the nearest neighbours as NPtr, with respect to each of the four STtype

	qstr = fmt.Sprintf("CREATE OR REPLACE FUNCTION GetFwdNodes(start NodePtr,exclude NodePtr[],sttype int)\n"+
		"RETURNS NodePtr[] AS $fn$\n" +
		"DECLARE \n" +
		"    neighbours NodePtr[];\n" +
		"    fwdlinks Link[];\n" +
		"    lnk Link;\n" +
		"BEGIN\n" +

		"    fwdlinks =GetNeighboursByType(start,sttype);\n"+

		"    IF fwdlinks IS NULL THEN\n" +
		"        RETURN '{}';\n" +
		"    END IF;\n" +

		"    neighbours := ARRAY[]::NodePtr[];\n" +

		"    FOREACH lnk IN ARRAY fwdlinks\n" +
		"    LOOP\n"+
		"      IF exclude is not NULL AND NOT lnk.dst=ANY(exclude) THEN\n" +
		"         neighbours := array_append(neighbours, lnk.dst);\n" +
		"      END IF; \n" +
		"    END LOOP;\n" +

		"    RETURN neighbours; \n" +
		"END ;\n" +
		"$fn$ LANGUAGE plpgsql;\n")

	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	row.Close()

	qstr = fmt.Sprintf("CREATE OR REPLACE FUNCTION GetFwdLinks(start NodePtr,exclude NodePtr[],sttype int)\n"+
		"RETURNS Link[] AS $fn$\n" +
		"DECLARE \n" +
		"    neighbours Link[];\n" +
		"    fwdlinks Link[];\n" +
		"    lnk Link;\n" +
		"BEGIN\n" +

		"    fwdlinks =GetNeighboursByType(start,sttype);\n"+

		"    IF fwdlinks IS NULL THEN\n" +
		"        RETURN '{}';\n" +
		"    END IF;\n" +
		"    neighbours := ARRAY[]::Link[];\n" +
		"    FOREACH lnk IN ARRAY fwdlinks\n" +
		"    LOOP\n"+
		"      IF exclude is not NULL AND NOT lnk.dst=ANY(exclude) THEN\n" +
		"         neighbours := array_append(neighbours, lnk);\n" +
		"      END IF; \n" + 
		"    END LOOP;\n" +
		"    RETURN neighbours; \n" +
		"END ;\n" +
		"$fn$ LANGUAGE plpgsql;\n")
	
	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	row.Close()
	
	// Get the forward cone / half-ball as NPtr

	qstr = "CREATE OR REPLACE FUNCTION FwdConeAsNodes(start NodePtr,sttype INT, maxdepth INT)\n"+
		"RETURNS NodePtr[] AS $fn$\n" +
		"DECLARE \n" +
		"    nextlevel NodePtr[];\n" +
		"    partlevel NodePtr[];\n" +
		"    level NodePtr[] = ARRAY[start]::NodePtr[];\n" +
		"    exclude NodePtr[] = ARRAY['(0,0)']::NodePtr[];\n" +
		"    cone NodePtr[];\n" +
		"    neigh NodePtr;\n" +
		"    frn NodePtr;\n" +
		"    counter int := 0;\n" +

		"BEGIN\n" +

		"LOOP\n" +
		"  EXIT WHEN counter = maxdepth+1;\n" +

		"  IF level IS NULL THEN\n" +
		"     RETURN cone;\n" +
		"  END IF;\n" +

		"  nextlevel := ARRAY[]::NodePtr[];\n" +

		"  FOREACH frn IN ARRAY level "+
		"  LOOP \n"+
		"     nextlevel = array_append(nextlevel,frn);\n" +
		"  END LOOP;\n" +

		"  IF nextlevel IS NULL THEN\n" +
		"     RETURN cone;\n" +
		"  END IF;\n" +

		"  FOREACH neigh IN ARRAY nextlevel LOOP \n"+
		"    IF NOT neigh = ANY(exclude) THEN\n" +
		"      cone = array_append(cone,neigh);\n" +
		"      exclude := array_append(exclude,neigh);\n" +
		"      partlevel := GetFwdNodes(neigh,exclude,sttype);\n" +
		"    END IF;" +
		"    IF partlevel IS NOT NULL THEN\n" +
		"         level = array_cat(level,partlevel);\n"+
		"    END IF;\n" +
		"  END LOOP;\n" +

		// Next, continue, foreach
		"  counter = counter + 1;\n" +
		"END LOOP;\n" +
		
		"RETURN cone; \n" +
		"END ;\n" +
		"$fn$ LANGUAGE plpgsql;\n"
	
	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}
	
	row.Close()
	
          /* e.g. select unnest(fwdconeaslinks) from FwdConeAsLinks('(4,1)',1,4);
                           unnest                           
             ------------------------------------------------------------
              (0,0,{},"(4,1)")
              (77,0.34,"{ ""fairy castles"", ""angel air"" }","(4,2)")
              (77,0.34,"{ ""fairy castles"", ""angel air"" }","(4,3)")
              (77,0.34,"{ ""steamy hot tubs"" }","(4,5)")
              (77,0.34,"{ ""fairy castles"", ""angel air"" }","(4,4)")
              (77,0.34,"{ ""steamy hot tubs"", ""lady gaga"" }","(4,6)")
             (6 rows)

          */

	qstr = "CREATE OR REPLACE FUNCTION FwdConeAsLinks(start NodePtr,sttype INT,maxdepth INT)\n"+
		"RETURNS Link[] AS $fn$\n" +
		"DECLARE \n" +
		"    nextlevel Link[];\n" +
		"    partlevel Link[];\n" +
		"    level Link[] = ARRAY[]::Link[];\n" +
		"    exclude NodePtr[] = ARRAY['(0,0)']::NodePtr[];\n" +
		"    cone Link[];\n" +
		"    neigh Link;\n" +
		"    frn Link;\n" +
		"    counter int := 0;\n" +

		"BEGIN\n" +

		"level := GetSingletonAsLinkArray(start);\n"+

		"LOOP\n" +
		"  EXIT WHEN counter = maxdepth+1;\n" +

		"  IF level IS NULL THEN\n" +
		"     RETURN cone;\n" +
		"  END IF;\n" +

		"  nextlevel := ARRAY[]::Link[];\n" +

		"  FOREACH frn IN ARRAY level "+
		"  LOOP \n"+
		"     nextlevel = array_append(nextlevel,frn);\n" +
		"  END LOOP;\n" +

		"  IF nextlevel IS NULL THEN\n" +
		"     RETURN cone;\n" +
		"  END IF;\n" +

		"  FOREACH neigh IN ARRAY nextlevel LOOP \n"+
		"    IF NOT neigh.Dst = ANY(exclude) THEN\n" +
		"      cone = array_append(cone,neigh);\n" +
		"      exclude := array_append(exclude,neigh.Dst);\n" +
		"      partlevel := GetFwdLinks(neigh.Dst,exclude,sttype);\n" +
		"    END IF;" +
		"    IF partlevel IS NOT NULL THEN\n" +
		"         level = array_cat(level,partlevel);\n"+
		"    END IF;\n" +
		"  END LOOP;\n" +

		"  counter = counter + 1;\n" +
		"END LOOP;\n" +

		"RETURN cone; \n" +
		"END ;\n" +
		"$fn$ LANGUAGE plpgsql;\n"

	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	row.Close()

	// Orthogonal (depth first) paths from origin spreading out

	qstr = "CREATE OR REPLACE FUNCTION FwdPathsAsLinks(start NodePtr,sttype INT,maxdepth INT)\n"+
		"RETURNS Text AS $fn$\n" +
		"DECLARE\n" +
		"   hop Text;\n" +
		"   path Text;\n"+
		"   summary_path Text[];\n"+
		"   exclude NodePtr[] = ARRAY['(0,0)']::NodePtr[];\n" +
		"   ret_paths Text;\n" +
		"   startlnk Link;"+

		"BEGIN\n" +

		"startlnk := GetSingletonAsLink(start);\n"+
		"path := Format('%s',startlnk::Text);\n"+
		"ret_paths := SumFwdPaths(startlnk,path,sttype,1,maxdepth,exclude);" +

		"RETURN ret_paths; \n" +
		"END ;\n" +
		"$fn$ LANGUAGE plpgsql;\n"

        // select FwdPathsAsLinks('(4,1)',1,3)

	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	row.Close()

	// Return end of path branches as aggregated text summaries

	qstr = "CREATE OR REPLACE FUNCTION SumFwdPaths(start Link,path TEXT, sttype INT,depth int, maxdepth INT,exclude NodePtr[])\n"+
		"RETURNS Text AS $fn$\n" +
		"DECLARE \n" + 
		"    fwdlinks Link[];\n" +
		"    empty Link[] = ARRAY[]::Link[];\n" +
		"    lnk Link;\n" +
		"    fwd Link;\n" +
		"    ret_paths Text;\n" +
		"    appendix Text;\n" +
		"    tot_path Text;\n"+
		"BEGIN\n" +

		"IF depth = maxdepth THEN\n"+
//		"  RAISE NOTICE 'Xend path %',path;"+
		"  ret_paths := Format('%s\n%s',ret_paths,path);\n"+
		"  RETURN ret_paths;\n"+
		"END IF;\n"+

		"fwdlinks := GetFwdLinks(start.Dst,exclude,sttype);\n" +

		"FOREACH lnk IN ARRAY fwdlinks LOOP \n" +
		"   IF NOT lnk.Dst = ANY(exclude) THEN\n"+
		"      exclude = array_append(exclude,lnk.Dst);\n" +
		"      IF lnk IS NULL THEN" +
		"         ret_paths := Format('%s\n%s',ret_paths,path);\n"+
//		"         RAISE NOTICE 'Yend path %',tot_path;"+
		"      ELSE"+
		"         tot_path := Format('%s;%s',path,lnk::Text);\n"+
		"         appendix := SumFwdPaths(lnk,tot_path,sttype,depth+1,maxdepth,exclude);\n" +
		"         IF appendix IS NOT NULL THEN\n"+
		"            ret_paths := Format('%s\n%s',ret_paths,appendix);\n"+
		"         END IF;"+
		"      END IF;"+
		"   END IF;"+
		"END LOOP;"+

		"RETURN ret_paths; \n" +
		"END ;\n" +
		"$fn$ LANGUAGE plpgsql;\n"

	// select SumFwdPaths('(4,1)',1,1,3);

	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	row.Close()
}

// **************************************************************************
// Retrieve
// **************************************************************************

func GetDBNodePtrMatchingName(ctx PoSST,s string) []NodePtr {

	search := "%"+s+"%"

	qstr := fmt.Sprintf("select NPtr from Node where S LIKE '%s'",search)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY GetNodePtrMatchingName Failed",err)
	}

	var whole string
	var n NodePtr
	var retval []NodePtr

	for row.Next() {		
		err = row.Scan(&whole)
		fmt.Sscanf(whole,"(%d,%d)",&n.Class,&n.CPtr)
		retval = append(retval,n)
	}

	row.Close()
	return retval

}

// **************************************************************************

func GetDBNodePtrMatchingNCC(ctx PoSST,chap,nm string ,cn []string) []NodePtr {

	// Match name, context, chapter

	name := "%"+nm+"%"
	context := FormatSQLStringArray(cn)
	chapter := "%"+chap+"%"

	qstr := fmt.Sprintf("WITH matching_nodes AS "+
		"  (SELECT NFrom,ctx,match_context(ctx,%s) AS match FROM NodeArrowNode)"+
		"     SELECT DISTINCT nfrom FROM matching_nodes "+
		"      JOIN Node ON nptr=nfrom WHERE match=true AND S LIKE '%s' AND chap LIKE '%s'",context,name,chapter)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY GetNodePtrMatchingNCC Failed",err)
	}

	var whole string
	var n NodePtr
	var retval []NodePtr

	for row.Next() {		
		err = row.Scan(&whole)
		fmt.Sscanf(whole,"(%d,%d)",&n.Class,&n.CPtr)
		retval = append(retval,n)
	}

	row.Close()
	return retval

}

// **************************************************************************

func GetDBNodeByNodePtr(ctx PoSST,db_nptr NodePtr) Node {

	im_nptr,cached := NODE_CACHE[db_nptr]

	if cached {
		return GetNodeFromPtr(im_nptr)
	}

	qstr := fmt.Sprintf("select L,S,Chap from Node where NPtr='(%d,%d)'::NodePtr",db_nptr.Class,db_nptr.CPtr)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("GetNodeByNodePointer Failed:",err)
	}

	var n Node
	var count int = 0

	n.NPtr = db_nptr

	for row.Next() {		
		err = row.Scan(&n.L,&n.S,&n.Chap)
		count++
	}

	if count > 1 {
		fmt.Println("GetNodeByNodePtr returned too many matches (multi-model conflict?):",count,"for ptr",db_nptr)
		os.Exit(-1)
	}

	row.Close()

	if !cached {
		CacheNode(n)
	}

	return n
}

// **************************************************************************

func GetDBArrowsMatchingArrowName(ctx PoSST,s string) []ArrowPtr {

	var list []ArrowPtr

	if ARROW_DIRECTORY_TOP == 0 {
		DownloadArrowsFromDB(ctx)
	}

	for a := range ARROW_DIRECTORY {

		if SimilarString(s,ARROW_DIRECTORY[a].Long) || SimilarString(s,ARROW_DIRECTORY[a].Short) {
			list = append(list,ARROW_DIRECTORY[a].Ptr)
		}
	}

	return list
}

// **************************************************************************

func GetDBNodeArrowNodeMatchingArrowPtrs(ctx PoSST,arrows []ArrowPtr) []NodeArrowNode {

	var intarrows []int

	for i := range arrows {
		intarrows = append(intarrows,int(arrows[i]))
	}

	qstr := fmt.Sprintf("SELECT NFrom,STType,Arr,Wgt,Ctx,NTo FROM NodeArrowNode where Arr=ANY(%s::int[])",FormatSQLIntArray(intarrows))

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("GetDBNodeArrowNodeMatchingArrowPtrs Failed:",err,qstr)
	}

	var from_node string
	var to_node string
	var actx string
	var st,arr int
	var wgt float64

	var nfr,nto NodePtr
	var nan NodeArrowNode
	var nanlist []NodeArrowNode

	for row.Next() {		
		err = row.Scan(&from_node,&st,&arr,&wgt,&actx,&to_node)

		fmt.Sscanf(from_node,"(%d,%d)",&nfr.Class,&nfr.CPtr)
		fmt.Sscanf(to_node,"(%d,%d)",&nto.Class,&nto.CPtr)

		nan.NFrom = nfr
		nan.STType = st
		nan.Arr = ArrowPtr(arr)
		nan.Wgt = wgt
		nan.Ctx = ParseSQLArrayString(actx)
		nan.NTo = nto

		nanlist = append(nanlist,nan)

	}

	row.Close()

	return nanlist
}

// **************************************************************************

func GetDBArrowByName(ctx PoSST,name string) ArrowPtr {

	ptr, ok := ARROW_SHORT_DIR[name]
	
	// If not, then check longname
	
	if !ok {
		ptr, ok = ARROW_LONG_DIR[name]
		
		if !ok {
			DownloadArrowsFromDB(ctx)

			ptr, ok = ARROW_SHORT_DIR[name]
			
			// If not, then check longname
			
			if !ok {
				ptr, ok = ARROW_LONG_DIR[name]
				fmt.Println(ERR_NO_SUCH_ARROW,name)
				os.Exit(-1)
			}
		}
	}
	return ptr
}

// **************************************************************************

func GetDBArrowByPtr(ctx PoSST,arrowptr ArrowPtr) ArrowDirectory {

	if ARROW_DIRECTORY_TOP > 0 {
		a := ARROW_DIRECTORY[arrowptr]
		return a
	}

	DownloadArrowsFromDB(ctx)

	if len(ARROW_DIRECTORY) < int(arrowptr) {
		fmt.Println(ERR_NO_SUCH_ARROW,arrowptr)
		os.Exit(-1)
	}

	return ARROW_DIRECTORY[arrowptr]

}

// **************************************************************************

func CacheNode(n Node) {

	NODE_CACHE[n.NPtr] = AppendTextToDirectory(n)
}

// **************************************************************************

func DownloadArrowsFromDB(ctx PoSST) {

	// These must be ordered to match in-memory array

	qstr := fmt.Sprintf("SELECT STAindex,Long,Short,ArrPtr FROM ArrowDirectory ORDER BY ArrPtr")

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY Download Arrows Failed",err)
	}

	var staidx int
	var long string
	var short string
	var ptr ArrowPtr
	var ad ArrowDirectory

	for row.Next() {		
		err = row.Scan(&staidx,&long,&short,&ptr)
		ad.STAindex = staidx
		ad.Long = long
		ad.Short = short
		ad.Ptr = ptr

		ARROW_DIRECTORY = append(ARROW_DIRECTORY,ad)
		ARROW_SHORT_DIR[short] = ARROW_DIRECTORY_TOP
		ARROW_LONG_DIR[long] = ARROW_DIRECTORY_TOP

		if ad.Ptr != ARROW_DIRECTORY_TOP {
			fmt.Println(ERR_MEMORY_DB_ARROW_MISMATCH,ad,ad.Ptr,ARROW_DIRECTORY_TOP)
			os.Exit(-1)
		}

		ARROW_DIRECTORY_TOP++
	}

	row.Close()

	// Get Inverses

	qstr = fmt.Sprintf("SELECT Plus,Minus FROM ArrowInverses ORDER BY Plus")

	row, err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY Download Inverses Failed",err)
	}

	var plus,minus ArrowPtr

	for row.Next() {		

		err = row.Scan(&plus,&minus)

		if err != nil {
			fmt.Println("QUERY Download Arrows Failed",err)
		}

		INVERSE_ARROWS[plus] = minus
	}
}

// **************************************************************************

func GetFwdConeAsNodes(ctx PoSST, start NodePtr, sttype,depth int) []NodePtr {

	qstr := fmt.Sprintf("select unnest(fwdconeasnodes) from FwdConeAsNodes('(%d,%d)',%d,%d);",start.Class,start.CPtr,sttype,depth)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY to FwdConeAsNodes Failed",err)
	}

	var whole string
	var n NodePtr
	var retval []NodePtr

	for row.Next() {		
		err = row.Scan(&whole)
		fmt.Sscanf(whole,"(%d,%d)",&n.Class,&n.CPtr)
		retval = append(retval,n)
	}

	row.Close()
	return retval
}

// **************************************************************************

func GetFwdConeAsLinks(ctx PoSST, start NodePtr, sttype,depth int) []Link {

	qstr := fmt.Sprintf("select unnest(fwdconeaslinks) from FwdConeAsLinks('(%d,%d)',%d,%d);",start.Class,start.CPtr,sttype,depth)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY to FwdConeAsLinkss Failed",err)
	}

	var whole string
	var retval []Link

	for row.Next() {		
		err = row.Scan(&whole)
		l := ParseSQLLinkString(whole)
		retval = append(retval,l)
	}

	row.Close()

	return retval
}

// **************************************************************************

func GetFwdPathsAsLinks(ctx PoSST, start NodePtr, sttype,depth int) ([][]Link,int) {

	qstr := fmt.Sprintf("select FwdPathsAsLinks from FwdPathsAsLinks('(%d,%d)',%d,%d);",start.Class,start.CPtr,sttype,depth)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY to FwdPathsAsLinks Failed",err)
	}

	var whole string
	var retval [][]Link

	for row.Next() {		
		err = row.Scan(&whole)
		retval = ParseLinkPath(whole)
	}

	row.Close()
	return retval,len(retval)
}

// **************************************************************************

func PrintLinkPath(ctx PoSST, alt_paths [][]Link, p int, prefix string) {

	if len(alt_paths[p]) > 1 {

		path_start := GetDBNodeByNodePtr(ctx,alt_paths[p][0].Dst)		
		
		fmt.Print(prefix," from ",path_start.S)
		
		for l := 1; l < len(alt_paths[p]); l++ {
			
			arr := GetDBArrowByPtr(ctx,alt_paths[p][l].Arr)
			
			if l < len(alt_paths[p]) {
				fmt.Print("  -(",arr.Long,")->  ")
			}
			
			nextnode := GetDBNodeByNodePtr(ctx,alt_paths[p][l].Dst)
			fmt.Print(nextnode.S)
		}
		
		fmt.Println()
	}
}

// **************************************************************************
// Retrieve Analysis
// **************************************************************************

func GetMatroidArrayByArrow(ctx PoSST, context,chapter string) map[ArrowPtr][]NodePtr {

          /* arr |             x             
            -----+---------------------------
              18 | {"(2,4)","(3,4)","(4,4)"}
             138 | {"(4,4)","(0,4)"}
              97 | {"(1,2)"}
              96 | {"(0,4)"}
             137 | {"(1,4)","(0,3)"}
              52 | {"(0,4)"}
              53 | {"(0,2)"} */

	var qplus string

	// Postgres && operator on arrays is SET OVERLAP .. how to solve this?

	if context != "any" {
		qplus += fmt.Sprintf("WHERE %s LIKE ANY(CTX)",context)
	}

	qstr := "SELECT arr,array_agg(DISTINCT NTo) FROM NodeArrowNode GROUP BY arr"

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY GetMatroidArrayByArrow Failed",err)
	}

	var arry string
	var arr ArrowPtr
	var retval = make(map[ArrowPtr][]NodePtr)

	for row.Next() {		
		err = row.Scan(&arr,&arry)
		retval[arr] = ParseSQLNPtrArray(arry)
	}

	row.Close()
	
	return retval
}

// **************************************************************************

func GetMatroidArrayBySSType(ctx PoSST) map[int][]NodePtr {


          /* sttype |             array_agg             
             --------+-----------------------------------
                 -3 | {"(0,4)","(4,4)"}
                 -2 | {"(1,2)"}
                 -1 | {"(0,2)"}
                  1 | {"(0,4)","(2,4)","(3,4)","(4,4)"}
                  2 | {"(0,4)"}
                  3 | {"(0,3)","(1,4)"} */


	qstr := "SELECT sttype, array_agg(DISTINCT NTo) FROM NodeArrowNode GROUP BY Sttype order by sttype"

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY GetMatroidArrayByArrow Failed",err)
	}

	var arry string
	var sttype int
	var retval = make(map[int][]NodePtr)

	for row.Next() {		
		err = row.Scan(&sttype,&arry)
		retval[sttype] = ParseSQLNPtrArray(arry)
	}

	row.Close()
	
	return retval
}

// **************************************************************************

func GetMatroidHistogramByArrow(ctx PoSST) map[ArrowPtr]int {

/* arr | count 
-----+-------
  18 |     3
  52 |     1
  53 |     1
  96 |     1
  97 |     1
 137 |     2
 138 |     2 */

	qstr := "SELECT arr,count(NTo) FROM NodeArrowNode GROUP BY arr order by Arr"

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY GetMatroidArrayByArrow Failed",err)
	}

	var freq int
	var arr ArrowPtr
	var retval = make(map[ArrowPtr]int)

	for row.Next() {		
		err = row.Scan(&arr,&freq)
		retval[arr] = freq
	}

	row.Close()
	
	return retval
}

// **************************************************************************

func GetMatroidHistogramBySSType(ctx PoSST) map[int]int {

/* sttype | x 
--------+---
      1 | 4
     -1 | 1
      2 | 1
     -2 | 1
      3 | 2
     -3 | 2 */

	qstr := "SELECT sttype,count(NTo) FROM NodeArrowNode GROUP BY Sttype"

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY GetMatroidArrayByArrow Failed",err)
	}

	var freq int
	var sttype int
	var retval = make(map[int]int)

	for row.Next() {		
		err = row.Scan(&sttype,&freq)
		retval[sttype] = freq
	}
	return retval
}

// **************************************************************************

func GetMatroidNodesByArrow(ctx PoSST) []ArrowMatroid {

/*  nfrom | arr | array_agg 
-------+-----+-----------
 (4,0) |  18 | {"(2,4)"}
 (4,2) |  18 | {"(3,4)"}
 (4,3) |  18 | {"(4,4)"}
 (2,0) |  52 | {"(0,4)"}
 (4,0) |  53 | {"(0,2)"}
 (2,1) |  96 | {"(0,4)"}
 (4,0) |  97 | {"(1,2)"}
 (4,0) | 137 | {"(1,4)"}
 (4,4) | 137 | {"(0,3)"}
 (3,0) | 138 | {"(4,4)"}
 (4,1) | 138 | {"(0,4)"} */

	qstr := "SELECT NFrom,Arr,array_agg(NTo) FROM NodeArrowNode GROUP BY Arr,Nfrom HAVING count(NTo) > 1 ORDER BY Arr "

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY GetMatroidArrayByArrow Failed",err)
	}

	var nptr,arry string
	var retval []ArrowMatroid
	var this ArrowMatroid

	for row.Next() {		
		err = row.Scan(&nptr,&this.Arr,&arry)
		fmt.Sscanf(nptr,"(%d,%d)",&this.NFrom.Class,&this.NFrom.CPtr)
		this.NTo = ParseSQLNPtrArray(arry)
		retval = append(retval,this)
	}

	row.Close()

	return retval
}

// **************************************************************************

func GetMatroidNodesBySTType(ctx PoSST) []STTypeMatroid {

/*  nfrom | sttype | array_agg 
-------+--------+-----------
 (3,0) |     -3 | {"(4,4)"}
 (4,1) |     -3 | {"(0,4)"}
 (4,0) |     -2 | {"(1,2)"}
 (4,0) |     -1 | {"(0,2)"}
 (2,0) |      1 | {"(0,4)"}
 (4,0) |      1 | {"(2,4)"}
 (4,2) |      1 | {"(3,4)"}
 (4,3) |      1 | {"(4,4)"}
 (2,1) |      2 | {"(0,4)"}
 (4,0) |      3 | {"(1,4)"}
 (4,4) |      3 | {"(0,3)"}*/

	qstr := "SELECT NFrom,sttype,array_agg(NTo) FROM NodeArrowNode GROUP BY sttype,Nfrom HAVING count(NTo) > 1 ORDER BY sttype"

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY GetMatroidArrayBySTType Failed",err)
	}

	var retval []STTypeMatroid
	var this STTypeMatroid
	var nptr,arr,arry string

	for row.Next() {		
		err = row.Scan(&nptr,&arr,&arry)
		err = row.Scan(&nptr,&this.STType,&arry)
		fmt.Sscanf(nptr,"(%d,%d)",&this.NFrom.Class,&this.NFrom.CPtr)
		this.NTo = ParseSQLNPtrArray(arry)
		retval = append(retval,this)
	}

	row.Close()
	return retval
}

// **************************************************************************
// Tools
// **************************************************************************

func ParseSQLNPtrArray(s string) []NodePtr {

	stringify := ParseSQLArrayString(s)

	var retval []NodePtr
	var nptr NodePtr

	for n := 0; n < len(stringify); n++ {
		fmt.Sscanf(stringify[n],"(%d,%d)",&nptr.Class,&nptr.CPtr)
		retval = append(retval,nptr)
	}

	return retval
}

// **************************************************************************

func ParseSQLArrayString(whole_array string) []string {

	// array as {"(1,2,3)","(4,5,6)"}

      	var l []string

    	whole_array = strings.Replace(whole_array,"{","",-1)
    	whole_array = strings.Replace(whole_array,"}","",-1)
	whole_array = strings.Replace(whole_array,"\",\"",";",-1)
	whole_array = strings.Replace(whole_array,"\"","",-1)
	
        items := strings.Split(whole_array,";")

	for i := range items {

	    s := strings.TrimSpace(items[i])

	    l = append(l,s)
	    }

	return l
}

// **************************************************************************

func FormatSQLIntArray(array []int) string {

        if len(array) == 0 {
		return "'{ }'"
        }

	var ret string = "'{ "
	
	for i := 0; i < len(array); i++ {
		ret += fmt.Sprintf("%d",array[i])
	    if i < len(array)-1 {
	    ret += ", "
	    }
        }

	ret += " }' "

	return ret
}

// **************************************************************************

func FormatSQLStringArray(array []string) string {

        if len(array) == 0 {
		return "'{ }'"
        }

	var ret string = "'{ "
	
	for i := 0; i < len(array); i++ {
		ret += fmt.Sprintf("\"%s\"",array[i])
	    if i < len(array)-1 {
	    ret += ", "
	    }
        }

	ret += " }' "

	return ret
}

// **************************************************************************

func ParseSQLLinkString(s string) Link {

        // e.g. (77,0.34,"{ ""fairy castles"", ""angel air"" }","(4,2)")
	// This feels dangerous. Is postgres consistent here?

      	var l Link

    	s = strings.Replace(s,"(","",-1)
    	s = strings.Replace(s,")","",-1)
	s = strings.Replace(s,"\"\"",";",-1)
	s = strings.Replace(s,"\"","",-1)
	
        items := strings.Split(s,",")

	// Arrow type
	fmt.Sscanf(items[0],"%d",&l.Arr)

	// Link weight
	fmt.Sscanf(items[1],"%f",&l.Wgt)

	// These are the context array

	var array []string

	for i := 2; i <= len(items)-3; i++ {
		items[i] = strings.Replace(items[i],"{","",-1)
		items[i] = strings.Replace(items[i],"}","",-1)
		items[i] = strings.Replace(items[i],";","",-1)
		items[i] = strings.TrimSpace(items[i])
		array = append(array,items[i])
	}

	l.Ctx = array

	// the last two are the NPtr

	fmt.Sscanf(items[len(items)-2],"%d",&l.Dst.Class)
	fmt.Sscanf(items[len(items)-1],"%d",&l.Dst.CPtr)

	return l
}

//**************************************************************

func ParseLinkPath(s string) [][]Link {

	// Each path will start on a new line, with comma sep Link encodings

	var array [][]Link
	var index int = 0

	lines := strings.Split(s,"\n")

	for line := range lines {

		if len(lines[line]) > 0 {

			links := strings.Split(lines[line],";")

			if len(links) < 2 {
				continue
			}

			array = append(array,make([]Link,0))

			for l := 0; l < len(links); l++ {

				lnk := ParseSQLLinkString(links[l])
				array[index] = append(array[index],lnk)
			}
			index++
		}
	}

	if index < 1 {
		return nil
	}
	return array
}

//**************************************************************

func StorageClass(s string) (int,int) {
	
	var spaces int = 0

	var l = len(s)
	
	for i := 0; i < l; i++ {
		
		if s[i] == ' ' {
			spaces++
		}
		
		if spaces > 2 {
			break
		}
	}
	
	// Text usage tends to fall into a number of different roles, with a power law
	// frequency of occurrence in a text, so let's classify in order of likely usage
	// for small and many, we use a hashmap/btree
	
	switch spaces {
	case 0:
		return l,N1GRAM
	case 1:
		return l,N2GRAM
	case 2:
		return l,N3GRAM
	}
	
	// For longer strings, a linear search is probably fine here
        // (once it gets into a database, it's someone else's problem)
	
	if l < 128 {
		return l,LT128
	}
	
	if l < 1024 {
		return l,LT1024
	}
	
	return l,GT1024
}

// **************************************************************************

func STTypeDBChannel(sttype int) string {

	// This expects the range for sttype to be unshifted 0,+/-

	var link_channel string
	switch sttype {

	case NEAR:
		link_channel = I_NEAR
	case LEADSTO:
		link_channel = I_PLEAD
	case CONTAINS:
		link_channel = I_PCONT
	case EXPRESS:
		link_channel = I_PEXPR
	case -LEADSTO:
		link_channel = I_MLEAD
	case -CONTAINS:
		link_channel = I_MCONT
	case -EXPRESS:
		link_channel = I_MEXPR
	default:
		fmt.Println(ERR_ILLEGAL_LINK_CLASS,sttype)
		os.Exit(-1)
	}

	return link_channel
}

// **************************************************************************

func STIndexToSTType(stindex int) int {

	// Convert shifted array index to symmetrical type

	return stindex - ST_ZERO
}

// **************************************************************************

func STTypeName(sttype int) string {

	switch sttype {
	case -EXPRESS:
		return "-Properties"
	case -CONTAINS:
		return "-Contains"
	case -LEADSTO:
		return "-LeadsTo"
	case NEAR:
		return "Similarity"
	case LEADSTO:
		return "+LeadsTo"
	case CONTAINS:
		return "+Contains"
	case EXPRESS:
		return "+Properties"
	}

	return "Unknown ST type"
}

// **************************************************************************

func Already (s string, cone map[int][]string) bool {

	for l := range cone {
		for n := 0; n < len(cone[l]); n++ {
			if s == cone[l][n] {
				return true
			}
		}
	}

	return false
}

// **************************************************************************

func EscapeString(s string) string {

	return strings.Replace(s, "'", ".", -1)
}

// **************************************************************************

func SimilarString(s1,s2 string) bool {

	// Placeholder

	// Need to handle pluralisation patterns etc... multi-language

	if strings.Contains(s2,s1) || strings.Contains(s1,s2) {
		return true
	}

	return false
}

//****************************************************************************

func NewLine(n int) {

	if n % 8 == 0 {
		fmt.Println()
	}
}
