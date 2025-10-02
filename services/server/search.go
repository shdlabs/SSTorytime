package server

import (
	"fmt"
	"net/http"
	"strings"

	SST "github.com/shdlabs/SSTorytime/services/sstorytime"
)

// SearchHandler handles N4L search requests
func SearchHandler(ctx SST.PoSST) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST", "GET":
			name := r.FormValue("name")
			nclass := r.FormValue("nclass")
			ncptr := r.FormValue("ncptr")
			chapcontext := r.FormValue("chapcontext")

			if name == "\\lastnptr" {
				if chapcontext != "" && chapcontext != "any" {
					UpdateLastSawSection(ctx, w, r, chapcontext)
				}
				UpdateLastSawNPtr(ctx, w, r, nclass, ncptr, chapcontext)
				return
			}

			if name == "" && len(nclass) > 0 && len(ncptr) > 0 {
				var a, b int
				fmt.Sscanf(nclass, "%d", &a)
				fmt.Sscanf(ncptr, "%d", &b)
				nstr := fmt.Sprintf("(%d,%d)", a, b)
				name = name + nstr
			}

			fmt.Println("\\nReceived command:", name)

			ambient, key, _ := SST.GetTimeContext()

			if len(name) == 0 || name == "\\remind" {
				name = "any \\chapter reminders \\context " + key + " " + ambient
			}

			if len(name) == 0 || name == "\\help" {
				name = "\\notes \\chapter \"help and search\" \\limit 40"
			}

			search := SST.DecodeSearchField(name)
			HandleSearch(ctx, search, name, w, r)

		default:
			http.Error(w, "Not supported", http.StatusMethodNotAllowed)
		}
	}
}

// Helper function to convert string slice to string for logging
func stringSliceToString(slice []string) string {
	if slice == nil {
		return "nil"
	}
	return strings.Join(slice, ", ")
}
