# SSTCanvas3D Graphics Library - Moved

The SSTCanvas3D graphics library and all related demo files have been moved to a dedicated directory for better organization.

## 📁 New Location

All SSTCanvas3D files are now located in:
```
/home/alex/SSTorytime/src/server/sst-canvas/
```

## 📋 What Was Moved

The following files have been relocated:

### Core Library Files
- `sstcanvas3d.js` → `src/server/sst-canvas/sstcanvas3d.js`
- `sstcanvas3d-examples.js` → `src/server/sst-canvas/sstcanvas3d-examples.js`
- `README-sstcanvas3d.md` → `src/server/sst-canvas/README-sstcanvas3d.md`

### Demo & Examples
- `chapter-visualization.html` → `src/server/sst-canvas/chapter-visualization.html`
- `chapter-data.json` → `src/server/sst-canvas/chapter-data.json`
- `chapter-data-enhanced.json` → `src/server/sst-canvas/chapter-data-enhanced.json`

### Data Extraction Tools
- `extract-chapter-data.sh` → `src/server/sst-canvas/extract-chapter-data.sh`
- `extract-chapter-data.js` → `src/server/sst-canvas/extract-chapter-data.js`

### Documentation
- `README-chapter-example.md` → `src/server/sst-canvas/README.md` (comprehensive guide)

## 🚀 Quick Access

### View the Demo
```bash
cd /home/alex/SSTorytime/src/server/sst-canvas
python3 -m http.server 8000
# Open: http://localhost:8000/chapter-visualization.html
```

### Extract Fresh Data
```bash
cd /home/alex/SSTorytime/src/server/sst-canvas
./extract-chapter-data.sh
```

### Read Documentation
```bash
cd /home/alex/SSTorytime/src/server/sst-canvas
cat README.md                    # Overview and usage guide
cat README-sstcanvas3d.md       # Complete API documentation
```

## 🎯 Benefits of New Organization

✅ **Cleaner Structure** - All related files in one dedicated directory
✅ **Better Documentation** - Comprehensive README with all details
✅ **Easier Discovery** - Clear location for graphics library components
✅ **Improved Maintenance** - Self-contained module for easier updates
✅ **Enhanced Organization** - Separates demo code from core server files

For complete documentation and usage instructions, see the files in `/home/alex/SSTorytime/src/server/sst-canvas/`.