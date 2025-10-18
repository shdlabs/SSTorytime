#!/bin/bash

# Demo script for N4L validator watch mode
# Shows how validation updates in real-time as you edit files

VALIDATOR="/home/alex/SSTorytime/cmd/n4l-validate/n4l-validate"
TEST_FILE="/tmp/n4l_watch_demo.n4l"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "========================================="
echo "  N4L Validator Watch Mode Demo"
echo "========================================="
echo ""
echo "This demo shows the validator watching a file"
echo "and re-validating automatically on changes."
echo ""

# Create initial valid file
cat > "$TEST_FILE" << 'EOF'
- Watch mode demo

+ demo context

 First node (contain) second node
 
 Second node (belong) first node
EOF

echo -e "${GREEN}Starting validator in watch mode...${NC}"
echo "File: $TEST_FILE"
echo ""
echo "The validator will watch for changes."
echo "In another terminal, try editing: $TEST_FILE"
echo ""
echo "Example edits to try:"
echo "  1. Add a valid node: Third node (hasX) value"
echo "  2. Add an error: Broken quote: \"oops"
echo "  3. Fix the error and save again"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop watching${NC}"
echo ""

# Start watching
$VALIDATOR -w -v "$TEST_FILE"
