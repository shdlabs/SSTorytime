package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/shdlabs/SSTorytime/services/text2n4l"
)

// Text2N4LRequest defines the structure for text processing requests
type Text2N4LRequest struct {
	Text       string  `json:"text"`
	Percentage float64 `json:"percentage"`
}

// Text2N4LResponse defines the structure for text processing responses
type Text2N4LResponse struct {
	N4LContent string `json:"n4l_content"`
	Stats      struct {
		TotalSentences    int     `json:"total_sentences"`
		SelectedSentences int     `json:"selected_sentences"`
		FinalFraction     float64 `json:"final_fraction"`
		RequestedFraction float64 `json:"requested_fraction"`
	} `json:"stats"`
	Error string `json:"error,omitempty"`
}

// enableCORS adds CORS headers to allow cross-origin requests
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// processTextHandler handles text processing requests
func processTextHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var req Text2N4LRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Text2N4LResponse{
			Error: "Invalid JSON request: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate input
	if strings.TrimSpace(req.Text) == "" {
		response := Text2N4LResponse{
			Error: "Text field is required and cannot be empty",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set default percentage if not provided or invalid
	if req.Percentage <= 0 || req.Percentage > 100 {
		req.Percentage = 10.0
	}

	// Process the text
	result, err := text2n4l.ProcessTextContent(req.Text, req.Percentage)
	if err != nil {
		response := Text2N4LResponse{
			Error: "Failed to process text: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Prepare successful response
	response := Text2N4LResponse{
		N4LContent: result.N4LContent,
	}
	response.Stats.TotalSentences = result.TotalSentences
	response.Stats.SelectedSentences = result.SelectedSentences
	response.Stats.FinalFraction = result.FinalFraction
	response.Stats.RequestedFraction = result.RequestedFraction

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// healthHandler provides a simple health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "text2n4l-standalone",
	})
}

func main() {
	// Setup routes
	http.HandleFunc("/process", processTextHandler)
	http.HandleFunc("/health", healthHandler)

	// Serve static files from the static directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	port := "3001"
	fmt.Printf("Text2N4L Standalone Server starting on port %s\n", port)
	fmt.Printf("Web interface: http://localhost:%s\n", port)
	fmt.Printf("API endpoint: http://localhost:%s/process\n", port)
	fmt.Println("Press Ctrl+C to stop the server")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
