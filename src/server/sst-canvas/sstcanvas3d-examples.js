/**
 * SSTCanvas3D Library Usage Example
 * 
 * This demonstrates how to use the extracted SSTCanvas3D graphics library
 * to create 3D visualizations without modifying the original main.js file.
 */
import { SSTCanvas3D } from './sstcanvas3d.js';
// Example 1: Basic setup
function createBasicVisualization()
{
    // Initialize the 3D canvas
    const canvas3d = new SSTCanvas3D('canvasContainer', {
        width: 800,
        height: 600,
        mobile: window.innerWidth < 450 ? 0.5 : 1
    });

    // Clear and draw a grid
    canvas3d.clear();
    canvas3d.drawWelcomeImage();

    return canvas3d;
}

// Example 2: Creating nodes and arrows
function createNodeExample(canvas3d)
{
    // Draw different types of nodes
    canvas3d.drawEvent(0, 0, 0);           // Red event node at origin
    canvas3d.drawThing(1, 0, 0);           // Green thing node
    canvas3d.drawConcept(-1, 0, 0);        // Blue concept node

    // Add labels
    canvas3d.drawLabel(0, 0, 0, "Main Event", 14, "black");
    canvas3d.drawLabel(1, 0, 0, "Related Thing", 12, "darkgreen");
    canvas3d.drawLabel(-1, 0, 0, "Concept", 12, "darkblue");

    // Draw arrows between nodes
    canvas3d.drawLeadsToArrow(0, 0, 0, 1, 0, 0);      // Event leads to thing
    canvas3d.drawExpressesArrow(-1, 0, 0, 0, 0, 0);   // Concept expresses event
}

// Example 3: Interactive 3D visualization
function createInteractiveVisualization()
{
    const canvas3d = new SSTCanvas3D('interactiveCanvas');

    // Set up initial view
    canvas3d.setViewingAngle(Math.PI / 6, Math.PI / 8);
    canvas3d.setObserverPosition(2, 1, -2);

    // Draw a simple graph structure
    const events = [
        { x: 0, y: 0, z: 0, text: "Start" },
        { x: 1, y: 0.5, z: 0, text: "Process" },
        { x: 2, y: 1, z: 0, text: "Result" }
    ];

    // Plot events and connections
    for (let i = 0; i < events.length; i++)
    {
        const event = events[i];
        canvas3d.drawEvent(event.x, event.y, event.z);
        canvas3d.drawLabel(event.x, event.y, event.z, event.text, 12, "black");

        if (i > 0)
        {
            const prev = events[i - 1];
            canvas3d.drawLeadsToArrow(prev.x, prev.y, prev.z, event.x, event.y, event.z);
        }
    }

    return canvas3d;
}

// Example 4: Animation loop
function animateVisualization(canvas3d)
{
    let angle = 0;

    function animate()
    {
        // Clear canvas
        canvas3d.clear();

        // Update viewing angle for rotation
        canvas3d.setViewingAngle(angle, Math.PI / 9);

        // Redraw grid and content
        canvas3d.drawGrid(0, 0, 2);
        createNodeExample(canvas3d);

        // Update angle
        angle += 0.01;

        // Continue animation
        requestAnimationFrame(animate);
    }

    animate();
}

// Example 5: Responsive canvas
function createResponsiveVisualization()
{
    const canvas3d = new SSTCanvas3D('responsiveCanvas');

    // Handle window resize
    window.addEventListener('resize', () =>
    {
        canvas3d.resizeCanvas(window.innerWidth * 0.8, window.innerHeight * 0.6);
        // Redraw content after resize
        canvas3d.clear();
        canvas3d.drawWelcomeImage();
        createNodeExample(canvas3d);
    });

    return canvas3d;
}

// Example 6: Using with semantic spacetime data
function plotSemanticSpacetimeGraph(canvas3d, eventData)
{
    /**
     * Expected eventData format:
     * {
     *   XYZ: { X: number, Y: number, Z: number },
     *   Text: string,
     *   Orbits: [
     *     [{ XYZ: {X,Y,Z}, OOO: {X,Y,Z}, Text: string }], // Im3
     *     [{ XYZ: {X,Y,Z}, OOO: {X,Y,Z}, Text: string }], // Im2
     *     [{ XYZ: {X,Y,Z}, OOO: {X,Y,Z}, Text: string }], // Im1
     *     [{ XYZ: {X,Y,Z}, OOO: {X,Y,Z}, Text: string }], // In0
     *     [{ XYZ: {X,Y,Z}, OOO: {X,Y,Z}, Text: string }], // Il1
     *     [{ XYZ: {X,Y,Z}, OOO: {X,Y,Z}, Text: string }], // Ic2
     *     [{ XYZ: {X,Y,Z}, OOO: {X,Y,Z}, Text: string }]  // Ie3
     *   ]
     * }
     */

    if (eventData && eventData.length > 0)
    {
        let lastEvent = null;

        eventData.forEach(event =>
        {
            canvas3d.plotGraphics(event, lastEvent);
            lastEvent = event;
        });
    }
}

// Usage examples - uncomment to test
/*
// Basic usage
document.addEventListener('DOMContentLoaded', () => {
    const canvas3d = createBasicVisualization();
    createNodeExample(canvas3d);
});

// Interactive usage
document.addEventListener('DOMContentLoaded', () => {
    const canvas3d = createInteractiveVisualization();
    animateVisualization(canvas3d);
});

// Responsive usage
document.addEventListener('DOMContentLoaded', () => {
    createResponsiveVisualization();
});
*/