#!/bin/bash

# SSTorytime Chapter Data Extraction Script
# 
# This script extracts JSON data from the SSTorytime API using the \chapter \limit 25 command
# and saves it to the public directory for use with the canvas visualization.

echo "🔍 Extracting SSTorytime Chapter Data..."

# Check if server is running
if ! curl -s http://localhost:8080/status > /dev/null 2>&1; then
    echo "❌ Error: SSTorytime server is not running on localhost:8080"
    echo "Please start the server first with: make run"
    exit 1
fi

# Extract the data
echo "📡 Fetching data from API..."

# Try different search queries to get varied data
queries=(
    "\chapter \limit 25"
    # "any \\chapter \\limit 25" 
    # "\\notes \\chapter \\limit 25"
    # "story \\limit 25"
    # "knowledge \\limit 25"
)

for i in "${!queries[@]}"; do
    query="${queries[$i]}"
    filename="chapter-data-$i.json"
    
    echo "🔍 Trying query: $query"
    
    response=$(curl -s -X POST -F "search=$query" http://localhost:8080/searchN4L)
    
    if [ $? -eq 0 ] && [ ! -z "$response" ]; then
        echo "$response" | jq . > "$filename" 2>/dev/null
        
        if [ $? -eq 0 ]; then
            echo "✅ Saved to $filename"
            
            # Extract some stats
            events=$(echo "$response" | jq '.Content | length' 2>/dev/null)
            echo "   📊 Events found: $events"
            
            # If this is the first successful extraction, also save as default
            if [ $i -eq 0 ]; then
                cp "$filename" "chapter-data.json"
                echo "   💾 Also saved as chapter-data.json (default)"
            fi
        else
            echo "⚠️  Response received but not valid JSON"
            echo "$response" > "$filename.raw"
            echo "   Raw response saved to $filename.raw"
        fi
    else
        echo "❌ Failed to get response"
    fi
    
    echo ""
done

echo "🎯 Data extraction complete!"
echo ""
echo "📁 Files created in current directory:"
ls -la chapter-data*.json 2>/dev/null || echo "   No JSON files created"
echo ""
echo "🌐 To view the visualization:"
echo "   1. Open chapter-visualization.html in a web browser"
echo "   2. Or serve via HTTP: python3 -m http.server 8000"
echo "   3. Then visit: http://localhost:8000/chapter-visualization.html"