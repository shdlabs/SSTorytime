class ChapterTOCVisualization
{
  constructor()
  {
    this.canvas3d = null;
    this.tocData = null;
    this.animationId = null;
    this.isAnimating = true;
    this.rotationSpeed = 0.01;
    this.currentAngle = 0;
    this.showContext = true;
    this.selectedChapter = null;

    this.init();
  }

  async init()
  {
    try
    {
      console.log("🚀 Initializing Chapter TOC Visualization...");

      // Initialize the 3D canvas
      this.canvas3d = new SSTCanvas3D("canvasContainer", {
        width: window.innerWidth,
        height: window.innerHeight,
        mobile: window.innerWidth < 550 ? 0.5 : 1,
      });

      // Load initial data
      const data = document.getElementById("dataSource");
      await this.loadData(data.value);

      // Set up controls
      this.setupControls();

      // Start animation
      this.startAnimation();

      console.log("✅ Chapter TOC visualization initialized successfully");
    } catch (error)
    {
      console.error("❌ Failed to initialize visualization:", error);
      this.showError("Failed to load visualization: " + error.message);
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

      this.tocData = await response.json();
      this.updateInfoPanel();
      this.updateChapterList();
      this.renderVisualization();

      console.log("✅ TOC data loaded:", this.tocData);
    } catch (error)
    {
      console.error("❌ Error loading data:", error);
      throw error;
    }
  }

  renderVisualization()
  {
    if (!this.tocData || !this.tocData.Content)
    {
      console.warn("⚠️ No TOC data to render");
      return;
    }

    // Clear canvas
    this.canvas3d.clear();

    // Draw coordinate grid
    this.canvas3d.drawGrid(0, 0, 1);

    // Render each chapter
    this.tocData.Content.forEach((chapter, index) =>
    {
      this.renderChapter(chapter, index);
    });

    // Draw connections between chapters (optional)
    this.drawChapterConnections();
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
    if (this.selectedChapter === index)
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
    // Draw connections between chapters based on proximity or shared context
    const chapters = this.tocData.Content;

    for (let i = 0; i < chapters.length - 1; i++)
    {
      const chapter1 = chapters[i];
      const chapter2 = chapters[i + 1];

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
    if (!this.tocData) return;

    document.getElementById("currentTime").textContent = this.tocData.Time || "Unknown";
    document.getElementById("intentInfo").textContent = this.tocData.Intent || "None";
    document.getElementById("ambientInfo").textContent = this.tocData.Ambient || "None";

    // Calculate statistics
    const chapters = this.tocData.Content || [];
    let totalContexts = 0;
    let totalFragments = 0;

    chapters.forEach((chapter) =>
    {
      if (chapter.Context) totalContexts += chapter.Context.length;
      if (chapter.Single) totalFragments += chapter.Single.length;
      if (chapter.Common) totalFragments += chapter.Common.length;
    });

    document.getElementById("chapterCount").textContent = chapters.length;
    document.getElementById("contextCount").textContent = totalContexts;
    document.getElementById("fragmentCount").textContent = totalFragments;
  }

  updateChapterList()
  {
    const chapterList = document.getElementById("chapterList");
    if (!this.tocData || !this.tocData.Content)
    {
      chapterList.innerHTML = '<div class="loading">No chapters found</div>';
      return;
    }

    chapterList.innerHTML = "";

    this.tocData.Content.forEach((chapter, index) =>
    {
      const chapterItem = document.createElement("div");
      chapterItem.className = "chapter-item";
      chapterItem.dataset.index = index;

      const title = chapter.Chapter || `Unnamed Chapter ${index + 1}`;
      const coords = `(${chapter.XYZ.X.toFixed(2,)}, ${chapter.XYZ.Y.toFixed(2)}, ${chapter.XYZ.Z.toFixed(2)})`;

      chapterItem.innerHTML = `<div class="title">${title}</div> <div class="coords">${coords}</div>`;

      chapterItem.addEventListener("click", () =>
      {
        this.selectChapter(index);
      });

      chapterList.appendChild(chapterItem);
    });
  }

  selectChapter(index)
  {
    this.selectedChapter = index;

    // Update visual selection in list
    document.querySelectorAll(".chapter-item").forEach((item, i) =>
    {
      item.style.background = i === index ? "rgba(255,215,0,0.3)" : "rgba(255,255,255,0.1)";
      item.style.borderLeftColor = i === index ? "#FFD700" : "#3498db";
    });

    // Focus on selected chapter
    if (this.tocData.Content[index])
    {
      const chapter = this.tocData.Content[index];
      this.canvas3d.setObserverPosition(chapter.XYZ.X + 0.5, chapter.XYZ.Y + 0.3, chapter.XYZ.Z - 1);
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

    // Show context toggle
    document
      .getElementById("showContext")
      .addEventListener("change", (e) =>
      {
        this.showContext = e.target.checked;
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
      this.currentAngle = 0;
      this.canvas3d.setViewingAngle(Math.PI / 10, Math.PI / 10);
      this.canvas3d.setObserverPosition(1.5, 0.75, -1.5);
      document.getElementById("viewAngle").value = Math.PI / 10;
      document.getElementById("angleValue").textContent = (Math.PI / 10).toFixed(2);
      this.selectedChapter = null;
      this.updateChapterList();

      if (!this.isAnimating)
      {
        this.renderVisualization();
      }
    });

    // Focus chapter
    document
      .getElementById("focusChapter")
      .addEventListener("click", () =>
      {
        if (this.selectedChapter !== null)
        {
          const chapter = this.tocData.Content[this.selectedChapter];
          if (chapter)
          {
            this.canvas3d.setObserverPosition(chapter.XYZ.X, chapter.XYZ.Y, chapter.XYZ.Z - 1.5);

            if (!this.isAnimating)
            {
              this.renderVisualization();
            }
          }
        } else
        {
          alert("Please select a chapter from the list first");
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
  new ChapterTOCVisualization();
});