//******************************************************************
//
// Try out neighbour search for all ST stypes together
//
// Prepare:
// cd examples
// ../src/N4L-db -u chinese.n4l
//
//******************************************************************

package main

import (
	"fmt"
	"strings"
        SST "SSTorytime"
)

var path [8][]string

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	chapter := "reminders"

	DefDeleteChapter(ctx,chapter)

	SST.Close(ctx)
}

//******************************************************************

func DefDeleteChapter(ctx SST.PoSST,chapter string) {

	ctx.DB.QueryRow("drop function DeleteChapter")

	qstr := "CREATE OR REPLACE FUNCTION DeleteChapter(chapter text)\n"+
		"RETURNS boolean AS $fn$\n" +
		"DECLARE\n" +
		"   marked    NodePtr[];\n"+
		"   autoset   NodePtr[];\n"+
		"   nnptr     NodePtr;\n"+
		"   lnk       Link;\n"+
		"   links     Link[];\n"+
		"   ed_list   Link[];\n"+
		"   oleft     text;\n"+
		"   oright    text;\n"+
		"   chaparray text[];\n"+
		"   chaplist  text;\n"+
	        "   ed_chap   text;\n"+
		"   chp       text;\n"+
	
		"BEGIN \n"+

		// First get all NPtrs contained in the chapter for deletion
		// To avoid deleting overlaps, select only the automorphic links

		"chp := Format('%%%s%%',chapter);\n"+
		"SELECT array_agg(NPtr) into autoset FROM Node WHERE Chap LIKE chp;\n"+

		"IF autoset IS NULL THEN\n"+
		"   RETURN false;\n"+
		"END IF;\n"+

		"DELETE FROM NodeArrowNode WHERE NFrom = ANY(autoset) AND NTo = ANY(autoset);\n"+

		// Look for overlapping chapters

		"oleft := Format('%%%s,%%',chapter);\n"+
		"oright := Format('%%,%s%%',chapter);\n"+

		"SELECT array_agg(NPtr) into marked FROM Node WHERE Chap LIKE oleft OR Chap LIKE oright;\n"+

		"IF marked IS NULL THEN\n"+
		"   DELETE FROM Node WHERE Chap = chapter;\n"+
		"   RETURN true;\n"+
		"END IF;\n"+

		"FOREACH nnptr IN ARRAY marked LOOP\n"+
		"   SELECT Chap into chaplist FROM Node WHERE NPtr = nnptr;\n"+
		"   chaparray = string_to_array(chaplist,',');\n"+

		// Remove the chapter reference
		"IF chaparray IS NOT NULL AND array_length(chaparray,1) > 1 THEN"+
		"   FOREACH chp IN ARRAY chaparray LOOP\n"+
		"      IF NOT chp = chapter THEN"+
		"         IF length(ed_chap) > 0 THEN\n"+
		"            ed_chap = Format('%s,%s',ed_chap,chp);\n"+
		"         ELSE"+
		"            ed_chap = chp;"+
		"         END IF;"+
		"      END IF;"+
		"   END LOOP;"+
		"   UPDATE Node SET Chap = ed_chap WHERE NPtr = nnptr;\n"+
		"   marked = array_remove(marked,nnptr);"+
		"END IF;\n"

	for st := -SST.EXPRESS; st <= SST.EXPRESS; st++ {
		qstr += fmt.Sprintf(
			
			"SELECT %s into links FROM Node WHERE NPtr = nnptr;\n"+

			"   IF links IS NOT NULL THEN\n"+
			"      ed_list = ARRAY[]::Link[];\n"+         // delete reference links
			"      FOREACH lnk in ARRAY links LOOP\n"+
			"         IF NOT lnk.Dst = ANY(marked) THEN\n"+
			"            ed_list = array_append(ed_list,lnk);\n"+
			"         END IF;\n"+
			"      END LOOP;\n"+
			"      UPDATE Node SET %s = ed_list WHERE NPtr = nnptr;\n"+
			"   END IF;\n",
			SST.STTypeDBChannel(st),SST.STTypeDBChannel(st))
	}
	
	qstr += "END LOOP;\n"+

		"DELETE FROM Node WHERE Nptr = ANY(marked);\n"+
		"DELETE FROM Node WHERE Chap = chapter;\n"+

		"RETURN true;\n" +
		"END ;\n" +
		"$fn$ LANGUAGE plpgsql;\n"

	row,err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	row.Close()


	qstr = fmt.Sprintf("select DeleteChapter('%s')",chapter)


	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error running deletechapter function:",qstr,err)
	}

	row.Close()

}

//******************************************************************

func RemoveFromStringList(list,remove string) string {

	split := strings.Split(list,",")

	var retvar string

	for s := 0; s < len(split); s++ {

		lc := strings.ToLower(split[s])
		lt := strings.ToLower(remove)

// this won't work if there are accents
		if !strings.Contains(lc,lt) {
			retvar += split[s] + ","
		}
	}

	strings.Trim(retvar,",")
	fmt.Println("XX",retvar)
	return retvar
}

//******************************************************************

func UpdateDBNode(ctx SST.PoSST,nptr SST.NodePtr,edited string,list []SST.NodePtr) {

	node := SST.GetDBNodeByNodePtr(ctx,nptr)

	fmt.Println("\nEdit",node.Chap,"with",edited)

	for st := 0; st < SST.ST_TOP; st++ {
		for dst := range node.I[st] {
			fmt.Println("  delete lenk",node.I[st][dst].Dst)
		}
	}
}

//******************************************************************

func DeleteDBNodeArrowNode(ctx SST.PoSST,nptr SST.NodePtr) {

	fmt.Println("Remove NodeArrowNode",nptr)
}

//******************************************************************

func DeleteDBNode(ctx SST.PoSST,nptr SST.NodePtr) {

	fmt.Println("DELETE nptr",nptr)
}





