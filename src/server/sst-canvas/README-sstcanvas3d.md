# SSTCanvas3D Graphics Library

This library extracts all canvas and 3D graphics functionality from the SSTorytime server's `main.js` file, providing a reusable and modular graphics engine for semantic spacetime visualization.

## Features

- **3D Coordinate Transformations**: Convert 3D world coordinates to 2D screen coordinates with perspective projection
- **Canvas Management**: Create, resize, and manage HTML5 canvas elements
- **Node Rendering**: Draw different types of semantic nodes (events, things, concepts)
- **Arrow Drawing**: Render directional arrows representing semantic relationships
- **Grid Visualization**: Display 3D coordinate grids for spatial reference
- **Interactive Controls**: Support for viewing angle and observer position changes
- **Mobile Responsive**: Automatic scaling for mobile devices

## Quick Start

### Basic Setup

```html
<!DOCTYPE html>
<html>
<head>
    <title>SSTCanvas3D Example</title>
</head>
<body>
    <div id="canvasContainer"></div>
    <script src="sstcanvas3d.js"></script>
    <script>
        // Create a new 3D canvas
        const canvas3d = new SSTCanvas3D('canvasContainer', {
            width: 800,
            height: 600
        });
        
        // Draw a grid and some nodes
        canvas3d.drawGrid(0, 0, 2);
        canvas3d.drawEvent(0, 0, 0);
        canvas3d.drawLabel(0, 0, 0, "Hello World", 14, "black");
    </script>
</body>
</html>
```

## API Reference

### Constructor

```javascript
new SSTCanvas3D(containerId, options)
```

**Parameters:**
- `containerId` (string): ID of the HTML element to contain the canvas
- `options` (object): Configuration options
  - `width` (number): Canvas width in pixels (default: window.innerWidth)
  - `height` (number): Canvas height in pixels (default: window.innerHeight)
  - `mobile` (number): Mobile scaling factor (default: auto-detected)

### Core Methods

#### Canvas Management

```javascript
canvas3d.clear()                    // Clear the entire canvas
canvas3d.resizeCanvas(width, height) // Resize canvas to new dimensions
canvas3d.getCanvas()                // Get the HTML canvas element
canvas3d.getContext()               // Get the 2D rendering context
canvas3d.destroy()                  // Remove canvas and clean up
```

#### 3D Coordinate System

```javascript
canvas3d.transformX(x, y, z)        // Convert 3D X to screen X
canvas3d.transformY(x, y, z)        // Convert 3D Y to screen Y
canvas3d.horizon(x, y, z)           // Calculate distance from observer
canvas3d.alpha(x, y, z)             // Calculate transparency based on distance
```

#### Viewing Controls

```javascript
canvas3d.setViewingAngle(theta, phi) // Set rotation angles
canvas3d.setObserverPosition(x, y, z) // Set observer position
```

#### Drawing Primitives

```javascript
// 3D Lines and Shapes
canvas3d.drawLine3D(x0, y0, z0, x1, y1, z1, color, thickness)
canvas3d.drawGrid(centerX, centerZ, length)
canvas3d.drawLabel(x, y, z, text, size, color)

// Nodes
canvas3d.drawNode(x, y, z, radius, colorInner, colorOuter)
canvas3d.drawEvent(x, y, z)         // Red event node
canvas3d.drawThing(x, y, z)         // Green thing node  
canvas3d.drawConcept(x, y, z)       // Blue concept node

// Arrows
canvas3d.drawArrow(x0, y0, z0, x1, y1, z1, color, thickness)
canvas3d.drawLeadsToArrow(x0, y0, z0, x1, y1, z1)     // Temporal causation
canvas3d.drawContainsArrow(x0, y0, z0, x1, y1, z1)    // Spatial containment
canvas3d.drawExpressesArrow(x0, y0, z0, x1, y1, z1)   // Property expression
canvas3d.drawNearArrow(x0, y0, z0, x1, y1, z1)        // Similarity/proximity
```

#### Semantic Spacetime

```javascript
canvas3d.plotGraphics(event, lastEvent) // Plot complete event with orbits
canvas3d.renderOrbits(event)            // Render semantic relationships
canvas3d.drawPath(arrowType, x0, y0, z0, x1, y1, z1) // Draw typed connection
```

## Semantic Spacetime Indices

The library supports the following semantic relationship types:

| Index | Name | Description |
|-------|------|-------------|
| Im3 | "is a property expressed by" | Property attribution (-3) |
| Im2 | "is contained by" | Spatial containment (-2) |  
| Im1 | "comes from" | Temporal causation (-1) |
| In0 | "is near/similar to" | Proximity/similarity (0) |
| Il1 | "leads to" | Temporal causation (+1) |
| Ic2 | "contains" | Spatial containment (+2) |
| Ie3 | "expresses property" | Property attribution (+3) |

## Event Data Format

For semantic spacetime visualization, events should follow this structure:

```javascript
const event = {
    XYZ: { X: 1.0, Y: 0.5, Z: 0.0 },     // 3D position
    Text: "Event description",             // Display text
    Orbits: [                             // Semantic relationships (7 arrays)
        [],  // Im3: properties expressed by this event
        [],  // Im2: containers of this event  
        [],  // Im1: sources/causes of this event
        [],  // In0: similar/nearby events
        [],  // Il1: events this leads to
        [],  // Ic2: events contained by this
        []   // Ie3: properties expressed by this event
    ]
};
```

Each orbit array contains neighbor objects:

```javascript
const neighbor = {
    XYZ: { X: 2.0, Y: 1.0, Z: 0.5 },     // Neighbor position
    OOO: { X: 1.0, Y: 0.5, Z: 0.0 },     // Origin position (for arrows)
    Text: "Neighbor description"           // Display text
};
```

## Examples

### Simple Node Graph

```javascript
const canvas3d = new SSTCanvas3D('container');

// Draw nodes
canvas3d.drawEvent(0, 0, 0);
canvas3d.drawThing(1, 0, 0);
canvas3d.drawConcept(-1, 0, 0);

// Add labels
canvas3d.drawLabel(0, 0, 0, "Event", 12, "black");
canvas3d.drawLabel(1, 0, 0, "Thing", 12, "black");
canvas3d.drawLabel(-1, 0, 0, "Concept", 12, "black");

// Connect with arrows
canvas3d.drawLeadsToArrow(0, 0, 0, 1, 0, 0);
canvas3d.drawExpressesArrow(-1, 0, 0, 0, 0, 0);
```

### Animated Rotation

```javascript
const canvas3d = new SSTCanvas3D('container');
let angle = 0;

function animate() {
    canvas3d.clear();
    canvas3d.setViewingAngle(angle, Math.PI / 9);
    canvas3d.drawGrid(0, 0, 2);
    // ... draw your content ...
    angle += 0.01;
    requestAnimationFrame(animate);
}

animate();
```

### Responsive Canvas

```javascript
const canvas3d = new SSTCanvas3D('container');

window.addEventListener('resize', () => {
    canvas3d.resizeCanvas(window.innerWidth * 0.8, window.innerHeight * 0.6);
    // Redraw content after resize
    redrawVisualization();
});
```

## Browser Compatibility

- Modern browsers supporting HTML5 Canvas
- ES6+ features (classes, arrow functions, destructuring)
- Tested on Chrome, Firefox, Safari, Edge

## Integration

The library can be used with various module systems:

```javascript
// CommonJS
const SSTCanvas3D = require('./sstcanvas3d.js');

// AMD
define(['./sstcanvas3d.js'], function(SSTCanvas3D) { ... });

// Global (browser)
// Just include the script tag - SSTCanvas3D will be available globally
```

## Performance Notes

- Use `clear()` before redrawing to avoid visual artifacts
- For animations, use `requestAnimationFrame()` for smooth performance
- Mobile scaling is automatically applied for devices with width < 450px
- Transparency calculations improve depth perception but may impact performance on complex scenes

## Relationship to Original Code

This library extracts the following functions from the original `main.js`:

- Canvas creation and management functions
- 3D coordinate transformation mathematics
- Node and arrow drawing routines
- Grid and primitive shape rendering
- Semantic spacetime visualization logic

The original `main.js` file remains unchanged and continues to work as before. This library provides the same functionality in a reusable, modular form.