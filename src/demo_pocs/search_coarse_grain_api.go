//******************************************************************
//
// Find <end|start> transition matrix and calculate symmetries
//
//
//******************************************************************

package main

import (
	"fmt"

        SST "SSTorytime"
)

//******************************************************************

const (
	host     = "localhost"
	port     = 5432
	user     = "sstoryline"
	password = "sst_1234"
	dbname   = "sstoryline"
)

//******************************************************************

func main() {

	load_arrows := true
	ctx := SST.Open(load_arrows)

	// Contra colliding wavefronts as path integral solver

	const maxdepth = 4

	context := []string{""}
	chapter := "slit"

	start_bc := []string{"start"}
	end_bc := []string{"target_1","target_2","target_3"}

/*	context := []string{""}
	chapter := "slit"

	start_bc := []string{"A1"}
	end_bc := []string{"B6"}*/

	var start_set,end_set []SST.NodePtr

	for n := range start_bc {
		start_set = append(start_set,SST.GetDBNodePtrMatchingName(ctx,"",start_bc[n])...)
	}

	for n := range end_bc {
		end_set = append(end_set,SST.GetDBNodePtrMatchingName(ctx,"",end_bc[n])...)
	}

	solutions := SST.GetPathsAndSymmetries(ctx,start_set,end_set,chapter,context,maxdepth)

	var count int

	// ***** paths ****

	fmt.Println("-- T R E E ----------------------------------")
	fmt.Println("Path solution",count,"from",start_bc,"to",end_bc)
	
	for s := 0; s < len(solutions); s++ {
		prefix := fmt.Sprintf(" - story %d: ",s)
		SST.PrintLinkPath(ctx,solutions,s,prefix,"",nil)
	}
	count++
	fmt.Println("-------------------------------------------")

	// **** Process symmetries ***

	supernodes := SST.GetPathTransverseSuperNodes(ctx,solutions,maxdepth)

	fmt.Println("Look for coarse grains, final matroid:",len(supernodes))

	for g := range supernodes {
		fmt.Print("\n    - Super node ",g," = {")
		for n := range supernodes[g] {
			node :=SST.GetDBNodeByNodePtr(ctx,supernodes[g][n])
			fmt.Print(node.S,",")
		}
		fmt.Println("}")
	}

}

//******************************************************************

func ShowNode(ctx SST.PoSST,nptr []SST.NodePtr) string {

	var ret string

	for n := range nptr {
		node := SST.GetDBNodeByNodePtr(ctx,nptr[n])
		ret += node.S + ","
	}

	return ret
}

// **********************************************************

func ShowNodePath(ctx SST.PoSST,lnk []SST.Link) string {

	var ret string

	for n := range lnk {
		node := SST.GetDBNodeByNodePtr(ctx,lnk[n].Dst)
		arrs := SST.GetDBArrowByPtr(ctx,lnk[n].Arr).Long
		ret += fmt.Sprintf("(%s) -> %s ",arrs,node.S)
	}

	return ret
}









