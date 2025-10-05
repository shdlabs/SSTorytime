/**
 * SSTorytime Canvas3D Graphics Library
 * 
 * Extracted from main.js server package to create a reusable graphics library
 * for 3D canvas visualization, graph rendering, and interactive node display.
 * 
 * This library provides:
 * - 3D coordinate transformations
 * - Canvas management and drawing
 * - Node and arrow rendering
 * - Graph visualization
 * - Grid and geometric primitives
 * 
 * @author Mark Burgess & Alex (extracted)
 * @version 1.0.0
 */

class SSTCanvas3D
{
    constructor(containerId, options = {})
    {
        // Configuration
        this.containerId = containerId;
        this.options = {
            width: options.width || window.innerWidth,
            height: options.height || window.innerHeight,
            mobile: options.mobile || (window.innerWidth < 650 ? 0.4 : 1),
            ...options
        };

        // Initialize canvas and context
        this.canvas = null;
        this.ctx = null;
        this.mob = this.options.mobile;

        // 3D viewing parameters
        this.width = this.options.width;
        this.height = this.options.height;
        this.orgX = this.width / 2;
        this.orgY = this.height / 2;
        this.theta = Math.PI / 11;
        this.phi = Math.PI / 11;
        this.scale = 0.8;
        this.obsX = 1;
        this.obsY = 0.5;
        this.obsZ = -1;
        this.vpX = 0;
        this.vpY = 4;
        this.vpZ = 8;

        // Semantic spacetime indices
        this.ST_INDICES = {
            Im3: 0, // "is a property expressed by"
            Im2: 1, // "is contained by"
            Im1: 2, // "comes from"
            In0: 3, // "is near/similar to"
            Il1: 4, // "leads to"
            Ic2: 5, // "contains"
            Ie3: 6  // "expresses property"
        };

        this.ARROW_LABELS = [
            "is a property expressed by",
            "is contained by",
            "comes from",
            "is near/similar to",
            "leads to",
            "contains",
            "expresses property"
        ];

        this.init();
    }

    /**
     * Initialize the canvas and graphics context
     */
    init()
    {
        this.createCanvas();
        this.ctx = this.canvas.getContext("2d");
        this.ctx.beginPath();
    }

    /**
     * Create canvas element and append to container
     * @returns {HTMLCanvasElement} The created canvas
     */
    createCanvas()
    {
        // Remove existing canvas if present
        const oldCanvas = document.getElementById("sstCanvas3D");
        if (oldCanvas)
        {
            oldCanvas.remove();
        }

        const parent = document.getElementById(this.containerId);
        if (!parent)
        {
            throw new Error(`Container with ID '${this.containerId}' not found`);
        }

        this.canvas = document.createElement("canvas");
        this.canvas.id = "sstCanvas3D";
        this.canvas.width = this.width;
        this.canvas.height = this.height;

        parent.appendChild(this.canvas);
        return this.canvas;
    }

    /**
     * Resize canvas to fit container or specified dimensions
     */
    async resizeCanvas(width, height)
    {
        if (width && height)
        {
            this.width = width;
            this.height = height;
        }

        const scale = Math.min(
            window.innerWidth / this.canvas.width,
            window.innerHeight / this.canvas.height
        );

        this.canvas.style.width = `${Math.round(scale * this.canvas.width)}px`;
        this.canvas.style.height = `${Math.round(scale * this.canvas.height)}px`;

        // Update coordinate system
        this.orgX = this.width / 2;
        this.orgY = this.height / 2;
    }

    /**
     * Clear the entire canvas
     */
    clear()
    {
        this.ctx.clearRect(0, 0, this.width, this.height);
    }

    /**
     * Calculate horizon distance for depth perception
     * @param {number} x - X coordinate
     * @param {number} y - Y coordinate  
     * @param {number} z - Z coordinate
     * @returns {number} Distance from observer
     */
    horizon(x, y, z)
    {
        return Math.sqrt(
            (x - this.obsX) * (x - this.obsX) +
            (y - this.obsY) * (y - this.obsY) +
            (z - this.obsZ) * (z - this.obsZ)
        );
    }

    /**
     * Calculate alpha transparency based on distance
     * @param {number} x - X coordinate
     * @param {number} y - Y coordinate
     * @param {number} z - Z coordinate
     * @returns {number} Alpha value between 0.1 and 1
     */
    alpha(x, y, z)
    {
        let alpha = 1.5 / Math.sqrt(
            (x - this.obsX) * (x - this.obsX) +
            (y - this.obsY) * (y - this.obsY) +
            (z - this.obsZ) * (z - this.obsZ)
        );

        if (alpha > 1) return 1;
        if (alpha < 0.1) return 0.1;
        return alpha;
    }

    /**
     * Transform 3D X coordinate to 2D screen X
     * @param {number} x - 3D X coordinate
     * @param {number} y - 3D Y coordinate
     * @param {number} z - 3D Z coordinate
     * @returns {number} Screen X coordinate
     */
    transformX(x, y, z)
    {
        const scale = (this.scale * this.width) / (1 + 1.2 * this.horizon(x, y, z));
        return this.orgX + scale * (x * Math.cos(this.theta) + z * Math.cos(this.phi));
    }

    /**
     * Transform 3D Y coordinate to 2D screen Y
     * @param {number} x - 3D X coordinate
     * @param {number} y - 3D Y coordinate
     * @param {number} z - 3D Z coordinate
     * @returns {number} Screen Y coordinate
     */
    transformY(x, y, z)
    {
        const scale = (this.scale * this.width) / (1 + 1.5 * this.horizon(x, y, z));
        return this.height - this.orgY - scale * (y + z * Math.sin(this.phi) - x * Math.sin(this.theta));
    }

    /**
     * Draw a 3D line
     * @param {number} x0 - Start X coordinate
     * @param {number} y0 - Start Y coordinate
     * @param {number} z0 - Start Z coordinate
     * @param {number} x1 - End X coordinate
     * @param {number} y1 - End Y coordinate
     * @param {number} z1 - End Z coordinate
     * @param {string} color - Line color
     * @param {number} thickness - Line thickness
     */
    drawLine3D(x0, y0, z0, x1, y1, z1, color, thickness)
    {
        this.ctx.save();
        this.ctx.beginPath();

        const xStart = this.transformX(x0, y0, z0);
        const yStart = this.transformY(x0, y0, z0);
        const xEnd = this.transformX(x1, y1, z1);
        const yEnd = this.transformY(x1, y1, z1);

        this.ctx.moveTo(xStart, yStart);
        this.ctx.lineTo(xEnd, yEnd);
        this.ctx.strokeStyle = color;
        this.ctx.lineWidth = thickness;
        this.ctx.stroke();
        this.ctx.closePath();
        this.ctx.restore();
    }

    /**
     * Draw a 3D grid
     * @param {number} centerX - Grid center X
     * @param {number} centerZ - Grid center Z  
     * @param {number} length - Grid size
     */
    drawGrid(centerX, centerZ, length)
    {
        this.ctx.save();

        // Draw grid lines
        for (let i = -length; i <= length; i += length / 25)
        {
            // Parallel to X axis
            this.drawLine3D(
                centerX - length / 2, 0, centerZ + i,
                centerX + length / 2, 0, centerZ + i,
                "oklch(27.1% 0.000 0.0)", 0.1 * this.mob
            );
            // Parallel to Z axis
            this.drawLine3D(
                centerX + i, 0, centerZ - length / 2,
                centerX + i, 0, centerZ + length / 2,
                "oklch(25.1% 0.009 0.0)", 0.1 * this.mob
            );
        }

        // Draw coordinate axes
        this.drawLine3D(-length / 2, 0, 0, length / 2, 0, 0, "oklch(59.1% 0.293 322.896)", 3 * this.mob);
        this.drawLine3D(0, 0, -length / 2, 0, 0, length / 2, "oklch(58.8% 0.158 241.966)", 3 * this.mob);
        this.drawLine3D(0, -length / 2, 0, 0, length / 2, 0, "oklch(69.6% 0.17 162.48)", 3 * this.mob);

        this.ctx.restore();
    }

    /**
     * Draw a text label at 3D position
     * @param {number} x - 3D X coordinate
     * @param {number} y - 3D Y coordinate
     * @param {number} z - 3D Z coordinate
     * @param {string} text - Label text
     * @param {number} size - Font size
     * @param {string} color - Text color
     */
    drawLabel(x, y, z, text, size, color)
    {
        this.ctx.save();

        const font = `bold ${size * this.mob}px sans-serif`;
        const screenX = this.transformX(x, y, z) + 30;
        const screenY = this.transformY(x, y, z);

        this.ctx.beginPath();
        this.ctx.font = font;
        this.ctx.fillStyle = color;
        this.ctx.fillText(text, screenX, screenY);

        this.ctx.restore();
    }

    /**
     * Draw a 3D node (sphere)
     * @param {number} x - 3D X coordinate
     * @param {number} y - 3D Y coordinate
     * @param {number} z - 3D Z coordinate
     * @param {number} radius - Node radius
     * @param {string} colorInner - Inner color
     * @param {string} colorOuter - Outer color
     */
    drawNode(x, y, z, radius, colorInner, colorOuter)
    {
        this.ctx.save();
        this.ctx.beginPath();

        const screenX = this.transformX(x, y, z);
        const screenY = this.transformY(x, y, z);
        const adjustedRadius = (radius * 1.6) / this.horizon(x, y, z);

        // Create gradient
        const grad = this.ctx.createLinearGradient(
            screenX, screenY,
            screenX + adjustedRadius, screenY + adjustedRadius
        );
        grad.addColorStop(0, colorOuter);
        grad.addColorStop(1, colorInner);

        this.ctx.arc(screenX, screenY, adjustedRadius, 0, Math.PI * 2);
        this.ctx.fillStyle = grad;
        this.ctx.fill();

        this.ctx.restore();
    }

    /**
     * Draw a 3D arrow
     * @param {number} x0 - Start X coordinate
     * @param {number} y0 - Start Y coordinate
     * @param {number} z0 - Start Z coordinate
     * @param {number} x1 - End X coordinate
     * @param {number} y1 - End Y coordinate
     * @param {number} z1 - End Z coordinate
     * @param {string} color - Arrow color
     * @param {number} thickness - Arrow thickness
     */
    drawArrow(x0, y0, z0, x1, y1, z1, color, thickness)
    {
        this.ctx.save();

        // Draw the line
        this.drawLine3D(x0, y0, z0, x1, y1, z1, color, thickness);

        // Draw arrowhead
        const fromX = this.transformX(x0, y0, z0);
        const fromY = this.transformY(x0, y0, z0);
        const toX = this.transformX(x1, y1, z1);
        const toY = this.transformY(x1, y1, z1);

        const scale = 1.1 - z1;
        const angle = Math.atan2(toY - fromY, toX - fromX);
        const headAngle = Math.PI / 12;
        const headLen = 12 * scale;
        const nodeRadius = 10 * scale;

        this.ctx.beginPath();
        this.ctx.strokeStyle = color;
        this.ctx.lineWidth = thickness;

        // Draw arrowhead
        this.ctx.moveTo(
            toX - nodeRadius * Math.cos(angle),
            toY - nodeRadius * Math.sin(angle)
        );
        this.ctx.lineTo(
            toX - headLen * Math.cos(angle - headAngle),
            toY - headLen * Math.sin(angle - headAngle)
        );
        this.ctx.moveTo(
            toX - nodeRadius * Math.cos(angle),
            toY - nodeRadius * Math.sin(angle)
        );
        this.ctx.lineTo(
            toX - headLen * Math.cos(angle + headAngle),
            toY - headLen * Math.sin(angle + headAngle)
        );
        this.ctx.stroke();
        this.ctx.restore();
    }

    /**
     * Node type renderers
     */
    drawEvent(x, y, z)
    {
        this.drawNode(x, y, z, 6 * this.mob, "darkred", "red");
    }

    drawThing(x, y, z)
    {
        this.drawNode(x, y, z, 4 * this.mob, "darkgreen", "lightgreen");
    }

    drawConcept(x, y, z)
    {
        this.drawNode(x, y, z, 4 * this.mob, "darkblue", "lightblue");
    }

    /**
     * Arrow type renderers
     */
    drawLeadsToArrow(x0, y0, z0, x1, y1, z1)
    {
        this.drawArrow(x0, y0, z0, x1, y1, z1, "darkred", 3 * this.mob);
    }

    drawContainsArrow(x0, y0, z0, x1, y1, z1)
    {
        this.drawArrow(x0, y0, z0, x1, y1, z1, "lightblue", 2 * this.mob);
    }

    drawExpressesArrow(x0, y0, z0, x1, y1, z1)
    {
        this.drawArrow(x0, y0, z0, x1, y1, z1, "orange", 2 * this.mob);
    }

    drawNearArrow(x0, y0, z0, x1, y1, z1)
    {
        this.drawArrow(x0, y0, z0, x1, y1, z1, "darkgrey", 1 * this.mob);
    }

    /**
     * Draw path between events based on arrow type
     * @param {number} arrowType - Type of arrow connection
     * @param {number} thisX - Current node X
     * @param {number} thisY - Current node Y
     * @param {number} thisZ - Current node Z
     * @param {number} lastX - Previous node X
     * @param {number} lastY - Previous node Y
     * @param {number} lastZ - Previous node Z
     */
    drawPath(arrowType, thisX, thisY, thisZ, lastX, lastY, lastZ)
    {
        if (lastX === 0 && lastY === 0 && lastZ === 0) return;

        switch (arrowType)
        {
        case -3:
            this.drawExpressesArrow(thisX, thisY, thisZ, lastX, lastY, lastZ);
            break;
        case -2:
            this.drawContainsArrow(thisX, thisY, thisZ, lastX, lastY, lastZ);
            break;
        case -1:
            this.drawLeadsToArrow(thisX, thisY, thisZ, lastX, lastY, lastZ);
            break;
        case 0:
            this.drawLeadsToArrow(thisX, thisY, thisZ, lastX, lastY, lastZ);
            break;
        case 1:
            this.drawLeadsToArrow(lastX, lastY, lastZ, thisX, thisY, thisZ);
            break;
        case 2:
            this.drawContainsArrow(lastX, lastY, lastZ, thisX, thisY, thisZ);
            break;
        case 3:
            this.drawExpressesArrow(lastX, lastY, lastZ, thisX, thisY, thisZ);
            break;
        default:
            console.warn("Unknown arrow type:", arrowType);
        }
    }

    /**
     * Plot a complete event with its orbits
     * @param {Object} event - Event object with XYZ coordinates and orbits
     * @param {Object} lastEvent - Previous event for path drawing
     */
    plotGraphics(event, lastEvent)
    {
        const tx = event.XYZ.X;
        const ty = event.XYZ.Y;
        const tz = event.XYZ.Z;

        this.drawEvent(tx, ty, tz);
        this.drawLabel(tx, ty, tz, event.Text.slice(0, 25), 12, "black");

        if (lastEvent && lastEvent !== event)
        {
            const lx = lastEvent.XYZ.X;
            const ly = lastEvent.XYZ.Y;
            const lz = lastEvent.XYZ.Z;

            this.drawLeadsToArrow(lx, ly, lz, tx, ty, tz);
        }

        // Render orbits based on semantic spacetime indices
        this.renderOrbits(event);
    }

    /**
     * Render event orbits based on semantic spacetime relationships
     * @param {Object} event - Event with orbits array
     */
    renderOrbits(event)
    {
        const orbitRenderers = [
            { index: this.ST_INDICES.Il1, nodeType: 'event', arrowFunc: 'drawLeadsToArrow', direction: 'outgoing' },
            { index: this.ST_INDICES.Im1, nodeType: 'event', arrowFunc: 'drawLeadsToArrow', direction: 'incoming' },
            { index: this.ST_INDICES.Ic2, nodeType: 'thing', arrowFunc: 'drawContainsArrow', direction: 'outgoing' },
            { index: this.ST_INDICES.Im2, nodeType: 'thing', arrowFunc: 'drawContainsArrow', direction: 'incoming' },
            { index: this.ST_INDICES.Ie3, nodeType: 'concept', arrowFunc: 'drawExpressesArrow', direction: 'outgoing' },
            { index: this.ST_INDICES.Im3, nodeType: 'concept', arrowFunc: 'drawExpressesArrow', direction: 'incoming' },
            { index: this.ST_INDICES.In0, nodeType: 'event', arrowFunc: 'drawNearArrow', direction: 'outgoing' }
        ];

        orbitRenderers.forEach(renderer =>
        {
            const orbit = event.Orbits[renderer.index];
            if (!orbit) return;

            orbit.forEach(neighbor =>
            {
                // Draw the node
                switch (renderer.nodeType)
                {
                case 'event':
                    this.drawEvent(neighbor.XYZ.X, neighbor.XYZ.Y, neighbor.XYZ.Z);
                    break;
                case 'thing':
                    this.drawThing(neighbor.XYZ.X, neighbor.XYZ.Y, neighbor.XYZ.Z);
                    break;
                case 'concept':
                    this.drawConcept(neighbor.XYZ.X, neighbor.XYZ.Y, neighbor.XYZ.Z);
                    break;
                }

                // Draw the arrow
                if (renderer.direction === 'outgoing')
                {
                    this[renderer.arrowFunc](
                        neighbor.OOO.X, neighbor.OOO.Y, neighbor.OOO.Z,
                        neighbor.XYZ.X, neighbor.XYZ.Y, neighbor.XYZ.Z
                    );
                } else
                {
                    this[renderer.arrowFunc](
                        neighbor.XYZ.X, neighbor.XYZ.Y, neighbor.XYZ.Z,
                        neighbor.OOO.X, neighbor.OOO.Y, neighbor.OOO.Z
                    );
                }
            });
        });
    }

    /**
     * Draw a welcome/demo visualization
     */
    drawWelcomeImage()
    {
        this.drawGrid(0, 0, 1);
        // Add additional welcome graphics as needed
    }

    /**
     * Update viewing angle
     * @param {number} theta - Rotation around Y axis
     * @param {number} phi - Rotation around X axis
     */
    setViewingAngle(theta, phi)
    {
        this.theta = theta;
        this.phi = phi;
    }

    /**
     * Update observer position
     * @param {number} x - Observer X position
     * @param {number} y - Observer Y position
     * @param {number} z - Observer Z position
     */
    setObserverPosition(x, y, z)
    {
        this.obsX = x;
        this.obsY = y;
        this.obsZ = z;
    }

    /**
     * Get current canvas element
     * @returns {HTMLCanvasElement} The canvas element
     */
    getCanvas()
    {
        return this.canvas;
    }

    /**
     * Get canvas 2D context
     * @returns {CanvasRenderingContext2D} The 2D context
     */
    getContext()
    {
        return this.ctx;
    }

    /**
     * Destroy the canvas and clean up
     */
    destroy()
    {
        if (this.canvas)
        {
            this.canvas.remove();
            this.canvas = null;
            this.ctx = null;
        }
    }
}

// Export for module systems or global use
if (typeof module !== 'undefined' && module.exports)
{
    module.exports = SSTCanvas3D;
} else if (typeof define === 'function' && define.amd)
{
    define([], function () { return SSTCanvas3D; });
} else
{
    window.SSTCanvas3D = SSTCanvas3D;
}