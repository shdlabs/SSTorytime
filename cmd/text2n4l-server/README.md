# Text2N4L Standalone Tool

A standalone web application for converting text to N4L (Narrative for Learning) format using Promise Theory analysis.

## Features

- **Two-panel interface**: Input text on the left, N4L output on the right
- **Adjustable selection percentage**: Control how much of the text to process (default: 10%)
- **Real-time processing**: Instant conversion with visual feedback
- **Edit capabilities**: Modify the generated N4L output as needed
- **Copy functionality**: Easy copying of results to clipboard
- **Responsive design**: Works on desktop and mobile devices

## Quick Start

### Option 1: Run the Standalone Server

1. Build and run the server:
```bash
cd cmd/text2n4l-server
go mod tidy
go build -o text2n4l-server .
./text2n4l-server
```

2. Open your browser to http://localhost:3001

### Option 2: Build Executable

To create a standalone executable:

```bash
cd cmd/text2n4l-server
go build -o text2n4l-server .
```

Then run `./text2n4l-server` from anywhere.

## How to Use

1. **Input Text**: Paste your text in the left panel
2. **Set Parameters**: Adjust the selection percentage (1-100%)
3. **Process**: Click "Process Text" to analyze
4. **Review**: Examine the N4L output in the right panel
5. **Edit**: Modify the output as needed
6. **Copy**: Use the Copy button to save results

## About N4L Format

N4L (Narrative for Learning) format extracts high-intentionality sentences that are significant for knowledge representation and Promise Theory analysis. The algorithm uses:

- **Dynamic running assessment**: Identifies sentences with high intentionality in context
- **Static post-hoc assessment**: Analyzes sentence importance after full document processing
- **Contextual annotations**: Adds metadata for Promise Theory analysis

## API Endpoints

The server exposes these endpoints:

- `GET /` - Web interface
- `POST /process` - Text processing API
- `GET /health` - Health check

### API Usage

```bash
curl -X POST http://localhost:3001/process \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Your text here...",
    "percentage": 10
  }'
```

## Requirements

- Go 1.21+
- SSTorytime libraries (included as local modules)

## Architecture

- **Backend**: Go HTTP server with CORS support
- **Frontend**: Vanilla HTML/CSS/JavaScript (no frameworks)
- **Processing**: Uses optimized text2n4l package with Promise Theory algorithms
- **Standalone**: No external dependencies beyond Go standard library

## Configuration

The server runs on port 3001 by default. To change:

1. Edit the `port` variable in `main.go`
2. Rebuild the server
3. Update the API URL in `static/index.html` if needed

## Development

To modify the interface:

1. Edit `static/index.html` for UI changes
2. Edit `main.go` for server logic
3. Rebuild and restart the server

## Performance

The tool includes several optimizations:
- Efficient string processing with `strings.Builder`
- Smart buffering for large text inputs
- Minimal memory allocations
- Fast N4L content generation

## License

Part of the SSTorytime project.