import React, { useRef, useEffect, useCallback, useState } from 'react';

// Types for the canvas system
interface CanvasNode {
  x3d: number;
  y3d: number;
  z3d: number;
  x2d: number;
  y2d: number;
  radius: number;
  text: string;
  type: 'event' | 'concept' | 'thing';
  searchQuery: string;
}

interface Canvas3DProps {
  className?: string;
  width?: number;
  height?: number;
  onNodeClick?: (searchQuery: string) => void;
}

// 3D Canvas Configuration
const CANVAS_CONFIG = {
  // Default view parameters
  THETA: Math.PI / 9,
  PHI: Math.PI / 9,
  SCALE: 0.9,
  OBS_X: 1,
  OBS_Y: 0.5,
  OBS_Z: -1,
  
  // Visual settings
  GRID_SIZE: 1,
  MOB_SCALE: 1, // Mobile scaling factor
};

export const Canvas3D: React.FC<Canvas3DProps> = ({ 
  className = "", 
  width, 
  height,
  onNodeClick 
}) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const ctxRef = useRef<CanvasRenderingContext2D | null>(null);
  const interactiveNodesRef = useRef<CanvasNode[]>([]);
  const tooltipRef = useRef<HTMLDivElement | null>(null);
  
  const [dimensions, setDimensions] = useState({
    width: width || 800,
    height: height || 600
  });
  
  const [viewParams, setViewParams] = useState({
    ORGX: (width || 800) / 2,
    ORGY: (height || 600) / 2,
    THETA: CANVAS_CONFIG.THETA,
    PHI: CANVAS_CONFIG.PHI,
    SCALE: CANVAS_CONFIG.SCALE,
    WIDTH: width || 800,
    HEIGHT: height || 600
  });

  // Update mobile scale factor based on screen size
  const getMobScale = useCallback(() => {
    if (typeof window === 'undefined') return 1;
    if (window.innerWidth < 450) return 0.4;
    if (window.innerWidth < 1300) return 0.6;
    return 1;
  }, []);

  // 3D to 2D coordinate transformation functions
  const horizon = useCallback((x: number, y: number, z: number) => {
    return Math.sqrt(
      (x - CANVAS_CONFIG.OBS_X) * (x - CANVAS_CONFIG.OBS_X) +
      (y - CANVAS_CONFIG.OBS_Y) * (y - CANVAS_CONFIG.OBS_Y) +
      (z - CANVAS_CONFIG.OBS_Z) * (z - CANVAS_CONFIG.OBS_Z)
    );
  }, []);

  const tx = useCallback((x: number, y: number, z: number) => {
    const scale = (viewParams.SCALE * viewParams.WIDTH) / (1 + 1.2 * horizon(x, y, z));
    return viewParams.ORGX + scale * (x * Math.cos(viewParams.THETA) + z * Math.cos(viewParams.PHI));
  }, [viewParams, horizon]);

  const ty = useCallback((x: number, y: number, z: number) => {
    const scale = (viewParams.SCALE * viewParams.WIDTH) / (1 + 1.5 * horizon(x, y, z));
    return viewParams.HEIGHT - viewParams.ORGY - scale * (y + z * Math.sin(viewParams.PHI) - x * Math.sin(viewParams.THETA));
  }, [viewParams, horizon]);

  // Drawing functions
  const drawLine = useCallback((
    x0: number, y0: number, z0: number,
    xp: number, yp: number, zp: number,
    color: string, thickness: number
  ) => {
    const ctx = ctxRef.current;
    if (!ctx) return;

    ctx.save();
    ctx.beginPath();
    
    const xb = tx(x0, y0, z0);
    const yb = ty(x0, y0, z0);
    const xe = tx(xp, yp, zp);
    const ye = ty(xp, yp, zp);

    ctx.moveTo(xb, yb);
    ctx.lineTo(xe, ye);
    ctx.strokeStyle = color;
    ctx.lineWidth = thickness;
    ctx.stroke();
    ctx.closePath();
    ctx.restore();
  }, [tx, ty]);

  const drawNode = useCallback((
    x: number, y: number, z: number,
    radius: number, color1: string, color2: string
  ) => {
    const ctx = ctxRef.current;
    if (!ctx) return;

    ctx.save();
    ctx.beginPath();
    
    const x0 = tx(x, y, z);
    const y0 = ty(x, y, z);
    const r = (radius * 1.6) / horizon(x, y, z);

    const grad = ctx.createLinearGradient(x0, y0, x0 + r, y0 + r);
    grad.addColorStop(0, color2);
    grad.addColorStop(1, color1);

    ctx.arc(x0, y0, r, 0, Math.PI * 2);
    ctx.fillStyle = grad;
    ctx.fill();
    ctx.restore();

    return { x2d: x0, y2d: y0, radius: r };
  }, [tx, ty, horizon]);

  const drawArrow = useCallback((
    x0: number, y0: number, z0: number,
    xp: number, yp: number, zp: number,
    color: string, thickness: number
  ) => {
    const ctx = ctxRef.current;
    if (!ctx) return;

    // Draw the line
    drawLine(x0, y0, z0, xp, yp, zp, color, thickness);

    // Draw arrowhead
    ctx.save();
    const frx = tx(x0, y0, z0);
    const fry = ty(x0, y0, z0);
    const tox = tx(xp, yp, zp);
    const toy = ty(xp, yp, zp);
    const scale = 1.1 - zp;
    const angle = Math.atan2(toy - fry, tox - frx);
    const headangle = Math.PI / 12;
    const headlen = 12 * scale;
    const noderadius = 10 * scale;

    ctx.beginPath();
    ctx.strokeStyle = color;
    ctx.lineWidth = thickness;
    
    // Left arrowhead line
    ctx.moveTo(
      tox - noderadius * Math.cos(angle),
      toy - noderadius * Math.sin(angle)
    );
    ctx.lineTo(
      tox - headlen * Math.cos(angle - headangle),
      toy - headlen * Math.sin(angle - headangle)
    );
    
    // Right arrowhead line
    ctx.moveTo(
      tox - noderadius * Math.cos(angle),
      toy - noderadius * Math.sin(angle)
    );
    ctx.lineTo(
      tox - headlen * Math.cos(angle + headangle),
      toy - headlen * Math.sin(angle + headangle)
    );
    
    ctx.stroke();
    ctx.restore();
  }, [drawLine, tx, ty]);

  const drawGrid = useCallback((length: number) => {
    const ctx = ctxRef.current;
    if (!ctx) return;

    const mob = getMobScale();
    ctx.save();

    // Draw grid lines
    for (let xi = -length; xi <= length; xi += 0.1) {
      drawLine(xi, 0, -length, xi, 0, length, "lightgrey", 1 * mob);
    }

    for (let zi = -length; zi <= length; zi += 0.1) {
      drawLine(-length, 0, zi, length, 0, zi, "lightgrey", 1 * mob);
    }

    // Draw axis lines
    drawLine(-length / 2, 0, 0, 0, 0, 0, "lightgrey", 1 * mob);
    drawLine(0, 0, -length / 2, 0, 0, 0, "lightgrey", 1 * mob);
    drawLine(0, -length / 2, 0, 0, length, 0, "lightgrey", 1 * mob);

    ctx.restore();
  }, [drawLine, getMobScale]);

  // Node drawing functions
  const drawEvent = useCallback((x: number, y: number, z: number) => {
    const mob = getMobScale();
    return drawNode(x, y, z, 6 * mob, "darkred", "red");
  }, [drawNode, getMobScale]);

  const drawConcept = useCallback((x: number, y: number, z: number) => {
    const mob = getMobScale();
    return drawNode(x, y, z, 4 * mob, "darkblue", "lightblue");
  }, [drawNode, getMobScale]);

  const drawThing = useCallback((x: number, y: number, z: number) => {
    const mob = getMobScale();
    return drawNode(x, y, z, 4 * mob, "darkgreen", "lightgreen");
  }, [drawNode, getMobScale]);

  // Arrow drawing functions for relationships
  const drawLeadsTo = useCallback((x0: number, y0: number, z0: number, xp: number, yp: number, zp: number) => {
    const mob = getMobScale();
    drawArrow(x0, y0, z0, xp, yp, zp, "darkred", 3 * mob);
  }, [drawArrow, getMobScale]);

  const drawContains = useCallback((x0: number, y0: number, z0: number, xp: number, yp: number, zp: number) => {
    const mob = getMobScale();
    drawArrow(x0, y0, z0, xp, yp, zp, "lightblue", 2 * mob);
  }, [drawArrow, getMobScale]);

  const drawExpresses = useCallback((x0: number, y0: number, z0: number, xp: number, yp: number, zp: number) => {
    const mob = getMobScale();
    drawArrow(x0, y0, z0, xp, yp, zp, "orange", 2 * mob);
  }, [drawArrow, getMobScale]);

  const drawNear = useCallback((x0: number, y0: number, z0: number, xp: number, yp: number, zp: number) => {
    const mob = getMobScale();
    drawArrow(x0, y0, z0, xp, yp, zp, "darkgrey", 1 * mob);
  }, [drawArrow, getMobScale]);

  // Interactive node management
  const addInteractiveNode = useCallback((
    x3d: number, y3d: number, z3d: number,
    text: string, type: 'event' | 'concept' | 'thing',
    searchQuery: string
  ) => {
    const x2d = tx(x3d, y3d, z3d);
    const y2d = ty(x3d, y3d, z3d);
    const mob = getMobScale();
    const radius = (4 * mob * 1.6) / horizon(x3d, y3d, z3d);

    interactiveNodesRef.current.push({
      x2d, y2d, radius, text, type, searchQuery, x3d, y3d, z3d
    });
  }, [tx, ty, horizon, getMobScale]);

  const clearInteractiveNodes = useCallback(() => {
    interactiveNodesRef.current = [];
  }, []);

  // Public drawing functions that add interactivity
  const drawInteractiveEvent = useCallback((x: number, y: number, z: number, text: string, searchQuery: string) => {
    addInteractiveNode(x, y, z, text, "event", searchQuery);
    drawEvent(x, y, z);
  }, [addInteractiveNode, drawEvent]);

  const drawInteractiveConcept = useCallback((x: number, y: number, z: number, text: string, searchQuery: string) => {
    addInteractiveNode(x, y, z, text, "concept", searchQuery);
    drawConcept(x, y, z);
  }, [addInteractiveNode, drawConcept]);

  const drawInteractiveThing = useCallback((x: number, y: number, z: number, text: string, searchQuery: string) => {
    addInteractiveNode(x, y, z, text, "thing", searchQuery);
    drawThing(x, y, z);
  }, [addInteractiveNode, drawThing]);

  // Canvas initialization and management
  const initCanvas = useCallback(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    ctxRef.current = ctx;
    clearInteractiveNodes();

    // Set canvas size
    canvas.width = dimensions.width;
    canvas.height = dimensions.height;

    // Update view parameters
    setViewParams(prev => ({
      ...prev,
      WIDTH: dimensions.width,
      HEIGHT: dimensions.height,
      ORGX: dimensions.width / 2,
      ORGY: dimensions.height / 2
    }));

    return ctx;
  }, [dimensions, clearInteractiveNodes]);

  // Clear canvas
  const clearCanvas = useCallback(() => {
    const ctx = ctxRef.current;
    const canvas = canvasRef.current;
    if (!ctx || !canvas) return;

    ctx.clearRect(0, 0, canvas.width, canvas.height);
  }, []);

  // Demo visualization function
  const drawDemo = useCallback(() => {
    clearCanvas();
    clearInteractiveNodes();
    
    // Draw grid
    drawGrid(CANVAS_CONFIG.GRID_SIZE);

    // Draw some demo nodes and relationships
    drawInteractiveConcept(0, 0, 0, "Central Concept", "\\central");
    drawInteractiveEvent(-0.5, 0, 0.3, "Event A", "\\event_a");
    drawInteractiveThing(0.5, 0, -0.3, "Thing B", "\\thing_b");

    // Draw relationships
    drawExpresses(0, 0, 0, -0.5, 0, 0.3);
    drawContains(0, 0, 0, 0.5, 0, -0.3);
    drawLeadsTo(-0.5, 0, 0.3, 0.5, 0, -0.3);

    // Draw orbital concepts
    const orbitRadius = 0.8;
    for (let a = 0; a < 2 * Math.PI; a += Math.PI / 3) {
      const x = orbitRadius * Math.cos(a);
      const z = orbitRadius * Math.sin(a);
      drawInteractiveConcept(x, 0, z, `Orbit ${Math.round(a * 3 / Math.PI)}`, `\\orbit_${Math.round(a * 3 / Math.PI)}`);
      drawNear(0, 0, 0, x, 0, z);
    }
  }, [clearCanvas, clearInteractiveNodes, drawGrid, drawInteractiveConcept, drawInteractiveEvent, drawInteractiveThing, drawExpresses, drawContains, drawLeadsTo, drawNear]);

  // Tooltip management
  const createTooltip = useCallback(() => {
    if (tooltipRef.current) return tooltipRef.current;

    const tooltip = document.createElement('div');
    tooltip.style.position = 'absolute';
    tooltip.style.background = 'rgba(0, 0, 0, 0.8)';
    tooltip.style.color = 'white';
    tooltip.style.padding = '5px 10px';
    tooltip.style.borderRadius = '4px';
    tooltip.style.fontSize = '12px';
    tooltip.style.pointerEvents = 'none';
    tooltip.style.zIndex = '1000';
    tooltip.style.maxWidth = '300px';
    tooltip.style.wordWrap = 'break-word';
    tooltip.style.display = 'none';
    document.body.appendChild(tooltip);

    tooltipRef.current = tooltip;
    return tooltip;
  }, []);

  const showTooltip = useCallback((text: string, x: number, y: number) => {
    const tooltip = createTooltip();
    tooltip.textContent = text;
    tooltip.style.left = (x + 10) + 'px';
    tooltip.style.top = (y - 5) + 'px';
    tooltip.style.display = 'block';
  }, [createTooltip]);

  const hideTooltip = useCallback(() => {
    if (tooltipRef.current) {
      tooltipRef.current.style.display = 'none';
    }
  }, []);

  // Mouse interaction
  const getMousePos = useCallback((canvas: HTMLCanvasElement, evt: MouseEvent) => {
    const rect = canvas.getBoundingClientRect();
    const scaleX = canvas.width / rect.width;
    const scaleY = canvas.height / rect.height;

    return {
      x: (evt.clientX - rect.left) * scaleX,
      y: (evt.clientY - rect.top) * scaleY
    };
  }, []);

  const findNodeAtPosition = useCallback((mouseX: number, mouseY: number) => {
    for (let i = interactiveNodesRef.current.length - 1; i >= 0; i--) {
      const node = interactiveNodesRef.current[i];
      const dx = mouseX - node.x2d;
      const dy = mouseY - node.y2d;
      const distance = Math.sqrt(dx * dx + dy * dy);

      if (distance <= node.radius) {
        return node;
      }
    }
    return null;
  }, []);

  // Setup canvas interactivity
  const setupCanvasInteractivity = useCallback(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    let hoveredNode: CanvasNode | null = null;

    const handleMouseMove = (evt: MouseEvent) => {
      const mousePos = getMousePos(canvas, evt);
      const node = findNodeAtPosition(mousePos.x, mousePos.y);

      if (node !== hoveredNode) {
        if (hoveredNode) {
          hideTooltip();
          canvas.style.cursor = 'default';
        }

        hoveredNode = node;

        if (node) {
          showTooltip(node.text, evt.clientX, evt.clientY);
          canvas.style.cursor = 'pointer';
        }
      } else if (node) {
        showTooltip(node.text, evt.clientX, evt.clientY);
      }
    };

    const handleMouseLeave = () => {
      hideTooltip();
      canvas.style.cursor = 'default';
      hoveredNode = null;
    };

    const handleClick = (evt: MouseEvent) => {
      const mousePos = getMousePos(canvas, evt);
      const node = findNodeAtPosition(mousePos.x, mousePos.y);

      if (node && node.searchQuery && onNodeClick) {
        onNodeClick(node.searchQuery);
      }
    };

    canvas.addEventListener('mousemove', handleMouseMove);
    canvas.addEventListener('mouseleave', handleMouseLeave);
    canvas.addEventListener('click', handleClick);

    return () => {
      canvas.removeEventListener('mousemove', handleMouseMove);
      canvas.removeEventListener('mouseleave', handleMouseLeave);
      canvas.removeEventListener('click', handleClick);
    };
  }, [getMousePos, findNodeAtPosition, hideTooltip, showTooltip, onNodeClick]);

  // Handle dimension changes
  useEffect(() => {
    if (!width && !height) {
      const updateDimensions = () => {
        if (canvasRef.current) {
          const rect = canvasRef.current.getBoundingClientRect();
          setDimensions({
            width: Math.max(400, rect.width || 800),
            height: Math.max(300, rect.height || 600)
          });
        }
      };

      updateDimensions();
      window.addEventListener('resize', updateDimensions);
      return () => window.removeEventListener('resize', updateDimensions);
    }
  }, [width, height]);

  // Initialize canvas and setup
  useEffect(() => {
    const ctx = initCanvas();
    if (!ctx) return;

    const cleanup = setupCanvasInteractivity();
    
    // Draw initial demo
    setTimeout(() => {
      drawDemo();
    }, 100);

    return () => {
      cleanup?.();
      if (tooltipRef.current) {
        document.body.removeChild(tooltipRef.current);
        tooltipRef.current = null;
      }
    };
  }, [initCanvas, setupCanvasInteractivity, drawDemo]);

  // Expose drawing functions for external use
  const canvasAPI = {
    clearCanvas,
    clearInteractiveNodes,
    drawGrid,
    drawInteractiveEvent,
    drawInteractiveConcept,
    drawInteractiveThing,
    drawLeadsTo,
    drawContains,
    drawExpresses,
    drawNear,
    drawDemo
  };

  // Attach API to canvas ref for external access
  useEffect(() => {
    if (canvasRef.current) {
      (canvasRef.current as any).canvasAPI = canvasAPI;
    }
  }, [canvasAPI]);

  return (
    <div className={`relative ${className}`}>
      <canvas
        ref={canvasRef}
        width={dimensions.width}
        height={dimensions.height}
        className="block border border-gray-300 rounded-lg"
        style={{
          width: dimensions.width,
          height: dimensions.height,
          background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)'
        }}
      />
    </div>
  );
};

export default Canvas3D;