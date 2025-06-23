//******************************************************************
//
// Experiment with old CFEngine context gathering approach
//
//******************************************************************

package main

import (
	"fmt"
	"time"
	"os"
	"os/exec"
        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	now := time.Now()
	c,slot := SST.DoNowt(now)

	input_stream := "/home/mark/Laptop/Work/SST/data_samples/MobyDick.dat"
	input_stream = "../../../org-42/roam/how-i-org.org"

	SST.ContextFromFile(input_stream)

	for n := SST.N_GRAM_MIN; n < SST.N_GRAM_MAX; n++ {
		for g := range SST.STM_NGRAM_RANK[n] {
			fmt.Println("ng",n,g)
		}
	}

	fmt.Println("TIME_CLASSES",c,"\nSLOT",slot)
	fmt.Println("INPUT_STREAM",input_stream)
	name,_ := os.Hostname()
	fmt.Println("HOST",name)

	cmd := exec.Command("ls", "-l") // "ls" is the command, "-l" is an argument
	err := cmd.Run()

	SST.Close(ctx)
}

