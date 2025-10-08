// Generic API Visualization v2.2 - Enhanced Arrows & Names
class GenericAPIVisualization
{
  constructor()
  {
    this.canvas3d = null;
    this.apiData = null;
    this.processedData = null;
    this.animationId = null;
    this.isAnimating = true;
    this.rotationSpeed = 0.01;
    this.currentAngle = 0;
    this.showContext = true;
    this.selectedItem = null;
    this.dataType = 'unknown';
    this.useSmartLabels = true; // Enable smart labels by default

    // Hover system
    this.hoveredNode = null;
    this.mouseX = 0;
    this.mouseY = 0;
    this.hoverTimeout = null;

    // Click/Focus system
    this.focusedNode = null;
    this.connectedNodes = new Set();
    this.connectedArrows = new Set();
    this.animationBeforeFocus = true;

    // Heat map system
    this.useHeatMap = false;
    this.heatMapCenter = null; // Node ID for the center of heat map
    this.maxDistance = 0; // Maximum distance for normalization
    this.heatMapColors = new Map(); // Cache for calculated colors

    this.init();
  }

  async init()
  {
    try
    {
      console.log("🚀 Initializing Generic API Visualization...");

      // Initialize the 3D canvas
      this.canvas3d = new SSTCanvas3D("canvasContainer", {
        width: window.innerWidth,
        height: window.innerHeight,
        mobile: window.innerWidth < 550 ? 0.5 : 1,
      });

      // Initialize hover system
      this.initHoverSystem();

      // Load initial data
      const data = document.getElementById("dataSource");
      await this.loadData(data.value);

      // Set up controls
      this.setupControls();

      // Start animation
      this.startAnimation();

      console.log("✅ Generic API visualization initialized successfully");
    } catch (error)
    {
      console.error("❌ Failed to initialize visualization:", error);
      this.showError("Failed to load visualization: " + error.message);
    }
  }

  // === SMART LABEL SYSTEM ===

  /**
   * Calculate distance from camera to a 3D point
   */
  calculateCameraDistance(x, y, z)
  {
    const camX = this.canvas3d.vpX;
    const camY = this.canvas3d.vpY;
    const camZ = this.canvas3d.vpZ;

    return Math.sqrt(
      Math.pow(x - camX, 2) +
      Math.pow(y - camY, 2) +
      Math.pow(z - camZ, 2)
    );
  }

  /**
   * Get node type significance priority (higher = more important)
   */
  getNodeSignificance(type)
  {
    const significance = {
      'event': 6,      // Highest priority
      'thing': 5,
      'concept': 4,
      'agent': 3,
      'expression': 2,
      'relation': 1    // Lowest priority
    };

    return significance[type] || 3; // Default medium significance
  }

  /**
   * Check if a label should be visible based on distance and significance
   */
  shouldShowLabel(x, y, z, type, labelType = 'main')
  {
    const distance = this.calculateCameraDistance(x, y, z);
    const significance = this.getNodeSignificance(type);

    // Base visibility thresholds
    const baseThresholds = {
      'main': 8,        // Node title labels
      'type': 6,        // Type indicator labels  
      'coords': 4,      // Coordinate labels
      'context': 5      // Context labels
    };

    const baseThreshold = baseThresholds[labelType] || 6;

    // Adjust threshold based on significance
    // More significant nodes show labels from further away
    const adjustedThreshold = baseThreshold + (significance - 3) * 1.5;

    return distance <= adjustedThreshold;
  }

  /**
   * Smart label drawing with distance and significance filtering
   */
  drawSmartLabel(x, y, z, text, size, color, type, labelType = 'main', nodeIndex = null)
  {
    // If smart labels are disabled, always show labels
    if (!this.useSmartLabels)
    {
      this.canvas3d.drawLabel(x, y, z, text, size, color);
      return true;
    }

    // Focus mode: only show labels for focused node and connected nodes
    if (this.focusedNode !== null && nodeIndex !== null)
    {
      if (nodeIndex === this.focusedNode || this.connectedNodes.has(nodeIndex))
      {
        // Show labels for focused node and connected nodes regardless of distance
        this.canvas3d.drawLabel(x, y, z, text, size, color);
        return true;
      }
      else
      {
        // Hide labels for non-connected nodes when in focus mode
        return false;
      }
    }

    // Normal mode: use distance and significance filtering
    if (this.shouldShowLabel(x, y, z, type, labelType))
    {
      this.canvas3d.drawLabel(x, y, z, text, size, color);
      return true;
    }
    return false;
  }

  // === HOVER SYSTEM ===

  /**
   * Initialize hover system with mouse tracking
   */
  initHoverSystem()
  {
    if (!this.canvas3d.canvas) return;

    const canvas = this.canvas3d.canvas;

    // Create hover tooltip
    this.createHoverTooltip();

    // Mouse move handler
    canvas.addEventListener('mousemove', (e) =>
    {
      const rect = canvas.getBoundingClientRect();
      this.mouseX = e.clientX - rect.left;
      this.mouseY = e.clientY - rect.top;

      // Clear existing hover timeout
      if (this.hoverTimeout)
      {
        clearTimeout(this.hoverTimeout);
      }

      // Check for hovered node after short delay to avoid excessive calculations
      this.hoverTimeout = setTimeout(() =>
      {
        this.checkNodeHover();
      }, 50);
    });

    // Mouse leave handler
    canvas.addEventListener('mouseleave', () =>
    {
      this.hideHoverTooltip();
      this.hoveredNode = null;
    });

    // Click handler for node focusing
    canvas.addEventListener('click', (e) =>
    {
      const rect = canvas.getBoundingClientRect();
      const clickX = e.clientX - rect.left;
      const clickY = e.clientY - rect.top;

      this.handleNodeClick(clickX, clickY);
    });
  }

  /**
   * Create hover tooltip element
   */
  createHoverTooltip()
  {
    // Remove existing tooltip
    const existing = document.getElementById('hover-tooltip');
    if (existing) existing.remove();

    const tooltip = document.createElement('div');
    tooltip.id = 'hover-tooltip';
    tooltip.style.cssText = `
      position: absolute;
      background: rgba(0, 0, 0, 0.9);
      color: white;
      padding: 8px 12px;
      border-radius: 6px;
      font-size: 12px;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      pointer-events: none;
      z-index: 1000;
      max-width: 250px;
      word-wrap: break-word;
      box-shadow: 0 2px 8px rgba(0,0,0,0.3);
      border: 1px solid #444;
      display: none;
    `;

    document.body.appendChild(tooltip);
    this.tooltip = tooltip;
  }

  /**
   * Check if mouse is hovering over any node
   */
  checkNodeHover()
  {
    if (!this.processedData || !this.processedData.items) return;

    let closestNode = null;
    let closestDistance = Infinity;
    let debugCount = 0;

    // Check all nodes
    for (let i = 0; i < this.processedData.items.length; i++)
    {
      const item = this.processedData.items[i];
      if (!item.XYZ) 
      {
        if (debugCount < 3) console.log(`⚠️ Node ${i} missing XYZ:`, item);
        debugCount++;
        continue;
      }

      const screenPos = this.canvas3d.project3D(item.XYZ.X, item.XYZ.Y, item.XYZ.Z);
      if (!screenPos) 
      {
        if (debugCount < 3) console.log(`⚠️ Node ${i} failed projection:`, item.XYZ);
        debugCount++;
        continue;
      }

      const distance = Math.sqrt(
        Math.pow(this.mouseX - screenPos.x, 2) +
        Math.pow(this.mouseY - screenPos.y, 2)
      );

      // Check if within node radius (base on node type significance)
      const nodeRadius = this.getNodeHoverRadius(item.type);

      // Debug nearby nodes
      if (distance < nodeRadius * 2 && debugCount < 5)
      {
        console.log(`🎯 Node ${i} (${item.type || 'unknown'}) at distance ${distance.toFixed(1)} (threshold: ${nodeRadius})`);
        debugCount++;
      }

      if (distance <= nodeRadius && distance < closestDistance)
      {
        closestDistance = distance;
        closestNode = { item, index: i, distance };
        console.log(`✅ Found hover node ${i} at distance ${distance.toFixed(1)}`);
      }
    }

    if (closestNode && closestNode !== this.hoveredNode)
    {
      this.hoveredNode = closestNode;
      this.showHoverTooltip(closestNode.item, closestNode.index);
    } else if (!closestNode && this.hoveredNode)
    {
      this.hideHoverTooltip();
      this.hoveredNode = null;
    }
  }

  /**
   * Get hover detection radius based on node type
   */
  getNodeHoverRadius(type)
  {
    const baseRadius = {
      'event': 25,
      'thing': 20,
      'concept': 20,
      'agent': 15,
      'expression': 12,
      'relation': 10,
      // TOC-specific node types
      'chapter': 25,      // Main chapters - large radius like events
      'context': 20,      // Context items - medium radius like concepts
      'fragment': 15      // Fragment items - smaller radius
    };

    return baseRadius[type] || 15;
  }

  /**
   * Show hover tooltip with node information
   */
  showHoverTooltip(item, index)
  {
    if (!this.tooltip) return;

    const title = item.title || `Node ${index}`;
    const type = item.type || 'unknown';
    const coords = `(${item.XYZ.X.toFixed(2)}, ${item.XYZ.Y.toFixed(2)}, ${item.XYZ.Z.toFixed(2)})`;

    let content = `<strong>${title}</strong><br>`;
    content += `<span style="color: #aaa;">Type:</span> ${type}<br>`;
    content += `<span style="color: #aaa;">Position:</span> ${coords}`;

    if (item.context)
    {
      content += `<br><span style="color: #aaa;">Context:</span> ${item.context}`;
    }

    this.tooltip.innerHTML = content;
    this.tooltip.style.display = 'block';
    this.tooltip.style.left = (this.mouseX + 15) + 'px';
    this.tooltip.style.top = (this.mouseY - 10) + 'px';

    // Adjust position if tooltip goes off screen
    const rect = this.tooltip.getBoundingClientRect();
    if (rect.right > window.innerWidth)
    {
      this.tooltip.style.left = (this.mouseX - rect.width - 15) + 'px';
    }
    if (rect.bottom > window.innerHeight)
    {
      this.tooltip.style.top = (this.mouseY - rect.height - 10) + 'px';
    }
  }

  /**
   * Hide hover tooltip
   */
  hideHoverTooltip()
  {
    if (this.tooltip)
    {
      this.tooltip.style.display = 'none';
    }
  }

  async loadData(filename)
  {
    try
    {
      console.log(`📡 Loading data from ${filename}...`);
      const response = await fetch(filename);
      if (!response.ok)
      {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      this.apiData = await response.json();

      // Detect and process data type
      this.dataType = this.detectDataType(this.apiData);
      console.log(`🔍 Detected data type: ${this.dataType}`);

      this.processedData = this.processAPIData(this.apiData);

      this.updateInfoPanel();
      this.updateItemList();
      this.updateHeatMapOptions(); // Update heat map center options
      this.renderVisualization();

      console.log("✅ API data loaded and processed:", this.processedData);
    } catch (error)
    {
      console.error("❌ Error loading data:", error);
      throw error;
    }
  }

  // === CLICK/FOCUS SYSTEM ===

  /**
   * Handle mouse click to focus on a node
   */
  handleNodeClick(clickX, clickY)
  {
    if (!this.processedData || !this.processedData.items) return;

    let clickedNode = null;
    let closestDistance = Infinity;
    let debugCount = 0;

    console.log(`🖱️ Click at (${clickX.toFixed(1)}, ${clickY.toFixed(1)})`);

    // Check all nodes for click collision
    for (let i = 0; i < this.processedData.items.length; i++)
    {
      const item = this.processedData.items[i];
      if (!item.XYZ) 
      {
        if (debugCount < 3) console.log(`⚠️ Click check: Node ${i} missing XYZ`);
        debugCount++;
        continue;
      }

      const screenPos = this.canvas3d.project3D(item.XYZ.X, item.XYZ.Y, item.XYZ.Z);
      if (!screenPos) 
      {
        if (debugCount < 3) console.log(`⚠️ Click check: Node ${i} failed projection`);
        debugCount++;
        continue;
      }

      const distance = Math.sqrt(
        Math.pow(clickX - screenPos.x, 2) +
        Math.pow(clickY - screenPos.y, 2)
      );

      const nodeRadius = this.getNodeHoverRadius(item.type);

      // Debug nearby nodes
      if (distance < nodeRadius * 2 && debugCount < 5)
      {
        console.log(`🎯 Click check: Node ${i} (${item.type || 'unknown'}) at distance ${distance.toFixed(1)} (threshold: ${nodeRadius})`);
        debugCount++;
      }

      if (distance <= nodeRadius && distance < closestDistance)
      {
        closestDistance = distance;
        clickedNode = { item, index: i, distance };
        console.log(`✅ Found click node ${i} at distance ${distance.toFixed(1)}`);
      }
    }

    if (clickedNode)
    {
      this.focusOnNode(clickedNode.index);
    }
    else
    {
      console.log(`❌ No node found at click position`);
      // Click on empty space - clear focus
      this.clearFocus();
    }
  }

  /**
   * Focus on a specific node by centering it and highlighting connections
   */
  focusOnNode(nodeIndex)
  {
    if (!this.processedData || !this.processedData.items[nodeIndex]) return;

    console.log(`🎯 Focusing on node ${nodeIndex}`);

    // Store animation state and stop animation
    this.animationBeforeFocus = this.isAnimating;
    if (this.isAnimating)
    {
      this.stopAnimation();
      document.getElementById("toggleAnimation").textContent = "▶️ Play";
    }

    // Set focused node
    this.focusedNode = nodeIndex;
    this.selectedItem = nodeIndex; // Also set as selected for highlighting

    // Find and highlight connected nodes and arrows
    this.highlightConnections(nodeIndex);

    // Activate heat map with focused node as center (after connections are established)
    this.activateHeatMapForFocus(nodeIndex);

    // Enhanced focus: calculate optimal view for connected nodes
    this.calculateOptimalFocusView(nodeIndex);

    // Force re-render
    this.renderVisualization();

    // Update UI to show focus state
    this.updateFocusInfo(nodeIndex);
  }

  /**
   * Calculate optimal camera position and zoom to frame the focused node and all connections
   */
  calculateOptimalFocusView(nodeIndex)
  {
    const focusedNode = this.processedData.items[nodeIndex];
    if (!focusedNode || !focusedNode.XYZ) return;

    // Start with the focused node coordinates
    let minX = focusedNode.XYZ.X;
    let maxX = focusedNode.XYZ.X;
    let minY = focusedNode.XYZ.Y;
    let maxY = focusedNode.XYZ.Y;
    let minZ = focusedNode.XYZ.Z;
    let maxZ = focusedNode.XYZ.Z;

    // Include all connected nodes in the bounding calculation
    this.connectedNodes.forEach(connectedIndex =>
    {
      const connectedNode = this.processedData.items[connectedIndex];
      if (connectedNode && connectedNode.XYZ)
      {
        minX = Math.min(minX, connectedNode.XYZ.X);
        maxX = Math.max(maxX, connectedNode.XYZ.X);
        minY = Math.min(minY, connectedNode.XYZ.Y);
        maxY = Math.max(maxY, connectedNode.XYZ.Y);
        minZ = Math.min(minZ, connectedNode.XYZ.Z);
        maxZ = Math.max(maxZ, connectedNode.XYZ.Z);
      }
    });

    // Calculate the center point of all connected nodes
    const centerX = (minX + maxX) / 2;
    const centerY = (minY + maxY) / 2;
    const centerZ = (minZ + maxZ) / 2;

    // Calculate the span of connected nodes
    const spanX = maxX - minX;
    const spanY = maxY - minY;
    const spanZ = maxZ - minZ;
    const maxSpan = Math.max(spanX, spanY, spanZ);

    // Calculate optimal zoom distance
    // For focus mode, we want to zoom IN (closer) to the focused node
    let optimalDistance;
    if (this.connectedNodes.size > 0)
    {
      // Even with connections, zoom in closer than the span to focus on the center
      optimalDistance = Math.max(maxSpan * 0.8, 1.5); // Much closer, minimum distance of 1.5
    } else
    {
      // Single node focus - zoom in very close
      optimalDistance = 1.0;
    }

    // For proper centering and zoom-in: adjust the viewpoint to center on the focused node
    // Keep objects in their original positions, move the viewpoint instead

    const focusX = focusedNode.XYZ.X;
    const focusY = focusedNode.XYZ.Y;
    const focusZ = focusedNode.XYZ.Z;

    // Store original origin and scale to restore later
    if (this.originalScale === undefined)
    {
      this.originalScale = this.canvas3d.scale;
      this.originalOrgX = this.canvas3d.orgX;
      this.originalOrgY = this.canvas3d.orgY;
    }

    // Calculate where the focused node appears on screen with current settings
    const currentScreenX = this.canvas3d.transformX(focusX, focusY, focusZ);
    const currentScreenY = this.canvas3d.transformY(focusX, focusY, focusZ);

    // Calculate offset needed to center the focused node
    const screenCenterX = this.canvas3d.width / 2;
    const screenCenterY = this.canvas3d.height / 2;
    const offsetX = screenCenterX - currentScreenX;
    const offsetY = screenCenterY - currentScreenY;

    // Adjust the screen origin to center the focused node
    this.canvas3d.orgX = this.originalOrgX + offsetX;
    this.canvas3d.orgY = this.originalOrgY + offsetY;

    // Zoom in by increasing the scale factor  
    if (this.connectedNodes.size > 0)
    {
      // Medium zoom for connected nodes
      this.canvas3d.scale = this.originalScale * 2.5;
    } else
    {
      // Close zoom for isolated nodes  
      this.canvas3d.scale = this.originalScale * 4.0;
    }

    console.log(`📷 Enhanced focus: Focused node(${focusX.toFixed(2)}, ${focusY.toFixed(2)}, ${focusZ.toFixed(2)}), Scale: ${this.canvas3d.scale.toFixed(2)}, Connected nodes: ${this.connectedNodes.size}`);
    console.log(`🎥 Viewpoint adjusted: orgX offset ${offsetX.toFixed(1)}, orgY offset ${offsetY.toFixed(1)}`);
  }

  /**
   * Find and highlight all nodes and arrows connected to the focused node
   */
  highlightConnections(nodeIndex)
  {
    this.connectedNodes.clear();
    this.connectedArrows.clear();

    if (!this.processedData || !this.processedData.arrows) return;

    // Find all arrows connected to this node
    for (let i = 0; i < this.processedData.arrows.length; i++)
    {
      const arrow = this.processedData.arrows[i];

      if (arrow.from === nodeIndex)
      {
        // Outgoing arrow
        this.connectedArrows.add(i);
        this.connectedNodes.add(arrow.to);
      }
      else if (arrow.to === nodeIndex)
      {
        // Incoming arrow
        this.connectedArrows.add(i);
        this.connectedNodes.add(arrow.from);
      }
    }

    console.log(`🔗 Found ${this.connectedArrows.size} connected arrows and ${this.connectedNodes.size} connected nodes`);
  }

  /**
   * Clear focus and restore normal view
   */
  clearFocus()
  {
    console.log("🔄 Clearing focus");

    this.focusedNode = null;
    this.selectedItem = null;
    this.connectedNodes.clear();
    this.connectedArrows.clear();

    // Restore original scale and screen origin
    if (this.originalScale !== undefined)
    {
      this.canvas3d.scale = this.originalScale;
    }
    if (this.originalOrgX !== undefined)
    {
      this.canvas3d.orgX = this.originalOrgX;
    }
    if (this.originalOrgY !== undefined)
    {
      this.canvas3d.orgY = this.originalOrgY;
    }

    // Restore animation if it was running before
    if (this.animationBeforeFocus && !this.isAnimating)
    {
      this.startAnimation();
      document.getElementById("toggleAnimation").textContent = "⏸️ Pause";
    }

    // Reset camera to default position
    this.resetCameraPosition();

    // Force re-render
    this.renderVisualization();

    // Clear focus info
    this.clearFocusInfo();
  }

  /**
   * Reset camera to default viewing position
   */
  resetCameraPosition()
  {
    // Reset viewport
    this.canvas3d.vpX = 0;
    this.canvas3d.vpY = 4;
    this.canvas3d.vpZ = 8;

    // Reset observer position (this is what actually controls the view)
    this.canvas3d.setObserverPosition(1, 0.5, -1);

    console.log("📷 Camera reset to default position");
  }

  /**
   * Update UI to show information about the focused node
   */
  updateFocusInfo(nodeIndex)
  {
    const node = this.processedData.items[nodeIndex];
    if (!node) return;

    // Create or update focus info panel
    let focusPanel = document.getElementById('focus-info');
    if (!focusPanel)
    {
      focusPanel = document.createElement('div');
      focusPanel.id = 'focus-info';
      focusPanel.style.cssText = `
        position: fixed;
        top: 10px;
        left: 10px;
        background: rgba(0, 0, 0, 0.9);
        color: white;
        padding: 15px;
        border-radius: 8px;
        max-width: 300px;
        z-index: 1001;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        border: 2px solid #FFD700;
        box-shadow: 0 4px 12px rgba(0,0,0,0.3);
      `;
      document.body.appendChild(focusPanel);
    }

    const title = node.title || `Node ${nodeIndex}`;
    const connections = this.connectedNodes.size;

    focusPanel.innerHTML = `
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
        <strong style="color: #FFD700;">🎯 ENHANCED FOCUS MODE</strong>
        <button onclick="window.visualization.clearFocus()" style="background: #666; border: none; color: white; border-radius: 4px; padding: 4px 8px; cursor: pointer;">✕</button>
      </div>
      <div><strong>Focused Node:</strong> ${title}</div>
      <div><strong>Type:</strong> ${node.type || 'unknown'}</div>
      <div><strong>Index:</strong> ${nodeIndex}</div>
      <div><strong>Connected Nodes:</strong> ${connections}</div>
      <div><strong>Connected Arrows:</strong> ${this.connectedArrows.size}</div>
      <div style="margin-top: 8px; padding: 6px; background: rgba(99, 102, 241, 0.1); border: 1px solid rgba(99, 102, 241, 0.3); border-radius: 4px; font-size: 11px;">
        🌡️ <strong>Heat Map Active:</strong><br>
        • Connected nodes colored by distance<br>
        • Red = close to focus, Blue = far<br>
        • Only connected cluster shown
      </div>
      <div style="margin-top: 8px; padding: 6px; background: rgba(255, 215, 0, 0.1); border-radius: 4px; font-size: 11px;">
        📱 <strong>Enhanced Focus:</strong><br>
        • Only connected nodes visible<br>
        • Optimal zoom and centering<br>
        • Click elsewhere to exit
      </div>
    `;
  }

  /**
   * Clear focus information panel
   */
  clearFocusInfo()
  {
    const focusPanel = document.getElementById('focus-info');
    if (focusPanel)
    {
      focusPanel.remove();
    }
  }

  detectDataType(data)
  {
    console.log("🔍 Detecting data type for:", data);

    if (!data || !data.Response)
    {
      console.log("⚠️ No data or Response field found");
      return 'unknown';
    }

    console.log(`📊 Response type: ${data.Response}`);

    switch (data.Response)
    {
    case 'TOC':
      return 'toc';
    case 'PageMap':
    case 'Orbits':
      return 'orbits';
    case 'ConePaths':
    case 'PathSolve':
      return 'paths';
    case 'Sequence':
      return 'sequence';
    case 'Stories':
      return 'stories';
    case 'ChapterContexts':
      return 'contexts';
    case 'Stats':
      return 'stats';
    case 'Arrows':
      return 'arrows';
    default:
      console.log(`❓ Unknown response type: ${data.Response}`);
      return 'generic';
    }
  }

  processAPIData(data)
  {
    console.log(`🔄 Processing ${this.dataType} data...`);

    switch (this.dataType)
    {
    case 'toc':
      return this.processTOCData(data);
    case 'orbits':
      return this.processOrbitData(data);
    case 'paths':
      return this.processPathData(data);
    case 'sequence':
      return this.processSequenceData(data);
    case 'stories':
      return this.processStoryData(data);
    case 'contexts':
      return this.processContextData(data);
    case 'arrows':
      return this.processArrowData(data);
    default:
      return this.processGenericData(data);
    }
  }

  processTOCData(data)
  {
    // Enhanced TOC processing with arrows
    const items = [];
    const arrows = [];

    if (data.Content)
    {
      data.Content.forEach((chapter, chapterIndex) =>
      {
        // Add main chapter
        items.push({
          ...chapter,
          type: 'chapter',
          title: chapter.Chapter || `Chapter ${chapterIndex + 1}`,
          dataType: 'chapter'
        });

        const chapterItemIndex = items.length - 1;

        // Add context items and create arrows
        if (chapter.Context)
        {
          chapter.Context.forEach((context) =>
          {
            if (context.XYZ)
            {
              items.push({
                ...context,
                type: 'context',
                title: context.Text || 'Context',
                dataType: 'context'
              });

              // Create arrow from chapter to context
              arrows.push({
                from: chapterItemIndex,
                to: items.length - 1,
                type: 300, // Contains relationship (lightblue)
                weight: 1
              });
            }
          });
        }

        // Add single items and create arrows
        if (chapter.Single)
        {
          chapter.Single.forEach((single) =>
          {
            if (single.XYZ)
            {
              items.push({
                ...single,
                type: 'fragment',
                title: single.Text || 'Fragment',
                dataType: 'fragment'
              });

              // Create arrow from chapter to fragment
              arrows.push({
                from: chapterItemIndex,
                to: items.length - 1,
                type: 302, // Expresses relationship (orange)
                weight: 1
              });
            }
          });
        }
      });
    }

    return {
      dataType: 'toc',
      title: 'Table of Contents',
      items: items,
      arrows: arrows,
      metadata: {
        response: data.Response,
        time: data.Time,
        intent: data.Intent,
        ambient: data.Ambient
      }
    };
  }

  processOrbitData(data)
  {
    console.log("🌍 Processing orbit data:", data);

    // Process PageMap/Orbit data structure
    const items = [];
    const arrows = [];

    if (data.Content && data.Content.Notes)
    {
      console.log(`📝 Found ${data.Content.Notes.length} note groups`);

      // First pass: create all items with their indices
      const itemIndexMap = new Map(); // Map from NPtr to array index

      data.Content.Notes.forEach((noteGroup, groupIndex) =>
      {
        console.log(`📋 Processing group ${groupIndex} with ${noteGroup.length} notes`);

        noteGroup.forEach((note, noteIndex) =>
        {
          if (note.XYZ)
          {
            const item = {
              id: `${note.NPtr?.Class || 0}-${note.NPtr?.CPtr || 0}`,
              title: note.Name || `Node ${groupIndex}-${noteIndex}`,
              chapter: note.Chp || '',
              context: note.Ctx || '',
              XYZ: note.XYZ,
              type: this.getNodeType(note),
              group: groupIndex,
              note: note,
              nptrClass: note.NPtr?.Class,
              nptrCPtr: note.NPtr?.CPtr
            };

            // Store mapping for arrow references
            if (note.NPtr)
            {
              itemIndexMap.set(`${note.NPtr.Class}-${note.NPtr.CPtr}`, items.length);
            }

            items.push(item);
            console.log(`✅ Added item: ${item.title} at (${item.XYZ.X}, ${item.XYZ.Y}, ${item.XYZ.Z})`);
          } else
          {
            console.log(`⚠️ Note ${noteIndex} in group ${groupIndex} missing XYZ coordinates`);
          }
        });
      });

      // Second pass: create arrows based on semantic relationships
      items.forEach((item, fromIndex) =>
      {
        if (item.note.Arr > 0)
        {
          // Create arrows based on actual arrow type from data
          const arrowType = item.note.Arr;

          // Find nearby items or items of compatible types for relationships
          items.forEach((targetItem, toIndex) =>
          {
            if (fromIndex !== toIndex && this.shouldCreateArrow(item, targetItem, arrowType))
            {
              arrows.push({
                from: fromIndex,
                to: toIndex,
                type: arrowType,
                weight: 1
              });
              console.log(`🔗 Added semantic arrow from ${fromIndex} to ${toIndex} (type: ${arrowType})`);
            }
          });
        }
      });

      // Third pass: Add more arrow variety for better visualization
      items.forEach((item, fromIndex) =>
      {
        items.forEach((targetItem, toIndex) =>
        {
          if (fromIndex !== toIndex)
          {
            const distance = this.calculateDistance(item.XYZ, targetItem.XYZ);

            // Create arrows based on proximity and types
            if (distance < 0.15)
            { // Very close items
              const arrowType = this.getRandomArrowType(item, targetItem);

              // Avoid duplicate arrows
              const existingArrow = arrows.find(a =>
                (a.from === fromIndex && a.to === toIndex) ||
                (a.from === toIndex && a.to === fromIndex)
              );

              if (!existingArrow && Math.random() < 0.4)
              { // 40% chance
                arrows.push({
                  from: fromIndex,
                  to: toIndex,
                  type: arrowType,
                  weight: 1
                });
                console.log(`🔗 Added proximity arrow from ${fromIndex} to ${toIndex} (type: ${arrowType})`);
              }
            }
          }
        });
      });

    } else
    {
      console.log("❌ No Content.Notes found in data");
    }

    const result = {
      dataType: 'orbits',
      title: data.Content?.Title || 'Orbit Visualization',
      context: data.Content?.Context || '',
      items: items,
      arrows: arrows,
      metadata: {
        response: data.Response,
        time: data.Time,
        intent: data.Intent,
        ambient: data.Ambient
      }
    };

    console.log(`✅ Processed orbit data: ${items.length} items, ${arrows.length} arrows`);
    return result;
  }

  shouldCreateArrow(fromItem, toItem, arrowType)
  {
    // Create arrows based on semantic relationships and proximity
    const distance = this.calculateDistance(fromItem.XYZ, toItem.XYZ);
    const maxDistance = 0.3; // Only connect nearby items

    // Type-based relationships with better logic
    if (arrowType === 298)
    { // "has value" relationship - LeadsTo
      return (fromItem.type === 'concept' && (toItem.type === 'thing' || toItem.type === 'expression')) && distance < maxDistance;
    } else if (arrowType === 276)
    { // Expresses relationship
      return (fromItem.type === 'expression' || fromItem.type === 'concept') && distance < maxDistance;
    } else if (arrowType === 258)
    { // General relationship
      return fromItem.type !== toItem.type && distance < maxDistance;
    }

    return distance < maxDistance; // Default proximity-based connection
  }

  calculateDistance(xyz1, xyz2)
  {
    const dx = xyz1.X - xyz2.X;
    const dy = xyz1.Y - xyz2.Y;
    const dz = xyz1.Z - xyz2.Z;
    return Math.sqrt(dx * dx + dy * dy + dz * dz);
  }

  getRandomArrowType(fromItem, toItem)
  {
    // Create variety in arrow types based on item characteristics
    const distance = this.calculateDistance(fromItem.XYZ, toItem.XYZ);

    // Different logic for different type combinations
    if (fromItem.type === 'concept' && toItem.type === 'expression')
    {
      return 302; // expresses (orange)
    } else if (fromItem.type === 'concept' && (toItem.type === 'thing' || toItem.type === 'agent'))
    {
      return 300; // contains (lightblue)
    } else if (fromItem.type === 'expression' && toItem.type === 'concept')
    {
      return 298; // leadsTo (darkred)
    } else if (distance < 0.1)
    {
      return 258; // near (darkgrey) for very close items
    } else if (fromItem.XYZ.Z > toItem.XYZ.Z)
    {
      return 298; // leadsTo (sequence from higher to lower Z)
    } else if (Math.abs(fromItem.XYZ.Y - toItem.XYZ.Y) < 0.05)
    {
      return 300; // contains (same level)
    } else
    {
      // Random selection for more variety
      const types = [298, 300, 302, 258];
      return types[Math.floor(Math.random() * types.length)];
    }
  }

  getNodeType(note)
  {
    if (!note.NPtr) return 'unknown';

    switch (note.NPtr.Class)
    {
    case 0: return 'relation';
    case 1: return 'event';
    case 2: return 'agent';
    case 3: return 'thing';
    case 4: return 'concept';
    case 5: return 'expression';
    default: return 'unknown';
    }
  }

  processPathData(data)
  {
    // Handle ConePaths/PathSolve data
    return {
      type: 'paths',
      title: 'Path Visualization',
      items: this.extractNodesFromContent(data.Content),
      metadata: {
        response: data.Response,
        time: data.Time,
        intent: data.Intent,
        ambient: data.Ambient
      }
    };
  }

  processSequenceData(data)
  {
    return {
      type: 'sequence',
      title: 'Sequence Visualization',
      items: this.extractNodesFromContent(data.Content),
      metadata: {
        response: data.Response,
        time: data.Time,
        intent: data.Intent,
        ambient: data.Ambient
      }
    };
  }

  processStoryData(data)
  {
    return {
      type: 'stories',
      title: 'Story Visualization',
      items: this.extractNodesFromContent(data.Content),
      metadata: {
        response: data.Response,
        time: data.Time,
        intent: data.Intent,
        ambient: data.Ambient
      }
    };
  }

  processContextData(data)
  {
    return {
      type: 'contexts',
      title: 'Context Visualization',
      items: this.extractNodesFromContent(data.Content),
      metadata: {
        response: data.Response,
        time: data.Time,
        intent: data.Intent,
        ambient: data.Ambient
      }
    };
  }

  processArrowData(data)
  {
    return {
      type: 'arrows',
      title: 'Arrow Relationships',
      items: this.extractNodesFromContent(data.Content),
      metadata: {
        response: data.Response,
        time: data.Time,
        intent: data.Intent,
        ambient: data.Ambient
      }
    };
  }

  processGenericData(data)
  {
    return {
      type: 'generic',
      title: 'Generic Visualization',
      items: this.extractNodesFromContent(data.Content),
      metadata: {
        response: data.Response,
        time: data.Time,
        intent: data.Intent,
        ambient: data.Ambient
      }
    };
  }

  extractNodesFromContent(content)
  {
    // Generic node extraction for unknown data formats
    const items = [];

    if (Array.isArray(content))
    {
      content.forEach((item, index) =>
      {
        if (item && item.XYZ)
        {
          items.push({
            id: index,
            title: item.Name || item.Text || item.title || `Item ${index}`,
            XYZ: item.XYZ,
            type: 'generic',
            raw: item
          });
        }
      });
    } else if (content && typeof content === 'object')
    {
      // Try to find nested arrays with XYZ data
      Object.keys(content).forEach(key =>
      {
        if (Array.isArray(content[key]))
        {
          content[key].forEach((item, index) =>
          {
            if (item && item.XYZ)
            {
              items.push({
                id: `${key}-${index}`,
                title: item.Name || item.Text || `${key} ${index}`,
                XYZ: item.XYZ,
                type: 'generic',
                group: key,
                raw: item
              });
            }
          });
        }
      });
    }

    return items;
  }

  renderVisualization()
  {
    if (!this.processedData || !this.processedData.items)
    {
      console.warn("⚠️ No processed data to render");
      return;
    }

    // Clear canvas
    this.canvas3d.clear();

    // Draw coordinate grid
    this.canvas3d.drawGrid(0, 0, 1);

    // Render items based on data type
    switch (this.dataType)
    {
    case 'toc':
      this.renderTOCVisualization();
      break;
    case 'orbits':
      this.renderOrbitVisualization();
      break;
    default:
      this.renderGenericVisualization();
      break;
    }

    // Draw connections/arrows if available
    if (this.processedData.arrows)
    {
      this.drawArrows();
    }
  }

  renderTOCVisualization()
  {
    // Original TOC rendering logic with focus mode filtering
    this.processedData.items.forEach((chapter, index) =>
    {
      // In focus mode, only render focused node and connected nodes
      if (this.focusedNode !== null)
      {
        if (index !== this.focusedNode && !this.connectedNodes.has(index))
        {
          return; // Skip unconnected nodes
        }
      }

      this.renderTOCChapter(chapter, index);
    });
    this.drawTOCConnections();
  }

  renderOrbitVisualization()
  {
    // New orbit rendering logic with focus mode filtering
    console.log(`🌍 Rendering ${this.processedData.items.length} orbit nodes...`);

    this.processedData.items.forEach((item, index) =>
    {
      // In focus mode, only render focused node and connected nodes
      if (this.focusedNode !== null)
      {
        if (index !== this.focusedNode && !this.connectedNodes.has(index))
        {
          return; // Skip unconnected nodes
        }
      }

      this.renderOrbitNode(item, index);
    });

    // Draw arrows between nodes with focus filtering
    if (this.processedData.arrows)
    {
      this.processedData.arrows.forEach((arrow, arrowIndex) =>
      {
        // In focus mode, only render connected arrows
        if (this.focusedNode !== null)
        {
          if (!this.connectedArrows.has(arrowIndex))
          {
            return; // Skip unconnected arrows
          }
        }

        this.renderArrow(arrow);
      });
    }
  }

  renderGenericVisualization()
  {
    // Generic rendering for unknown data types with focus mode filtering
    this.processedData.items.forEach((item, index) =>
    {
      // In focus mode, only render focused node and connected nodes
      if (this.focusedNode !== null)
      {
        if (index !== this.focusedNode && !this.connectedNodes.has(index))
        {
          return; // Skip unconnected nodes
        }
      }

      this.renderGenericNode(item, index);
    });
  }

  renderTOCChapter(chapter, index)
  {
    const x = chapter.XYZ.X;
    const y = chapter.XYZ.Y;
    const z = chapter.XYZ.Z;

    // Get heat map color if enabled
    const itemId = chapter.id || index.toString();
    const heatMapColor = this.getNodeHeatMapColor(itemId);

    // Use specialized rendering based on item type
    if (chapter.type === 'chapter')
    {
      // Main chapters - use event rendering (large red nodes)
      this.canvas3d.drawEvent(x, y, z, heatMapColor);
    } else if (chapter.type === 'context')
    {
      // Context items - use concept rendering (blue nodes)
      this.canvas3d.drawConcept(x, y, z, heatMapColor);
    } else if (chapter.type === 'fragment')
    {
      // Fragment items - use thing rendering (green nodes)
      this.canvas3d.drawThing(x, y, z, heatMapColor);
    } else
    {
      // Fallback for original chapter structure
      this.canvas3d.drawEvent(x, y, z, heatMapColor);
    }

    // Draw title with type-appropriate handling
    let title;
    if (chapter.type === 'chapter')
    {
      title = chapter.title || chapter.Chapter || `Chapter ${index + 1}`;
    } else
    {
      title = chapter.title || chapter.Text || `Item ${index + 1}`;
    }

    // Use smart label system for chapter titles
    this.drawSmartLabel(x, y, z, title.slice(0, 40), 12, "white", chapter.type || 'chapter', 'main', index);

    // Draw type indicator - only if close
    if (chapter.type)
    {
      const colors = this.getNodeColors(chapter.type);
      this.drawSmartLabel(x, y - 0.05, z, `[${chapter.type}]`, 6, colors.fill, chapter.type, 'type', index);
    }

    // Draw coordinates info - only if very close
    this.drawSmartLabel(x, y - 0.08, z, `(${x.toFixed(2)}, ${y.toFixed(2)}, ${z.toFixed(2)})`, 5, "oklch(85.5% 0.138 181.071)", chapter.type || 'chapter', 'coords', index);

    // Render context fragments if enabled (for original structure)
    if (this.showContext && chapter.Context)
    {
      this.renderTOCContextFragments(chapter, x, y, z, index);
    }

    // Highlight selected item
    if (this.selectedItem === index)
    {
      this.canvas3d.drawNode(x, y, z, 7 * this.canvas3d.mob, "transparent", "#FFD700");
    }

    // Highlight hovered item
    if (this.hoveredNode && this.hoveredNode.index === index)
    {
      this.canvas3d.drawNode(x, y, z, 6 * this.canvas3d.mob, "transparent", "#FFFFFF40");
    }

    // Highlight connected nodes when focusing
    if (this.focusedNode !== null)
    {
      if (this.connectedNodes.has(index))
      {
        // Connected node - highlight with green glow
        this.canvas3d.drawNode(x, y, z, 8 * this.canvas3d.mob, "transparent", "#00FF0060");
      }
    }
  }

  renderOrbitNode(item, index)
  {
    const x = item.XYZ.X;
    const y = item.XYZ.Y;
    const z = item.XYZ.Z;

    let nodeSize = 4; // Default size

    // Get heat map color if enabled
    const itemId = item.id || index.toString();
    const heatMapColor = this.getNodeHeatMapColor(itemId);

    // Use specialized node drawing methods from SSTCanvas3D when available
    switch (item.type)
    {
    case 'event':
      this.canvas3d.drawEvent(x, y, z, heatMapColor);
      nodeSize = 6; // Events are larger
      break;
    case 'thing':
      this.canvas3d.drawThing(x, y, z, heatMapColor);
      nodeSize = 4; // Things are medium
      break;
    case 'concept':
      this.canvas3d.drawConcept(x, y, z, heatMapColor);
      nodeSize = 4; // Concepts are medium
      break;
    case 'agent':
      // Agent nodes - use orange color similar to 'thing' but different
      nodeSize = 3;
      if (heatMapColor)
      {
        this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, heatMapColor.inner, heatMapColor.outer);
      } else
      {
        this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, "#f39c12", "#e67e22");
      }
      break;
    case 'expression':
      // Expression nodes - use purple color
      nodeSize = 2.5;
      if (heatMapColor)
      {
        this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, heatMapColor.inner, heatMapColor.outer);
      } else
      {
        this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, "#9b59b6", "#8e44ad");
      }
      break;
    case 'relation':
      // Relation nodes - use gray color, smaller size
      nodeSize = 2;
      if (heatMapColor)
      {
        this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, heatMapColor.inner, heatMapColor.outer);
      } else
      {
        this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, "#95a5a6", "#7f8c8d");
      }
      break;
    default:
      // Fallback to generic rendering
      if (heatMapColor)
      {
        nodeSize = this.getNodeSize(item.type);
        this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, heatMapColor.inner, heatMapColor.outer);
      } else
      {
        const colors = this.getNodeColors(item.type);
        nodeSize = this.getNodeSize(item.type);
        this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, colors.fill, colors.stroke);
      }
    }

    // Draw label (smart truncation) - only if close enough
    let label = item.title || `Node ${index}`;
    if (label.length > 40)
    {
      // Try to find a good breaking point
      const words = label.split(' ');
      if (words.length > 3)
      {
        label = words.slice(0, 3).join(' ') + '...';
      } else
      {
        label = label.slice(0, 37) + '...';
      }
    }

    // Use smart label system
    this.drawSmartLabel(x, y, z, label, 8, "white", item.type, 'main', index);

    // Draw type indicator with appropriate color - only if very close
    const colors = this.getNodeColors(item.type);
    this.drawSmartLabel(x, y - 0.05, z, `[${item.type}]`, 6, colors.fill, item.type, 'type', index);

    // Draw coordinates info - only if extremely close
    this.drawSmartLabel(x, y - 0.08, z, `(${x.toFixed(2)}, ${y.toFixed(2)}, ${z.toFixed(2)})`, 5, "oklch(85.5% 0.138 181.071)", item.type, 'coords', index);

    // Draw context info if available - only if close
    if (item.context && this.showContext)
    {
      const contextLabel = item.context.length > 25 ? item.context.slice(0, 22) + '...' : item.context;
      this.drawSmartLabel(x, y + 0.05, z, contextLabel, 6, "lightblue", item.type, 'context', index);
    }

    // Highlight selected item
    if (this.selectedItem === index)
    {
      this.canvas3d.drawNode(x, y, z, (nodeSize + 2) * this.canvas3d.mob, "transparent", "#FFD700");
    }

    // Highlight hovered item with subtle glow
    if (this.hoveredNode && this.hoveredNode.index === index)
    {
      this.canvas3d.drawNode(x, y, z, (nodeSize + 1) * this.canvas3d.mob, "transparent", "#FFFFFF40");
    }

    // Highlight connected nodes when focusing
    if (this.focusedNode !== null)
    {
      if (this.connectedNodes.has(index))
      {
        // Connected node - highlight with green glow
        this.canvas3d.drawNode(x, y, z, (nodeSize + 1.5) * this.canvas3d.mob, "transparent", "#00FF0060");
      }
      else if (index !== this.focusedNode)
      {
        // Non-connected node - dim it
        this.canvas3d.ctx.save();
        this.canvas3d.ctx.globalAlpha = 0.4;
        // Re-draw with reduced opacity - this is a simple approach
        this.canvas3d.ctx.restore();
      }
    }
  }

  renderGenericNode(item, index)
  {
    const x = item.XYZ.X;
    const y = item.XYZ.Y;
    const z = item.XYZ.Z;

    let nodeSize = 4; // Default size

    // Get heat map color if enabled
    const itemId = item.id || index.toString();
    const heatMapColor = this.getNodeHeatMapColor(itemId);

    // Use specialized node drawing methods based on type if available
    if (item.type)
    {
      switch (item.type)
      {
      case 'event':
        this.canvas3d.drawEvent(x, y, z, heatMapColor);
        nodeSize = 6;
        break;
      case 'thing':
        this.canvas3d.drawThing(x, y, z, heatMapColor);
        nodeSize = 4;
        break;
      case 'concept':
        this.canvas3d.drawConcept(x, y, z, heatMapColor);
        nodeSize = 4;
        break;
      case 'agent':
        nodeSize = 3;
        if (heatMapColor)
        {
          this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, heatMapColor.inner, heatMapColor.outer);
        } else
        {
          this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, "#f39c12", "#e67e22");
        }
        break;
      case 'expression':
        nodeSize = 2.5;
        if (heatMapColor)
        {
          this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, heatMapColor.inner, heatMapColor.outer);
        } else
        {
          this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, "#9b59b6", "#8e44ad");
        }
        break;
      case 'relation':
        nodeSize = 2;
        if (heatMapColor)
        {
          this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, heatMapColor.inner, heatMapColor.outer);
        } else
        {
          this.canvas3d.drawNode(x, y, z, nodeSize * this.canvas3d.mob, "#95a5a6", "#7f8c8d");
        }
        break;
      default:
        // Default generic rendering
        if (heatMapColor)
        {
          this.canvas3d.drawNode(x, y, z, 4 * this.canvas3d.mob, heatMapColor.inner, heatMapColor.outer);
        } else
        {
          this.canvas3d.drawNode(x, y, z, 4 * this.canvas3d.mob, "#3498db", "#2980b9");
        }
      }
    }
    else
    {
      // Default generic rendering when no type is specified
      if (heatMapColor)
      {
        this.canvas3d.drawNode(x, y, z, 4 * this.canvas3d.mob, heatMapColor.inner, heatMapColor.outer);
      } else
      {
        this.canvas3d.drawNode(x, y, z, 4 * this.canvas3d.mob, "#3498db", "#2980b9");
      }
    }

    // Use smart label system for title
    const title = item.title || `Node ${index}`;
    this.drawSmartLabel(x, y, z, title.slice(0, 25), 8, "white", item.type || 'generic', 'main', index);

    // Draw type indicator if available
    if (item.type)
    {
      const colors = this.getNodeColors(item.type);
      this.drawSmartLabel(x, y - 0.05, z, `[${item.type}]`, 6, colors.fill, item.type, 'type', index);
    }

    // Draw coordinates info for debugging
    this.drawSmartLabel(x, y - 0.08, z, `(${x.toFixed(2)}, ${y.toFixed(2)}, ${z.toFixed(2)})`, 5, "oklch(85.5% 0.138 181.071)", item.type || 'generic', 'coords', index);

    // Highlight selected item
    if (this.selectedItem === index)
    {
      this.canvas3d.drawNode(x, y, z, (nodeSize + 2) * this.canvas3d.mob, "transparent", "#FFD700");
    }

    // Highlight hovered item with subtle glow
    if (this.hoveredNode && this.hoveredNode.index === index)
    {
      this.canvas3d.drawNode(x, y, z, (nodeSize + 1) * this.canvas3d.mob, "transparent", "#FFFFFF40");
    }

    // Highlight connected nodes when focusing
    if (this.focusedNode !== null)
    {
      if (this.connectedNodes.has(index))
      {
        // Connected node - highlight with green glow
        this.canvas3d.drawNode(x, y, z, (nodeSize + 1.5) * this.canvas3d.mob, "transparent", "#00FF0060");
      }
    }
  }

  getNodeColors(type)
  {
    const colorMap = {
      'relation': { fill: '#95a5a6', stroke: '#7f8c8d' },      // Gray
      'event': { fill: '#e74c3c', stroke: '#c0392b' },        // Red  
      'agent': { fill: '#f39c12', stroke: '#e67e22' },        // Orange
      'thing': { fill: '#27ae60', stroke: '#2ecc71' },        // Green
      'concept': { fill: '#3498db', stroke: '#2980b9' },      // Blue
      'expression': { fill: '#9b59b6', stroke: '#8e44ad' },   // Purple
      'chapter': { fill: '#e74c3c', stroke: '#c0392b' },      // Red for chapters
      'context': { fill: '#3498db', stroke: '#2980b9' },      // Blue for context
      'generic': { fill: '#34495e', stroke: '#2c3e50' }       // Dark gray
    };

    return colorMap[type] || colorMap['generic'];
  }

  getNodeSize(type)
  {
    const sizeMap = {
      'relation': 2,
      'event': 4,
      'agent': 4,
      'thing': 4,
      'concept': 4,
      'expression': 3,
      'chapter': 6,
      'context': 3,
      'generic': 3
    };

    return sizeMap[type] || 3;
  }

  renderArrow(arrow, highlight = null)
  {
    // Render arrow between two nodes using the SSTCanvas3D arrow system
    // (7 semantic relationship types mapped to 4 visual arrow styles)
    if (arrow.from < this.processedData.items.length && arrow.to < this.processedData.items.length)
    {
      const fromItem = this.processedData.items[arrow.from];
      const toItem = this.processedData.items[arrow.to];

      if (fromItem && toItem && fromItem.XYZ && toItem.XYZ)
      {
        const arrowType = this.convertArrowType(arrow.type);

        // Set canvas context for highlighting
        if (highlight !== null)
        {
          this.canvas3d.ctx.save();
          if (highlight === true)
          {
            // Enhanced visibility for connected arrows
            this.canvas3d.ctx.globalAlpha = 1.0;
            this.canvas3d.ctx.shadowColor = "#FFD700";
            this.canvas3d.ctx.shadowBlur = 8;
          }
          else
          {
            // Reduced visibility for non-connected arrows
            this.canvas3d.ctx.globalAlpha = 0.3;
          }
        }

        // Use the specific arrow methods based on converted type
        switch (arrowType)
        {
        case 'leadsTo':
          this.canvas3d.drawLeadsToArrow(
            fromItem.XYZ.X, fromItem.XYZ.Y, fromItem.XYZ.Z,
            toItem.XYZ.X, toItem.XYZ.Y, toItem.XYZ.Z
          );
          break;
        case 'contains':
          this.canvas3d.drawContainsArrow(
            fromItem.XYZ.X, fromItem.XYZ.Y, fromItem.XYZ.Z,
            toItem.XYZ.X, toItem.XYZ.Y, toItem.XYZ.Z
          );
          break;
        case 'expresses':
          this.canvas3d.drawExpressesArrow(
            fromItem.XYZ.X, fromItem.XYZ.Y, fromItem.XYZ.Z,
            toItem.XYZ.X, toItem.XYZ.Y, toItem.XYZ.Z
          );
          break;
        case 'near':
          this.canvas3d.drawNearArrow(
            fromItem.XYZ.X, fromItem.XYZ.Y, fromItem.XYZ.Z,
            toItem.XYZ.X, toItem.XYZ.Y, toItem.XYZ.Z
          );
          break;
        default:
          // Default to LeadsTo if unknown type
          this.canvas3d.drawLeadsToArrow(
            fromItem.XYZ.X, fromItem.XYZ.Y, fromItem.XYZ.Z,
            toItem.XYZ.X, toItem.XYZ.Y, toItem.XYZ.Z
          );
        }

        // Restore canvas context if highlighting was applied
        if (highlight !== null)
        {
          this.canvas3d.ctx.restore();
        }
      }
    }
  }

  convertArrowType(arrowValue)
  {
    // Convert from SSTorytime arrow values to SSTCanvas3D arrow system
    // (7 semantic relationship types mapped to 4 visual arrow styles)

    if (typeof arrowValue === 'string')
    {
      // Handle string-based arrow types
      switch (arrowValue.toLowerCase())
      {
      case 'leadsto':
      case 'leads_to':
      case 'sequence':
      case 'flow':
        return 'leadsTo';
      case 'contains':
      case 'containment':
      case 'hierarchy':
        return 'contains';
      case 'expresses':
      case 'expression':
      case 'meaning':
        return 'expresses';
      case 'near':
      case 'similarity':
      case 'proximity':
        return 'near';
      default:
        return 'leadsTo';
      }
    }

    if (arrowValue >= 298)
    {
      // High values like 298, 276, 258 - convert to semantic types
      switch (arrowValue)
      {
      case 298: return 'leadsTo';    // has value - sequence/flow
      case 299: return 'leadsTo';    // is a - sequence/flow
      case 300: return 'contains';   // contains - hierarchy
      case 301: return 'leadsTo';    // follows - sequence
      case 302: return 'expresses';  // expresses - meaning
      case 303: return 'contains';   // promises - commitment
      case 276: return 'expresses';  // expression relationship
      case 258: return 'near';       // general proximity
      default: return 'leadsTo';     // default
      }
    } else if (arrowValue >= 0 && arrowValue <= 10)
    {
      // Direct semantic types (0-10 range) - map to 7 semantic types then to 4 visual types
      switch (arrowValue)
      {
      case 0: return 'expresses';    // Im3: "is a property expressed by"
      case 1: return 'contains';     // Im2: "is contained by"
      case 2: return 'leadsTo';      // Im1: "comes from"
      case 3: return 'near';         // In0: "is near/similar to"
      case 4: return 'leadsTo';      // Il1: "leads to"
      case 5: return 'contains';     // Ic2: "contains"
      case 6: return 'expresses';    // Ie3: "expresses property"
      default: return 'leadsTo';
      }
    } else
    {
      // Default fallback
      return 'leadsTo';
    }
  }

  drawArrows()
  {
    if (!this.processedData.arrows) return;

    this.processedData.arrows.forEach((arrow, index) =>
    {
      // Check if this arrow is connected to the focused node
      const isConnected = this.connectedArrows.has(index);

      if (this.focusedNode !== null)
      {
        // If we have a focused node, dim non-connected arrows
        if (isConnected)
        {
          // Render connected arrows with enhanced visibility
          this.renderArrow(arrow, true);
        }
        else
        {
          // Render non-connected arrows with reduced opacity
          this.renderArrow(arrow, false);
        }
      }
      else
      {
        // Normal rendering when no node is focused
        this.renderArrow(arrow, null);
      }
    });
  }

  renderTOCContextFragments(chapter, centerX, centerY, centerZ, chapterIndex)
  {
    // Original TOC context rendering
    if (chapter.Context && chapter.Context.length > 0)
    {
      chapter.Context.forEach((context, i) =>
      {
        if (context && context.XYZ)
        {
          const cx = context.XYZ.X;
          const cy = context.XYZ.Y;
          const cz = context.XYZ.Z;

          // Use specialized context node rendering - concepts work best for context
          this.canvas3d.drawConcept(cx, cy, cz);

          // Add context type label - show if parent chapter is focused/connected or normal distance rules apply
          this.drawSmartLabel(cx, cy - 0.03, cz, "[context]", 5, "lightblue", 'context', 'type', chapterIndex);

          // Connect to main chapter
          this.canvas3d.drawLine3D(centerX, centerY, centerZ, cx, cy, cz, "oklch(68.5% 0.169 237.323)", 0.9);
        }
      });
    }

    // Render Single fragments (green)
    if (chapter.Single && chapter.Single.length > 0)
    {
      chapter.Single.forEach((single, i) =>
      {
        if (single && single.XYZ)
        {
          const sx = single.XYZ.X;
          const sy = single.XYZ.Y;
          const sz = single.XYZ.Z;

          // Use specialized single node rendering - things work best for single items
          this.canvas3d.drawThing(sx, sy, sz);

          // Add single type label - show if parent chapter is focused/connected or normal distance rules apply
          this.drawSmartLabel(sx, sy - 0.03, sz, "[single]", 5, "lightgreen", 'fragment', 'type', chapterIndex);

          // Connect to main chapter
          this.canvas3d.drawLine3D(centerX, centerY, centerZ, sx, sy, sz, "oklch(76.9% 0.17 142.685)", 0.7);
        }
      });
    }
  }

  drawTOCConnections()
  {
    // Optional: draw connections between TOC chapters
    // This could show chapter relationships or reading order
  }

  renderChapter(chapter, index)
  {
    const x = chapter.XYZ.X;
    const y = chapter.XYZ.Y;
    const z = chapter.XYZ.Z;

    // Draw main chapter node (larger, red)
    this.canvas3d.drawNode(x, y, z, 5 * this.canvas3d.mob, "red", "red");

    // Draw chapter title (handle empty chapter names)
    const chapterTitle = chapter.Chapter || `Chapter ${index + 1}`;
    this.canvas3d.drawLabel(x, y, z, chapterTitle.slice(0, 40), 12, "white");

    // Draw chapter coordinates info
    this.canvas3d.drawLabel(x, y - 0.05, z, `[${x.toFixed(2)}, ${y.toFixed(2)}, ${z.toFixed(2)}]`, 7, "oklch(85.5% 0.138 181.071)");

    // Render context fragments if enabled
    if (this.showContext)
    {
      this.renderContextFragments(chapter, x, y, z);
    }

    // Highlight selected chapter
    if (this.selectedItem === index)
    {
      this.canvas3d.drawNode(x, y, z, 7 * this.canvas3d.mob, "transparent", "#FFD700");
    }
  }

  renderContextFragments(chapter, centerX, centerY, centerZ)
  {
    // Render Context fragments (blue)
    if (chapter.Context && chapter.Context.length > 0)
    {
      chapter.Context.forEach((context, i) =>
      {
        if (context && context.XYZ)
        {
          const cx = context.XYZ.X;
          const cy = context.XYZ.Y;
          const cz = context.XYZ.Z;

          // Draw context node
          this.canvas3d.drawNode(cx, cy, cz, 3 * this.canvas3d.mob, "#2980b9", "#3498db");

          // Draw context text if available
          // if (context.Text)
          // {
          //   this.canvas3d.drawLabel(cx, cy, cz, context.Text.slice(0, 15), 7, "lightblue");
          // }

          // Connect to main chapter
          this.canvas3d.drawLine3D(centerX, centerY, centerZ, cx, cy, cz, "oklch(68.5% 0.169 237.323)", 0.9);
        }
      });
    }

    // Render Single contexts (orange)
    if (chapter.Single && chapter.Single.length > 0)
    {
      chapter.Single.forEach((single, i) =>
      {
        if (single && single.XYZ)
        {
          const sx = single.XYZ.X;
          const sy = single.XYZ.Y;
          const sz = single.XYZ.Z;

          // Draw single context node
          this.canvas3d.drawNode(sx, sy, sz, 2.5 * this.canvas3d.mob, "#e67e22", "#f39c12");

          // Draw single context text if available
          // if (single.Text)
          // {
          //   this.canvas3d.drawLabel(sx, sy, sz, single.Text.slice(0, 12), 7, "orange");
          // }

          // Connect to main chapter with dashed line effect
          this.canvas3d.drawLine3D(centerX, centerY, centerZ, sx, sy, sz, "oklch(92.4% 0.12 95.746)", 0.7);
        }
      });
    }

    // Render Common contexts (green)
    if (chapter.Common && chapter.Common.length > 0)
    {
      chapter.Common.forEach((common, i) =>
      {
        if (common && common.XYZ)
        {
          const gx = common.XYZ.X;
          const gy = common.XYZ.Y;
          const gz = common.XYZ.Z;

          // Draw common context node
          this.canvas3d.drawNode(gx, gy, gz, 2 * this.canvas3d.mob, "oklch(82.8% 0.189 84.429)", "oklch(64.8% 0.2 131.684)");

          // Draw common context text if available
          // if (common.Text)
          // {
          //   this.canvas3d.drawLabel(gx, gy, gz, common.Text.slice(0, 12), 7, "lightgreen");
          // }

          // Connect to main chapter
          this.canvas3d.drawLine3D(centerX, centerY, centerZ, gx, gy, gz, "oklch(96.7% 0.067 122.328)", 0.1);
        }
      });
    }
  }

  drawChapterConnections()
  {
    // Draw connections between items based on proximity or shared context
    if (!this.processedData || !this.processedData.items) return;

    const items = this.processedData.items;

    for (let i = 0; i < items.length - 1; i++)
    {
      const item1 = items[i];
      const item2 = items[i + 1];

      if (chapter1 && chapter2 && chapter1.XYZ && chapter2.XYZ)
      {
        // Calculate distance
        const dx = chapter2.XYZ.X - chapter1.XYZ.X;
        const dy = chapter2.XYZ.Y - chapter1.XYZ.Y;
        const dz = chapter2.XYZ.Z - chapter1.XYZ.Z;
        const distance = Math.sqrt(dx * dx + dy * dy + dz * dz);

        // Draw connection if chapters are reasonably close
        if (distance < 1.0)
        {
          this.canvas3d.drawLine3D(
            chapter1.XYZ.X, chapter1.XYZ.Y, chapter1.XYZ.Z, chapter2.XYZ.X, chapter2.XYZ.Y,
            chapter2.XYZ.Z, "rgba(255,255,255,0.5)", 1 * this.canvas3d.mob
          );
        }
      }
    }
  }

  updateInfoPanel()
  {
    console.log("📊 Updating info panel with:", this.processedData);

    if (!this.processedData)
    {
      console.log("⚠️ No processed data for info panel");
      return;
    }

    // Update time and intent info if available
    document.getElementById("currentTime").textContent = this.processedData.metadata.time || "Unknown";
    document.getElementById("intentInfo").textContent = this.processedData.metadata.intent || "None";
    document.getElementById("ambientInfo").textContent = this.processedData.metadata.ambient || "None";

    // Update statistics based on data type
    this.updateStatistics();
  }

  updateStatistics()
  {
    console.log("📈 Updating statistics with:", this.processedData);

    if (!this.processedData)
    {
      console.log("⚠️ No processed data for statistics");
      return;
    }

    const items = this.processedData.items || [];
    const arrows = this.processedData.arrows || [];

    console.log(`📊 Items: ${items.length}, Arrows: ${arrows.length}, Data type: ${this.processedData.dataType}`);

    if (this.processedData.dataType === 'toc')
    {
      // TOC-specific statistics
      let totalContexts = 0;
      let totalFragments = 0;

      items.forEach((chapter) =>
      {
        if (chapter.Context) totalContexts += chapter.Context.length;
        if (chapter.Single) totalFragments += chapter.Single.length;
        if (chapter.Common) totalFragments += chapter.Common.length;
      });

      document.getElementById("chapterCount").textContent = items.length;
      document.getElementById("contextCount").textContent = totalContexts;
      document.getElementById("fragmentCount").textContent = totalFragments;

    } else if (this.processedData.dataType === 'orbit')
    {
      // Orbit/PageMap statistics
      const typeCount = {};
      items.forEach(item =>
      {
        typeCount[item.type] = (typeCount[item.type] || 0) + 1;
      });

      document.getElementById("chapterCount").textContent = items.length;
      document.getElementById("contextCount").textContent = arrows.length;
      document.getElementById("fragmentCount").textContent = Object.keys(typeCount).length;

    } else
    {
      // Generic statistics
      document.getElementById("chapterCount").textContent = items.length;
      document.getElementById("contextCount").textContent = arrows.length;
      document.getElementById("fragmentCount").textContent = this.processedData.metadata.type || 'Generic';
    }
  }

  updateItemList()
  {
    const itemList = document.getElementById("chapterList");
    if (!this.processedData || !this.processedData.items)
    {
      itemList.innerHTML = '<div class="loading">No items found</div>';
      return;
    }

    itemList.innerHTML = "";
    console.log(`📋 Updating item list with ${this.processedData.items.length} items`);

    this.processedData.items.forEach((item, index) =>
    {
      const itemElement = document.createElement("div");
      itemElement.className = "chapter-item";
      itemElement.dataset.index = index;

      let title, coords, subtitle = '';

      if (this.processedData.dataType === 'toc')
      {
        title = item.title || item.Chapter || `Chapter ${index + 1}`;
        coords = `(${item.XYZ.X.toFixed(2)}, ${item.XYZ.Y.toFixed(2)}, ${item.XYZ.Z.toFixed(2)})`;
        subtitle = item.type ? `Type: ${item.type}` : (item.Context ? `${item.Context.length} contexts` : '');
      } else if (this.processedData.dataType === 'orbits')
      {
        // For orbit data, truncate long titles but keep them meaningful
        const fullTitle = item.title || `Item ${index + 1}`;
        title = fullTitle.length > 50 ? fullTitle.substring(0, 47) + '...' : fullTitle;
        coords = `(${item.XYZ.X.toFixed(2)}, ${item.XYZ.Y.toFixed(2)}, ${item.XYZ.Z.toFixed(2)})`;
        subtitle = `${item.type} | ${item.chapter || 'No Chapter'}`;
      } else
      {
        title = item.title || item.name || `Item ${index + 1}`;
        coords = `(${item.XYZ.X.toFixed(2)}, ${item.XYZ.Y.toFixed(2)}, ${item.XYZ.Z.toFixed(2)})`;
        subtitle = item.type ? `Type: ${item.type}` : '';
      }

      console.log(`📝 Item ${index}: ${title}`);

      itemElement.innerHTML = `
        <div class="title" title="${item.title || title}">${title}</div>
        ${subtitle ? `<div class="subtitle">${subtitle}</div>` : ''}
        <div class="coords">${coords}</div>
      `;

      itemElement.addEventListener("click", () =>
      {
        this.selectItem(index);
      });

      itemList.appendChild(itemElement);
    });
  }

  selectItem(index)
  {
    this.selectedItem = index;

    // Update visual selection in list
    document.querySelectorAll(".chapter-item").forEach((item, i) =>
    {
      item.style.background = i === index ? "rgba(255,215,0,0.3)" : "rgba(255,255,255,0.1)";
      item.style.borderLeftColor = i === index ? "#FFD700" : "#3498db";
    });

    // Focus on selected item
    if (this.processedData.items[index])
    {
      const item = this.processedData.items[index];
      this.canvas3d.setObserverPosition(item.XYZ.X + 0.5, item.XYZ.Y + 0.3, item.XYZ.Z - 1);
    }

    // Re-render to highlight selection
    if (!this.isAnimating)
    {
      this.renderVisualization();
    }
  }

  setupControls()
  {
    // Data source selector
    document
      .getElementById("dataSource")
      .addEventListener("change", async (e) =>
      {
        try
        {
          await this.loadData(e.target.value);
        } catch (error)
        {
          console.error("Failed to load data:", error);
          alert("Failed to load data: " + error.message);
        }
      });

    // Rotation speed control
    document
      .getElementById("rotationSpeed")
      .addEventListener("input", (e) =>
      {
        this.rotationSpeed = parseFloat(e.target.value);
        document.getElementById("speedValue").textContent =
          e.target.value;
      });

    // View angle control
    document
      .getElementById("viewAngle")
      .addEventListener("input", (e) =>
      {
        const angle = parseFloat(e.target.value);
        this.canvas3d.setViewingAngle(angle, this.canvas3d.phi);
        document.getElementById("angleValue").textContent = e.target.value;
        if (!this.isAnimating)
        {
          this.renderVisualization();
        }
      });

    // Smart labels toggle
    document
      .getElementById("smartLabels")
      .addEventListener("change", (e) =>
      {
        this.useSmartLabels = e.target.checked;
        if (!this.isAnimating)
        {
          this.renderVisualization();
        }
      });

    // Animation toggle
    document
      .getElementById("toggleAnimation")
      .addEventListener("click", () =>
      {
        if (this.isAnimating)
        {
          this.stopAnimation();
          document.getElementById("toggleAnimation").textContent = "▶️ Play";
        } else
        {
          this.startAnimation();
          document.getElementById("toggleAnimation").textContent =
            "⏸️ Pause";
        }
      });

    // Reset view
    document.getElementById("resetView").addEventListener("click", () =>
    {
      // Clear any focus first
      this.clearFocus();

      this.currentAngle = 0;
      this.canvas3d.setViewingAngle(Math.PI / 10, Math.PI / 10);
      this.canvas3d.setObserverPosition(1.5, 0.75, -1.5);
      document.getElementById("viewAngle").value = Math.PI / 10;
      document.getElementById("angleValue").textContent = (Math.PI / 10).toFixed(2);
      this.selectedItem = null;
      this.updateItemList();

      if (!this.isAnimating)
      {
        this.renderVisualization();
      }
    });

    // Focus item
    document
      .getElementById("focusChapter")
      .addEventListener("click", () =>
      {
        if (this.focusedNode !== null)
        {
          // If already focused, clear focus
          this.clearFocus();
        }
        else if (this.selectedItem !== null)
        {
          // Focus on the selected item
          this.focusOnNode(this.selectedItem);
        }
        else
        {
          alert("Please click on a node to focus, or select from the list first");
        }
      });

    // Heat map toggle
    document
      .getElementById("heatMapToggle")
      .addEventListener("change", (e) =>
      {
        this.useHeatMap = e.target.checked;

        // Show/hide heat map controls
        const heatMapControls = document.getElementById("heatMapControls");
        if (e.target.checked)
        {
          heatMapControls.style.display = "flex";
          this.calculateHeatMap();
        } else
        {
          heatMapControls.style.display = "none";
          this.heatMapColors.clear();
        }

        if (!this.isAnimating)
        {
          this.renderVisualization();
        }
      });

    // Heat map center selector
    document
      .getElementById("heatMapCenter")
      .addEventListener("change", (e) =>
      {
        this.heatMapCenter = e.target.value;
        if (this.useHeatMap)
        {
          this.calculateHeatMap();
          if (!this.isAnimating)
          {
            this.renderVisualization();
          }
        }
      });

    // Control panel toggle
    document
      .getElementById("toggleControls")
      .addEventListener("click", () =>
      {
        const content = document.getElementById("controlContent");
        const button = document.getElementById("toggleControls");

        if (content.classList.contains("collapsed"))
        {
          content.classList.remove("collapsed");
          button.textContent = "−";
        } else
        {
          content.classList.add("collapsed");
          button.textContent = "+";
        }
      });

    // Window resize handler
    window.addEventListener("resize", () =>
    {
      this.canvas3d.resizeCanvas(
        Math.min(window.innerWidth, 400),
        Math.min(window.innerHeight, 700),
      );
      if (!this.isAnimating)
      {
        this.renderVisualization();
      }
    });
  }

  startAnimation()
  {
    this.isAnimating = true;

    const animate = () =>
    {
      if (!this.isAnimating) return;

      // Update rotation
      this.currentAngle += this.rotationSpeed;
      this.canvas3d.setViewingAngle(this.currentAngle, this.canvas3d.phi);

      // Re-render
      this.renderVisualization();

      // Continue animation
      this.animationId = requestAnimationFrame(animate);
    };

    animate();
  }

  stopAnimation()
  {
    this.isAnimating = false;
    if (this.animationId)
    {
      cancelAnimationFrame(this.animationId);
      this.animationId = null;
    }
  }

  // === HEAT MAP SYSTEM ===

  /**
   * Calculate heat map colors for all nodes based on distance from center
   */
  calculateHeatMap()
  {
    if (!this.processedData || !this.processedData.items || !this.heatMapCenter)
    {
      console.log("📊 Heat map: No data or center node selected");
      return;
    }

    console.log(`📊 Calculating heat map from center: ${this.heatMapCenter}`);

    // Find the center node
    const centerNode = this.processedData.items.find((item, index) =>
      (item.id && item.id === this.heatMapCenter) || index.toString() === this.heatMapCenter
    );
    if (!centerNode)
    {
      console.warn("❌ Heat map center node not found:", this.heatMapCenter);
      return;
    }

    // Determine which nodes to process
    let nodesToProcess = [];

    if (this.focusedNode !== null && this.connectedNodes.size > 0)
    {
      // In focus mode: only process connected nodes + the focused node
      console.log(`🎯 Heat map in focus mode: processing ${this.connectedNodes.size + 1} connected nodes`);

      // Add the focused node
      nodesToProcess.push({
        item: this.processedData.items[this.focusedNode],
        index: this.focusedNode
      });

      // Add all connected nodes
      this.connectedNodes.forEach(nodeIndex =>
      {
        if (nodeIndex < this.processedData.items.length)
        {
          nodesToProcess.push({
            item: this.processedData.items[nodeIndex],
            index: nodeIndex
          });
        }
      });
    }
    else
    {
      // Normal mode: process all nodes
      console.log(`📊 Heat map in normal mode: processing all ${this.processedData.items.length} nodes`);
      nodesToProcess = this.processedData.items.map((item, index) => ({ item, index }));
    }

    // Calculate all distances and find maximum
    const distances = new Map();
    this.maxDistance = 0;

    nodesToProcess.forEach(({ item, index }) =>
    {
      const itemId = item.id || index.toString();

      if (itemId === this.heatMapCenter)
      {
        distances.set(itemId, 0);
        return;
      }

      const distance = this.canvas3d.calculate3DDistance(
        centerNode.XYZ.X, centerNode.XYZ.Y, centerNode.XYZ.Z,
        item.XYZ.X, item.XYZ.Y, item.XYZ.Z
      );

      distances.set(itemId, distance);
      this.maxDistance = Math.max(this.maxDistance, distance);
    });

    // Calculate normalized distances and colors
    this.heatMapColors.clear();
    distances.forEach((distance, nodeId) =>
    {
      const normalizedDistance = this.maxDistance > 0 ? distance / this.maxDistance : 0;
      const color = this.canvas3d.getHeatMapColor(normalizedDistance);
      this.heatMapColors.set(nodeId, color);
    });

    const mode = this.focusedNode !== null && this.connectedNodes.size > 0 ? "focus" : "normal";
    console.log(`📊 Heat map calculated (${mode} mode): ${this.heatMapColors.size} nodes, max distance: ${this.maxDistance.toFixed(2)}`);
  }

  /**
   * Get heat map color for a node (if heat map is enabled)
   */
  getNodeHeatMapColor(nodeId)
  {
    if (!this.useHeatMap || !this.heatMapColors.has(nodeId))
    {
      return null;
    }
    return this.heatMapColors.get(nodeId);
  }

  /**
   * Activate heat map with focused node as center
   */
  activateHeatMapForFocus(nodeIndex)
  {
    if (!this.processedData || !this.processedData.items[nodeIndex]) return;

    console.log(`🌡️ Activating heat map for focused node ${nodeIndex}`);

    // Get the focused item
    const focusedItem = this.processedData.items[nodeIndex];
    const itemId = focusedItem.id || nodeIndex.toString();

    // Enable heat map
    this.useHeatMap = true;
    this.heatMapCenter = itemId;

    // Update UI controls
    const heatMapToggle = document.getElementById("heatMapToggle");
    const heatMapCenter = document.getElementById("heatMapCenter");
    const heatMapControls = document.getElementById("heatMapControls");

    if (heatMapToggle) heatMapToggle.checked = true;
    if (heatMapCenter) heatMapCenter.value = itemId;
    if (heatMapControls) heatMapControls.style.display = "flex";

    // Calculate heat map
    this.calculateHeatMap();

    console.log(`🎯 Heat map activated with center: ${itemId}`);
  }

  /**
   * Update heat map center options based on current data
   */
  updateHeatMapOptions()
  {
    const selector = document.getElementById("heatMapCenter");
    if (!selector || !this.processedData || !this.processedData.items) return;

    // Clear existing options
    selector.innerHTML = '<option value="">Select center node...</option>';

    // Add all nodes as options
    this.processedData.items.forEach((item, index) =>
    {
      const option = document.createElement("option");
      option.value = item.id || index.toString();
      option.textContent = `${item.title || item.Chapter || item.id || `Node ${index}`} (${item.type || 'node'})`;
      selector.appendChild(option);
    });

    // Auto-select first node if none selected
    if (!this.heatMapCenter && this.processedData.items.length > 0)
    {
      this.heatMapCenter = this.processedData.items[0].id || "0";
      selector.value = this.heatMapCenter;
    }
  }

  showError(message)
  {
    const container = document.getElementById("canvasContainer");
    container.innerHTML = `
          <div class="error">
            <h3>⚠️ Error</h3>
             <p>${message}</p>
             <button onclick="location.reload()" style="padding: 10px 20px; margin-top: 15px; background: #3498db; color: white; border: none; border-radius: 5px; cursor: pointer;">
                 🔄 Reload Page
             </button>
          </div>
      `;
  }
}

// Initialize when page loads
document.addEventListener("DOMContentLoaded", () =>
{
  window.visualization = new GenericAPIVisualization();
});