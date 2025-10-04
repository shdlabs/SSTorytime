# SSTorytime Chapter Visualization Example

This example demonstrates how to extract JSON data from the SSTorytime API using the `\chapter \limit 25` command and visualize it using the new SSTCanvas3D graphics library.

## 📁 Files Included

### Core Library Files
- **`sstcanvas3d.js`** - The extracted 3D graphics library
- **`sstcanvas3d-examples.js`** - Usage examples for the library
- **`README-sstcanvas3d.md`** - Complete library documentation

### Visualization Example
- **`chapter-visualization.html`** - Complete interactive 3D visualization
- **`chapter-data.json`** - Real data from `\chapter \limit 25` command
- **`chapter-data-enhanced.json`** - Enhanced demo data with full orbit relationships

### Data Extraction Tools
- **`extract-chapter-data.sh`** - Bash script to extract fresh data from API
- **`extract-chapter-data.js`** - Node.js script for data extraction

## 🚀 Quick Start

### Option 1: View the Visualization Directly

1. **Open the visualization:**
   ```bash
   # Simply open in browser
   open chapter-visualization.html
   
   # Or serve via HTTP
   python3 -m http.server 8000
   # Then visit: http://localhost:8000/chapter-visualization.html
   ```

2. **Use the controls:**
   - Switch between real and demo data
   - Adjust rotation speed and viewing angle
   - Pause/play animation
   - Reset or center the view

### Option 2: Extract Fresh Data

1. **Make sure SSTorytime server is running:**
   ```bash
   cd /home/alex/SSTorytime
   make run
   ```

2. **Extract data using bash script:**
   ```bash
   cd /home/alex/SSTorytime/src/server/public
   ./extract-chapter-data.sh
   ```

3. **Or extract using Node.js:**
   ```bash
   cd /home/alex/SSTorytime/src/server/public
   node extract-chapter-data.js
   ```

## 🎯 Features Demonstrated

### SSTCanvas3D Library Usage
- **3D Canvas Creation** - Automatic canvas setup and management
- **Coordinate Transformations** - 3D to 2D projection with perspective
- **Node Rendering** - Different node types (events, things, concepts)
- **Arrow Drawing** - Semantic relationship arrows
- **Interactive Controls** - Rotation, zoom, and view controls

### Semantic Spacetime Visualization
- **Event Networks** - Main events connected by temporal arrows
- **Orbit Relationships** - 7 types of semantic connections:
  - `Im3`: "is a property expressed by" (incoming concepts)
  - `Im2`: "is contained by" (incoming containers)
  - `Im1`: "comes from" (incoming causation)
  - `In0`: "is near/similar to" (proximity)
  - `Il1`: "leads to" (outgoing causation)
  - `Ic2`: "contains" (outgoing containers)
  - `Ie3`: "expresses property" (outgoing concepts)

### Real-time Data Integration
- **API Integration** - Direct connection to SSTorytime searchN4L endpoint
- **JSON Processing** - Automatic parsing of server responses
- **Dynamic Updates** - Switch between data sources in real-time

## 📊 Data Format

The visualization expects data in this format:

```json
{
  "Response": "Orbits",
  "Content": [
    {
      "Text": "Event description",
      "L": 38,
      "Chap": "chapter_name",
      "Context": "context_name",
      "XYZ": { "X": 0.4, "Y": 0, "Z": 0 },
      "Orbits": [
        [{"Text": "...", "XYZ": {...}, "OOO": {...}}],  // Im3
        [{"Text": "...", "XYZ": {...}, "OOO": {...}}],  // Im2
        [{"Text": "...", "XYZ": {...}, "OOO": {...}}],  // Im1
        [{"Text": "...", "XYZ": {...}, "OOO": {...}}],  // In0
        [{"Text": "...", "XYZ": {...}, "OOO": {...}}],  // Il1
        [{"Text": "...", "XYZ": {...}, "OOO": {...}}],  // Ic2
        [{"Text": "...", "XYZ": {...}, "OOO": {...}}]   // Ie3
      ]
    }
  ],
  "Time": "Fri:Hr18:Qu2-Min25_30",
  "Intent": "Chapter:reminders,context...",
  "Ambient": "N_Autumn, S_Spring, Evening..."
}
```

## 🎮 Interactive Controls

### Visualization Controls
- **Data Source** - Switch between real and demo data
- **Rotation Speed** - Control automatic rotation speed (0-0.05)
- **View Angle** - Manually adjust viewing angle
- **Animation** - Pause/play the rotation animation
- **Reset View** - Return to default viewing position
- **Center View** - Focus on coordinate origin

### Visual Elements
- 🔴 **Red Spheres** - Main events from chapters
- 🟢 **Green Spheres** - Things/objects in relationships
- 🔵 **Blue Spheres** - Concepts and properties
- ➡️ **Red Arrows** - Temporal causation (leads to)
- ➡️ **Blue Arrows** - Spatial containment
- ➡️ **Orange Arrows** - Property expression
- ➡️ **Gray Arrows** - Similarity/proximity

## 🔧 Technical Details

### Library Architecture
- **Class-based Design** - Modern ES6+ JavaScript
- **Module Support** - Works with CommonJS, AMD, and global scope
- **Mobile Responsive** - Automatic scaling for mobile devices
- **Performance Optimized** - Efficient canvas operations

### Browser Requirements
- Modern browser with HTML5 Canvas support
- ES6+ features (classes, arrow functions, fetch API)
- Recommended: Chrome, Firefox, Safari, Edge

### Server Integration
- Compatible with SSTorytime Go server
- Uses standard HTTP POST requests
- JSON response parsing
- CORS-friendly for local development

## 🎨 Customization

### Visual Styling
You can customize the visualization by modifying:
- Node colors and sizes in the SSTCanvas3D methods
- Arrow styles and thickness
- Grid appearance and density
- Label fonts and colors

### Data Processing
Extend the visualization by:
- Adding new semantic relationship types
- Implementing custom node shapes
- Creating animated transitions
- Adding interactive selection

### Example Custom Node:
```javascript
// In your custom visualization
canvas3d.drawNode(x, y, z, radius, 'purple', 'lavender');
canvas3d.drawLabel(x, y, z, 'Custom Node', 14, 'purple');
```

## 🐛 Troubleshooting

### Common Issues

1. **Blank Canvas**
   - Check browser console for errors
   - Verify data files exist and are valid JSON
   - Ensure SSTCanvas3D library loaded correctly

2. **No Data Displayed**
   - Verify SSTorytime server is running on localhost:8080
   - Check that data extraction completed successfully
   - Try the enhanced demo data as fallback

3. **Performance Issues**
   - Reduce rotation speed
   - Use fewer data points
   - Check for JavaScript errors in console

### Debug Mode
Enable debug logging by opening browser developer tools and running:
```javascript
// See current visualization state
console.log(window.chapterViz);

// Check loaded data
console.log(window.chapterViz.chapterData);

// Test canvas operations
window.chapterViz.canvas3d.clear();
window.chapterViz.canvas3d.drawEvent(0, 0, 0);
```

## 📈 Extending the Example

This example serves as a foundation for more complex visualizations:

1. **Add Animation Effects** - Smooth node transitions, orbit animations
2. **Implement Interaction** - Click to explore nodes, hover for details
3. **Create Filters** - Show/hide specific relationship types
4. **Add Temporal Views** - Time-based navigation through events
5. **Export Capabilities** - Save visualizations as images or videos

The SSTCanvas3D library provides all the tools needed for advanced semantic spacetime visualizations while maintaining the same visual fidelity as the original SSTorytime application.