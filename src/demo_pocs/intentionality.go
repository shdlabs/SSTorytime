//
// Test intent assessment
//

package main

import (
	"fmt"
        SST "SSTorytime"
)

//**************************************************************
// BEGIN
//**************************************************************

func main() {

	SST.MemoryInit()
	s1 := "My heart will go on"
	s2 := "Dee doo run run run, dee doo run run"
	s3 := "I wish to go fishing"
	
	for i := 0; i < 5; i++ {
		
		SST.Fractionate2Learn(s1,len(s1),SST.STM_NGRAM_RANK,SST.N_GRAM_MIN)
		SST.Fractionate2Learn(s2,len(s2),SST.STM_NGRAM_RANK,SST.N_GRAM_MIN)
		SST.Fractionate2Learn(s3,len(s3),SST.STM_NGRAM_RANK,SST.N_GRAM_MIN)
		
		L:= 10
		i1 := SST.AssessIntent(s1,L,SST.STM_NGRAM_RANK,SST.N_GRAM_MIN)
		i2 := SST.AssessIntent(s2,L,SST.STM_NGRAM_RANK,SST.N_GRAM_MIN)
		i3 := SST.AssessIntent(s3,L,SST.STM_NGRAM_RANK,SST.N_GRAM_MIN)
		
		fmt.Println("AD HOC",i,s1,i1)
		fmt.Println("AD HOC",i,s2,i2)
		fmt.Println("AD HOC",i,s3,i3)
	}
}
