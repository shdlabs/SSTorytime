//******************************************************************
//
// Study graph properties
// 
//
//******************************************************************

package main

import (
	"fmt"

        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := true
	ctx := SST.Open(load_arrows)

	sttypes := []int{1}

	adj,lookup := GetDBAdjacentNodePtrBySTType(ctx,sttypes)

	fmt.Println(adj,lookup)

	SST.Close(ctx)
}

// **************************************************************************

func GetDBAdjacentNodePtrBySTType(ctx SST.PoSST,sttypes []int) ([][]float32,[]SST.NodePtr) {

	// Return a connected adjacency matrix for the subgraph and a lookup table
	// A bit memory intensive, but possibly unavoidable

	var qstr,qwhere,qsearch string
	var dim = len(sttypes)

	if dim > 4 {
		fmt.Println("Maximum 4 sttypes in GetDBAdjacentNodePtrBySTType")
		return nil,nil
	}

	for st := 0; st < len(sttypes); st++ {

		stname := SST.STTypeDBChannel(sttypes[st])
		qwhere += fmt.Sprintf("array_length(%s::text[],1) IS NOT NULL",stname)

		if st != dim-1 {
			qwhere += " AND "
		}

		qsearch += "," + stname

	}

	qstr = fmt.Sprintf("SELECT NPtr%s FROM Node WHERE %s",qsearch,qwhere)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY GetDBAdjacentNodePtrBySTType Failed",err)
		return nil,nil
	}

	var linkstr = make ([]string,dim)
	var protoadj = make(map[int][]SST.Link)
	var lookup = make(map[SST.NodePtr]int)
	var rowindex int
	var nodekey []SST.NodePtr
	var counter int

	for row.Next() {		

		var n SST.NodePtr
		var nstr string

		switch dim {

		case 1: err = row.Scan(&nstr,&linkstr[0])
		case 2: err = row.Scan(&nstr,&linkstr[0],&linkstr[1])
		case 3: err = row.Scan(&nstr,&linkstr[0],&linkstr[1],&linkstr[2])
		case 4: err = row.Scan(&nstr,&linkstr[0],&linkstr[1],&linkstr[2],&linkstr[3])

		default:
			fmt.Println("Maximum 4 sttypes in GetDBAdjacentNodePtrBySTType - shouldn't happen")
			row.Close()
			return nil,nil
		}

		if err != nil {
			fmt.Println("Error scanning sql data",err)
			row.Close()
			return nil,nil
		}

		fmt.Sscanf(nstr,"(%d,%d)",&n.Class,&n.CPtr)

		// idempotently gather nptrs into a map, keeping linked nodes close in order

		index,already := lookup[n]

		if already {
			rowindex = index
		} else {
			rowindex = counter
			lookup[n] = counter
			counter++
			nodekey = append(nodekey,n)
		}

		for lnks := range linkstr {

			links := SST.ParseLinkArray(linkstr[lnks])
			
			// we have to go through one by one to avoid duplicates
			// and keep adjacent nodes closer in order
			
			for l := range links {
				
				_,already := lookup[links[l].Dst]
				
				if !already {
					lookup[links[l].Dst] = counter
					counter++
					nodekey = append(nodekey,links[l].Dst)
				}
			}
			protoadj[rowindex] = links // now have a sparse ordered repr		
		}
	}
	
	adj := make([][]float32,counter)

	for r := 0; r < counter; r++ {

		adj[r] = make([]float32,counter)

		row := protoadj[r]

		for l := 0; l < len(row); l++ {
			lnk := row[l]
			c := lookup[lnk.Dst]
			adj[r][c] = lnk.Wgt
		}
	}

	row.Close()

	return adj,nodekey

}









