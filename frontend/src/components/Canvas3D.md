# Canvas3D Component

A React TypeScript component that provides 3D visualization functionality extracted from Mark's original main.js canvas implementation.

## Features

- **3D Coordinate System**: Full 3D to 2D projection with observer positioning
- **Interactive Nodes**: Clickable nodes with hover tooltips
- **Grid Rendering**: 3D grid background for spatial reference  
- **Multiple Node Types**: Support for Events, Concepts, and Things with different visual styles
- **Relationship Arrows**: Directional arrows showing relationships (LeadsTo, Contains, Expresses, Near)
- **Mouse Interaction**: Hover effects, tooltips, and click handling
- **Responsive Design**: Automatically adapts to container size
- **Demo Mode**: Built-in demonstration with sample nodes and relationships

## Usage

### Basic Usage

```tsx
import Canvas3D from './components/Canvas3D';

function MyComponent() {
  const handleNodeClick = (searchQuery: string) => {
    console.log('Node clicked:', searchQuery);
  };

  return (
    <Canvas3D 
      width={800}
      height={600}
      onNodeClick={handleNodeClick}
      className="my-canvas"
    />
  );
}
```

### Props

- `className?: string` - Additional CSS classes
- `width?: number` - Canvas width (defaults to 800px)
- `height?: number` - Canvas height (defaults to 600px)  
- `onNodeClick?: (searchQuery: string) => void` - Callback when interactive nodes are clicked

### API Access

The component exposes a canvas API through the canvas ref for programmatic control:

```tsx
const canvasRef = useRef<HTMLCanvasElement>(null);

// Access API methods
const canvasAPI = (canvasRef.current as any)?.canvasAPI;
if (canvasAPI) {
  canvasAPI.clearCanvas();
  canvasAPI.drawDemo();
  canvasAPI.drawInteractiveConcept(0, 0, 0, "My Concept", "\\search_query");
}
```

### Available API Methods

- `clearCanvas()` - Clear the entire canvas
- `clearInteractiveNodes()` - Remove all interactive node data
- `drawGrid(length: number)` - Draw 3D grid with specified size
- `drawInteractiveEvent(x, y, z, text, searchQuery)` - Draw clickable event node
- `drawInteractiveConcept(x, y, z, text, searchQuery)` - Draw clickable concept node  
- `drawInteractiveThing(x, y, z, text, searchQuery)` - Draw clickable thing node
- `drawLeadsTo(x0, y0, z0, x1, y1, z1)` - Draw "leads to" relationship arrow
- `drawContains(x0, y0, z0, x1, y1, z1)` - Draw "contains" relationship arrow
- `drawExpresses(x0, y0, z0, x1, y1, z1)` - Draw "expresses" relationship arrow
- `drawNear(x0, y0, z0, x1, y1, z1)` - Draw "near" relationship arrow
- `drawDemo()` - Draw demonstration scene with sample nodes

## 3D Coordinate System

The component uses a 3D coordinate system with:

- **Observer Position**: (1, 0.5, -1) - Controls perspective view
- **Viewing Angles**: THETA and PHI for rotation
- **Scale Factor**: Adjustable zoom level
- **Origin**: Canvas center serves as (0,0,0) in 3D space

## Node Types & Colors

- **Events**: Red nodes (6px radius) - Represent actions or happenings
- **Concepts**: Blue nodes (4px radius) - Represent ideas or abstractions  
- **Things**: Green nodes (4px radius) - Represent objects or entities

## Relationship Types & Colors

- **LeadsTo**: Dark red arrows (3px width) - Causal relationships
- **Contains**: Light blue arrows (2px width) - Containment relationships
- **Expresses**: Orange arrows (2px width) - Expression relationships
- **Near**: Dark grey arrows (1px width) - Proximity relationships

## Integration with Visualization Component

The Canvas3D component is integrated into the main Visualization component with a toggle to switch between 2D SVG and 3D Canvas modes:

```tsx
<Visualization 
  data={searchResult} 
  mode="canvas"  // or "svg"
  onNodeClick={handleNodeClick}
/>
```

Users can switch between visualization modes using the "2D View" and "3D View" buttons in the top-left corner.