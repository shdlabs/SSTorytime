//******************************************************************
//
// Experiment with old CFEngine context gathering approach
//
//******************************************************************

package main

import (
	"fmt"
	"time"
	"regexp"
	"io/ioutil"

//	"strings"

        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	now := time.Now()
	c,slot := DoNowt(now)

	fmt.Println("TIME_CLASSES",c,"\nSLOT",slot)


	txtcls := ContextFromFormatFile("/home/mark/Laptop/Work/SST/data_samples/MobyDick.dat")

	fmt.Println("FILE_SCAN_CLASSES",txtcls)

	SST.Close(ctx)
}

// ****************************************************************************
// Semantic 2D time
// ****************************************************************************

var GR_DAY_TEXT = []string{
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday",
        "Saturday",
        "Sunday",
    }
        
var GR_MONTH_TEXT = []string{
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December",
}
        
var GR_SHIFT_TEXT = []string{
        "Night",
        "Morning",
        "Afternoon",
        "Evening",
    }

// For second resolution Unix time

const CF_MONDAY_MORNING = 345200
const CF_MEASURE_INTERVAL = 5*60
const CF_SHIFT_INTERVAL = 6*3600

const MINUTES_PER_HOUR = 60
const SECONDS_PER_MINUTE = 60
const SECONDS_PER_HOUR = (60 * SECONDS_PER_MINUTE)
const SECONDS_PER_DAY = (24 * SECONDS_PER_HOUR)
const SECONDS_PER_WEEK = (7 * SECONDS_PER_DAY)
const SECONDS_PER_YEAR = (365 * SECONDS_PER_DAY)
const HOURS_PER_SHIFT = 6
const SECONDS_PER_SHIFT = (HOURS_PER_SHIFT * SECONDS_PER_HOUR)
const SHIFTS_PER_DAY = 4
const SHIFTS_PER_WEEK = (4*7)

// ****************************************************************************
// Semantic spacetime timeslots
// ****************************************************************************

func DoNowt(then time.Time) (string,string) {

	//then := given.UnixNano()

	// Time on the torus (donut/doughnut) (CFEngine style)
	// The argument is a Golang time unit e.g. then := time.Now()
	// Return a db-suitable keyname reflecting the coarse-grained SST time
	// The function also returns a printable summary of the time

	year := fmt.Sprintf("Yr%d",then.Year())
	month := GR_MONTH_TEXT[int(then.Month())-1]
	day := then.Day()
	hour := fmt.Sprintf("Hr%02d",then.Hour())
	mins := fmt.Sprintf("Min%02d",then.Minute())
	quarter := fmt.Sprintf("Q%d",then.Minute()/15 + 1)
	shift :=  fmt.Sprintf("%s",GR_SHIFT_TEXT[then.Hour()/6])

	//secs := then.Second()
	//nano := then.Nanosecond()

	dayname := then.Weekday()
	dow := fmt.Sprintf("%.3s",dayname)
	daynum := fmt.Sprintf("Day%d",day)

	// 5 minute resolution capture
        interval_start := (then.Minute() / 5) * 5
        interval_end := (interval_start + 5) % 60
        minD := fmt.Sprintf("Min%02d_%02d",interval_start,interval_end)

	var when string = fmt.Sprintf("%s,%s,%s,%s,%s at %s %s %s %s",shift,dayname,daynum,month,year,hour,mins,quarter,minD)
	var key string = fmt.Sprintf("%s:%s:%s",dow,hour,minD)

	return when, key
}

// ****************************************************************************

func GetUnixTimeKey(now int64) string {

	// Time on the torus (donut/doughnut) (CFEngine style)
	// The argument is in traditional UNIX "time_t" unit e.g. then := time.Unix()
	// This is a simple wrapper to DoNowt() returning only a db-suitable keyname

	t := time.Unix(now, 0)
	_,slot := DoNowt(t)

	return slot
}

//******************************************************************
// Read text file
//******************************************************************

func ContextFromFormatFile(name string) string {

	file := ReadFormatFile(name)
	return file
}

// *****************************************************************

func ReadFormatFile(filename string) string {

	// Read a string and strip out characters that can't be used in kenames
	// to yield a "pure" text for n-gram classification, with fewer special chars
	// The text marks end of sentence with a # for later splitting

	content, _ := ioutil.ReadFile(filename)

	// Start by stripping HTML / XML tags before para-split
	// if they haven't been removed already

	m1 := regexp.MustCompile("<[^>]*>") 
	cleaned := m1.ReplaceAllString(string(content),";") 
	return cleaned
}

