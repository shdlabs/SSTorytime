import React, { useState, useRef, useEffect } from 'react';
import { motion } from 'framer-motion';
import * as d3 from 'd3';
import { APIResponse, NodeOrbit, XYZ } from '../types/api';

interface VisualizationProps {
  data: APIResponse;
  className?: string;
}

interface VisualizationNode extends NodeOrbit {
  id: string;
  screenX: number;
  screenY: number;
  color: string;
  size: number;
  chapter?: string;
  type: 'context' | 'single' | 'common' | 'chapter';
}

export const Visualization: React.FC<VisualizationProps> = ({ data, className = "" }) => {
  const svgRef = useRef<SVGSVGElement>(null);
  const [nodes, setNodes] = useState<VisualizationNode[]>([]);
  const [selectedNode, setSelectedNode] = useState<string | null>(null);
  const [dimensions, setDimensions] = useState({ width: 800, height: 600 });

  // Update dimensions on resize
  useEffect(() => {
    const updateDimensions = () => {
      if (svgRef.current) {
        const rect = svgRef.current.getBoundingClientRect();
        setDimensions({
          width: Math.max(400, rect.width || 800),
          height: Math.max(300, rect.height || 600)
        });
      }
    };

    updateDimensions();
    window.addEventListener('resize', updateDimensions);
    return () => window.removeEventListener('resize', updateDimensions);
  }, []);

  // Process data into visualization nodes
  useEffect(() => {
    if (!('Content' in data) || !data.Content) return;

    const processedNodes: VisualizationNode[] = [];
    let nodeIndex = 0;

    // Handle different response types
    if (data.Response === 'Orbits' && Array.isArray(data.Content)) {
      // Simple orbits data - direct array of nodes
      data.Content.forEach((node: NodeOrbit) => {
        const processedNode = processNode(node, nodeIndex++, 'context');
        processedNodes.push(processedNode);
      });
    } else if (data.Response === 'TOC' && Array.isArray(data.Content)) {
      // Table of contents with chapters
      data.Content.forEach((chapter: any) => {
        // Add chapter as a central node
        if (chapter.Chapter) {
          const chapterNode: NodeOrbit = {
            Text: chapter.Chapter,
            XYZ: chapter.XYZ || { X: 0, Y: 0, Z: 0 },
            Reln: null
          };
          const processedChapter = processNode(chapterNode, nodeIndex++, 'chapter', chapter.Chapter);
          processedNodes.push(processedChapter);
        }

        // Add context nodes
        if (chapter.Context && Array.isArray(chapter.Context)) {
          chapter.Context.forEach((node: NodeOrbit) => {
            const processedNode = processNode(node, nodeIndex++, 'context', chapter.Chapter);
            processedNodes.push(processedNode);
          });
        }

        // Add single nodes
        if (chapter.Single && Array.isArray(chapter.Single)) {
          chapter.Single.forEach((node: NodeOrbit) => {
            const processedNode = processNode(node, nodeIndex++, 'single', chapter.Chapter);
            processedNodes.push(processedNode);
          });
        }

        // Add common nodes
        if (chapter.Common && Array.isArray(chapter.Common)) {
          chapter.Common.forEach((node: NodeOrbit) => {
            const processedNode = processNode(node, nodeIndex++, 'common', chapter.Chapter);
            processedNodes.push(processedNode);
          });
        }
      });
    } else if (Array.isArray(data.Content)) {
      // Generic array of content
      data.Content.forEach((item: any) => {
        // Try to extract nodes from various structures
        if (item.Text && item.XYZ) {
          // Direct node
          const processedNode = processNode(item, nodeIndex++, 'context');
          processedNodes.push(processedNode);
        } else if (item.Context) {
          // Has context array
          item.Context.forEach((node: NodeOrbit) => {
            const processedNode = processNode(node, nodeIndex++, 'context', item.Chapter);
            processedNodes.push(processedNode);
          });
        }
      });
    }

    setNodes(processedNodes);
  }, [data, dimensions]);

  // Helper function to process individual nodes
  const processNode = (
    node: NodeOrbit, 
    index: number, 
    type: 'context' | 'single' | 'common' | 'chapter',
    chapter?: string
  ): VisualizationNode => {
    // Convert 3D coordinates to 2D screen coordinates
    const screenX = ((node.XYZ.X + 1) / 2) * dimensions.width;
    const screenY = ((1 - node.XYZ.Y) / 2) * dimensions.height; // Flip Y axis
    
    // Calculate size based on type and Z coordinate
    let baseSize = 12;
    if (type === 'chapter') baseSize = 25;
    else if (type === 'context') baseSize = 18;
    else if (type === 'single') baseSize = 14;
    else if (type === 'common') baseSize = 20;
    
    const size = Math.max(6, Math.min(40, baseSize + (node.XYZ.Z || 0) * 10 + (node.XYZ.R || 0) * 8));
    
    // Generate color based on type and position
    let hue = (index * 137.508) % 360; // Golden angle distribution
    let saturation = 60;
    let lightness = 50;
    
    switch (type) {
      case 'chapter':
        hue = 220; // Blue for chapters
        saturation = 80;
        lightness = 60;
        break;
      case 'context':
        hue = 150; // Green for context
        saturation = 70;
        lightness = 55;
        break;
      case 'single':
        hue = 60; // Yellow for single
        saturation = 70;
        lightness = 60;
        break;
      case 'common':
        hue = 300; // Purple for common
        saturation = 70;
        lightness = 55;
        break;
    }
    
    // Adjust based on position
    saturation += (node.XYZ.Z || 0) * 15;
    lightness += (node.XYZ.Y || 0) * 15;
    
    const color = `hsl(${hue}, ${Math.max(20, Math.min(90, saturation))}%, ${Math.max(30, Math.min(80, lightness))}%)`;

    return {
      ...node,
      id: `node-${index}`,
      screenX: Math.max(size, Math.min(dimensions.width - size, screenX)),
      screenY: Math.max(size, Math.min(dimensions.height - size, screenY)),
      color,
      size,
      chapter,
      type
    };
  };

  // Handle node click
  const handleNodeClick = (nodeId: string) => {
    setSelectedNode(selectedNode === nodeId ? null : nodeId);
  };

  if (!('Content' in data) || !data.Content || (Array.isArray(data.Content) && data.Content.length === 0)) {
    return (
      <div className={`flex items-center justify-center h-64 ${className}`}>
        <div className="text-center text-gray-400">
          <div className="w-16 h-16 mx-auto mb-4 opacity-50">
            <svg viewBox="0 0 24 24" fill="none" className="w-full h-full">
              <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="2"/>
              <circle cx="12" cy="12" r="3" fill="currentColor"/>
            </svg>
          </div>
          <p>No visualization data available</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`relative ${className}`}>
      {/* SVG Visualization */}
      <motion.svg
        ref={svgRef}
        width="100%"
        height="400"
        className="bg-gray-900/30 rounded-xl border border-gray-700/50"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5 }}
      >
        {/* Background grid */}
        <defs>
          <pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse">
            <path d="M 40 0 L 0 0 0 40" fill="none" stroke="rgba(148, 163, 184, 0.1)" strokeWidth="1"/>
          </pattern>
        </defs>
        <rect width="100%" height="100%" fill="url(#grid)" />

        {/* Connection lines */}
        <g className="connections">
          {nodes.map((node, i) => 
            nodes.slice(i + 1).map((otherNode) => {
              const distance = Math.sqrt(
                Math.pow(node.screenX - otherNode.screenX, 2) + 
                Math.pow(node.screenY - otherNode.screenY, 2)
              );
              
              // Only draw connections for nearby nodes
              if (distance > 150) return null;
              
              const opacity = Math.max(0.1, 1 - distance / 150);
              
              return (
                <motion.line
                  key={`${node.id}-${otherNode.id}`}
                  x1={node.screenX}
                  y1={node.screenY}
                  x2={otherNode.screenX}
                  y2={otherNode.screenY}
                  stroke="rgba(148, 163, 184, 0.3)"
                  strokeWidth={opacity * 2}
                  initial={{ pathLength: 0 }}
                  animate={{ pathLength: 1 }}
                  transition={{ duration: 1, delay: i * 0.1 }}
                />
              );
            })
          )}
        </g>

        {/* Nodes */}
        <g className="nodes">
          {nodes.map((node, index) => (
            <motion.g
              key={node.id}
              initial={{ scale: 0, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              transition={{ 
                duration: 0.5, 
                delay: index * 0.05,
                type: "spring",
                stiffness: 200 
              }}
              whileHover={{ scale: 1.2 }}
              whileTap={{ scale: 0.9 }}
              style={{ cursor: 'pointer' }}
              onClick={() => handleNodeClick(node.id)}
            >
              {/* Node glow effect */}
              <circle
                cx={node.screenX}
                cy={node.screenY}
                r={node.size + 5}
                fill={node.color}
                opacity={selectedNode === node.id ? 0.3 : 0.1}
                className="transition-opacity duration-300"
              />
              
              {/* Main node */}
              <circle
                cx={node.screenX}
                cy={node.screenY}
                r={node.size}
                fill={node.color}
                stroke={selectedNode === node.id ? "#ffffff" : "rgba(255, 255, 255, 0.2)"}
                strokeWidth={selectedNode === node.id ? 3 : 1}
                className="transition-all duration-300"
              />
              
              {/* Node label */}
              <text
                x={node.screenX}
                y={node.screenY + node.size + 20}
                textAnchor="middle"
                className="fill-white text-xs font-medium pointer-events-none"
                opacity={selectedNode === node.id ? 1 : 0.7}
              >
                {node.Text.length > 15 ? `${node.Text.substring(0, 15)}...` : node.Text}
              </text>
            </motion.g>
          ))}
        </g>
      </motion.svg>

      {/* Node Details Panel */}
      {selectedNode && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: 20 }}
          className="absolute bottom-4 left-4 right-4 glass rounded-lg p-4 border border-gray-700/50"
        >
          {(() => {
            const node = nodes.find(n => n.id === selectedNode);
            if (!node) return null;
            
            return (
              <div>
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center space-x-2">
                    <div 
                      className="w-3 h-3 rounded-full" 
                      style={{ backgroundColor: node.color }}
                    ></div>
                    <h4 className="text-white font-medium">{node.Text}</h4>
                    <span className="text-xs px-2 py-1 bg-gray-700 rounded-full text-gray-300">
                      {node.type}
                    </span>
                  </div>
                  <button
                    onClick={() => setSelectedNode(null)}
                    className="text-gray-400 hover:text-white"
                  >
                    ✕
                  </button>
                </div>
                
                {node.chapter && (
                  <div className="mb-3 text-sm">
                    <span className="text-gray-400">Chapter:</span>
                    <span className="text-primary-400 ml-1">{node.chapter}</span>
                  </div>
                )}
                
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                  <div>
                    <span className="text-gray-400">X:</span>
                    <span className="text-white ml-1">{node.XYZ.X.toFixed(3)}</span>
                  </div>
                  <div>
                    <span className="text-gray-400">Y:</span>
                    <span className="text-white ml-1">{node.XYZ.Y.toFixed(3)}</span>
                  </div>
                  <div>
                    <span className="text-gray-400">Z:</span>
                    <span className="text-white ml-1">{node.XYZ.Z.toFixed(3)}</span>
                  </div>
                  {node.XYZ.R !== undefined && (
                    <div>
                      <span className="text-gray-400">R:</span>
                      <span className="text-white ml-1">{node.XYZ.R.toFixed(3)}</span>
                    </div>
                  )}
                </div>
                
                {node.Reln && (
                  <div className="mt-3 pt-3 border-t border-gray-700">
                    <span className="text-gray-400 text-sm">Relations:</span>
                    <pre className="text-white text-xs mt-1 bg-gray-800 p-2 rounded overflow-x-auto">
                      {JSON.stringify(node.Reln, null, 2)}
                    </pre>
                  </div>
                )}
              </div>
            );
          })()}
        </motion.div>
      )}

      {/* Controls */}
      <div className="absolute top-4 right-4 flex space-x-2">
        <button
          onClick={() => setSelectedNode(null)}
          className="glass px-3 py-2 rounded-lg text-white text-sm hover:bg-gray-700/50 transition-colors"
        >
          Clear Selection
        </button>
        
        <div className="glass px-3 py-2 rounded-lg text-white text-sm">
          {nodes.length} nodes
        </div>
        
        <div className="glass px-3 py-2 rounded-lg text-white text-xs">
          {data.Response}
        </div>
      </div>
    </div>
  );
};

export default Visualization;