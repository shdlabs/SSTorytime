# SSTCanvas3D Graphics Library & Demo

This directory contains the complete SSTCanvas3D graphics library and demonstration files, extracted from the main SSTorytime application to provide a reusable 3D visualization system for semantic spacetime networks.

## 📁 Directory Structure

```
sst-canvas/
├── README.md                      # This file - overview and usage guide
├── sstcanvas3d.js                # Main graphics library (extracted from main.js)
├── sstcanvas3d-examples.js       # Usage examples and patterns
├── README-sstcanvas3d.md         # Complete library API documentation
├── chapter-visualization.html    # Interactive demo application (orbits)
├── chapter-toc-demo.html         # NEW: Table of Contents visualization
├── chapter-data.json             # Real data from \chapter \limit 25 API
├── chapter-data-0.json           # Real TOC data from extraction script
├── chapter-data-enhanced.json    # Enhanced demo data with full orbits
├── extract-chapter-data.js       # Node.js data extraction script
└── extract-chapter-data.sh       # Bash data extraction script
```

## 🎯 Quick Start

### View the Demo

```bash
# Serve the files via HTTP
cd /home/alex/SSTorytime/src/server/sst-canvas
python3 -m http.server 8000

# Choose your visualization:
# http://localhost:8000/chapter-visualization.html     # Orbit-based events
# http://localhost:8000/chapter-toc-demo.html         # Table of Contents
```

### Extract Fresh Data

```bash
# Make sure SSTorytime server is running
cd /home/alex/SSTorytime
make run

# Extract data (choose one method)
cd src/server/sst-canvas
./extract-chapter-data.sh        # Bash script
node extract-chapter-data.js     # Node.js script
```

## 🔧 Library Components

### Core Library (`sstcanvas3d.js`)
- **Purpose**: Complete 3D graphics engine for semantic spacetime visualization
- **Features**: Canvas management, 3D transforms, node/arrow rendering, orbit visualization
- **Dependencies**: None (pure JavaScript)
- **Compatibility**: Modern browsers with HTML5 Canvas support

### Examples (`sstcanvas3d-examples.js`)
- **Purpose**: Demonstration patterns and usage examples
- **Content**: Basic setup, node creation, animation, responsive design
- **Use Case**: Learning how to integrate the library

### Demo Applications

#### Orbit Visualization (`chapter-visualization.html`)
- **Purpose**: Interactive visualization of event orbits and semantic relationships
- **Features**: Real-time controls, data switching, animation, statistics
- **Data Sources**: Event-based JSON with orbit relationships

#### Table of Contents (`chapter-toc-demo.html`)
- **Purpose**: 3D visualization of chapter structure and context fragments
- **Features**: Chapter selection, context filtering, TOC navigation
- **Data Sources**: Real TOC data from `\chapter \limit 25` API calls

## 📊 Data Integration

### Real Data (`chapter-data.json`, `chapter-data-0.json`)
- Source: SSTorytime API `/searchN4L` with `\chapter \limit 25`
- Format: Standard SSTorytime JSON response
- Content: Actual semantic spacetime events and relationships

#### Two Data Types:

**Orbit Data** (`chapter-data-enhanced.json`):
```json
{
  "Response": "Orbits",
  "Content": [
    {
      "Text": "Event description",
      "XYZ": {"X": 0.4, "Y": 0, "Z": 0},
      "Orbits": [/* 7 arrays of semantic relationships */]
    }
  ]
}
```

**TOC Data** (`chapter-data-0.json`):
```json
{
  "Response": "TOC",
  "Content": [
    {
      "Chapter": "Chapter title",
      "XYZ": {"X": 0.4, "Y": 0, "Z": 0},
      "Context": [/* Context fragments */],
      "Single": [/* Individual contexts */],
      "Common": [/* Shared contexts */]
    }
  ]
}
```

### Enhanced Data (`chapter-data-enhanced.json`)
- Source: Extended demo data with full orbit relationships
- Format: Same as real data but with populated orbits
- Content: Multiple events with all 7 semantic relationship types

### Data Extraction Scripts
- **`extract-chapter-data.sh`**: Bash script for quick data updates
- **`extract-chapter-data.js`**: Node.js script with better error handling
- Both scripts automatically test multiple search queries and save results

## 🎮 Interactive Features

### Visualization Controls
- **Rotation Speed**: Adjust automatic 3D rotation (0-0.05 rad/frame)
- **View Angle**: Manual control of viewing perspective
- **Animation Toggle**: Pause/resume automatic rotation
- **View Reset**: Return to default camera position
- **Center View**: Focus on coordinate origin

### Visual Elements
- 🔴 **Events**: Red spheres representing main chapter events
- 🟢 **Things**: Green spheres for objects and containers
- 🔵 **Concepts**: Blue spheres for properties and abstract concepts
- ➡️ **Arrows**: Colored directional arrows for semantic relationships
- 📝 **Labels**: Text annotations with event descriptions
- 🗂️ **Grid**: 3D coordinate reference system

### Semantic Relationships
- **Im3** (-3): "is a property expressed by" → Orange arrows
- **Im2** (-2): "is contained by" → Blue arrows
- **Im1** (-1): "comes from" → Red arrows
- **In0** (0): "is near/similar to" → Gray arrows
- **Il1** (+1): "leads to" → Red arrows
- **Ic2** (+2): "contains" → Blue arrows
- **Ie3** (+3): "expresses property" → Orange arrows

## 🔗 Integration Examples

### Basic Library Usage
```html
<!DOCTYPE html>
<html>
<head><title>SSTCanvas3D Example</title></head>
<body>
    <div id="container"></div>
    <script src="sstcanvas3d.js"></script>
    <script>
        const canvas = new SSTCanvas3D('container');
        canvas.drawEvent(0, 0, 0);
        canvas.drawLabel(0, 0, 0, "Hello World", 14, "white");
    </script>
</body>
</html>
```

### Data Visualization
```javascript
// Load and visualize chapter data
async function visualizeData() {
    const response = await fetch('chapter-data.json');
    const data = await response.json();
    
    const canvas = new SSTCanvas3D('container');
    data.Content.forEach(event => {
        canvas.plotGraphics(event, lastEvent);
        lastEvent = event;
    });
}
```

### Custom Animation
```javascript
const canvas = new SSTCanvas3D('container');
let angle = 0;

function animate() {
    canvas.clear();
    canvas.setViewingAngle(angle, Math.PI / 9);
    // ... render content ...
    angle += 0.01;
    requestAnimationFrame(animate);
}
```

## 🛠️ Development

### File Dependencies
- `chapter-visualization.html` → `sstcanvas3d.js` (main library)
- `chapter-visualization.html` → `chapter-data*.json` (data files)
- `extract-chapter-data.*` → SSTorytime server on localhost:8080

### Customization Points
- **Visual Styling**: Modify colors, sizes, and shapes in `sstcanvas3d.js`
- **Interaction**: Add mouse/touch controls to the visualization
- **Data Processing**: Extend orbit rendering for custom relationship types
- **Performance**: Implement level-of-detail for large datasets

### Testing Changes
```bash
# Test the library
python3 -m http.server 8000
# Open: http://localhost:8000/chapter-visualization.html

# Test data extraction
./extract-chapter-data.sh
# Check generated files for new data
```

## 📈 Performance Notes

### Optimization Tips
- Use `clear()` before each frame when animating
- Limit orbit depth for complex datasets
- Consider using `requestAnimationFrame()` for smooth animation
- Mobile devices automatically use 0.5x scaling factor

### Browser Compatibility
- **Required**: HTML5 Canvas, ES6+ features
- **Recommended**: Chrome, Firefox, Safari, Edge
- **Mobile**: iOS Safari, Chrome Mobile (with automatic scaling)

## 🎨 Visual Fidelity

This library maintains **exact visual compatibility** with the original SSTorytime main.js implementation:

- ✅ Same 3D coordinate transformations
- ✅ Identical node shapes and colors
- ✅ Same arrow geometry and styling
- ✅ Matching perspective projection
- ✅ Equivalent depth perception effects

The only difference is the modular, reusable architecture that makes these graphics capabilities available as a standalone library.

## 📚 Documentation

### Complete API Reference
See `README-sstcanvas3d.md` for the full library API documentation including:
- Constructor options and methods
- 3D coordinate system details
- Drawing primitive reference
- Semantic spacetime data formats
- Browser compatibility information

### Usage Examples
See `sstcanvas3d-examples.js` for practical implementation patterns:
- Basic setup and configuration
- Node and arrow creation
- Animation loops and interaction
- Responsive design techniques
- Data visualization workflows

## 🔄 Updates and Maintenance

### Keeping Data Fresh
Run the extraction scripts regularly to get updated chapter data:
```bash
# Weekly data refresh
cd /home/alex/SSTorytime/src/server/sst-canvas
./extract-chapter-data.sh
```

### Library Updates
The library is extracted from the stable main.js codebase. For updates:
1. Check for changes in `/home/alex/SSTorytime/src/server/public/main.js`
2. Re-extract graphics functions if needed
3. Test compatibility with existing visualizations

### Contributing
When adding features:
1. Maintain API compatibility
2. Update documentation
3. Test on mobile devices
4. Verify visual fidelity with original
5. Add examples for new functionality

This directory provides everything needed to create semantic spacetime visualizations using the proven SSTorytime graphics engine in a modular, reusable form.