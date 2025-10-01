# SSTorytime Modern Frontend

A modern, interactive, mobile-first frontend for the SSTorytime knowledge exploration system.

## ✨ Features

### 🎯 Core Functionality
- **Smart Search Interface**: Intelligent search with autocomplete and command suggestions
- **Interactive Visualization**: 3D node graph visualization with touch/mobile support
- **Real-time API Integration**: Seamless connection to SSTorytime backend
- **Mobile-First Design**: Responsive layout optimized for all device sizes

### 🎨 User Experience
- **Modern UI**: Glass morphism design with smooth animations
- **Dark Theme**: Eye-friendly dark theme with accent colors
- **Progressive Web App**: Installable as a native-like app
- **Accessibility**: Keyboard navigation and screen reader support

### 🔍 Search Capabilities
- **Natural Language**: Search with plain text queries
- **Command System**: Use special commands like `\help`, `\chapters`, `\stats`
- **Context Aware**: Search within specific chapters or contexts
- **Quick Actions**: Pre-defined searches and navigation shortcuts

### 📱 Mobile Optimizations
- **Touch Interactions**: Optimized for touch devices
- **Responsive Layout**: Adapts to any screen size
- **Fast Performance**: Optimized bundle size and loading
- **Offline Ready**: PWA capabilities for offline use

## 🚀 Technology Stack

- **React 18**: Modern React with hooks and concurrent features
- **TypeScript**: Full type safety and better developer experience
- **Vite**: Fast build tool and development server
- **Tailwind CSS**: Utility-first CSS framework
- **Framer Motion**: Smooth animations and interactions
- **D3.js**: Data visualization and SVG manipulation
- **Axios**: HTTP client for API requests

## 🛠️ Development

### Prerequisites
- Node.js 18+ and npm
- SSTorytime backend server running on port 8080

### Installation
```bash
cd frontend
npm install
```

### Development Server
```bash
npm run dev
```
Visit http://localhost:3000

### Build for Production
```bash
npm run build
```

### Preview Production Build
```bash
npm run preview
```

## 📡 API Integration

The frontend integrates with the SSTorytime backend API:

- **Endpoint**: `POST /searchN4L`
- **Proxy**: Vite dev server proxies `/searchN4L` to `localhost:8080`
- **Response Types**: Handles Orbits, PageMap, TOC, Sequence, PathSolve, Arrows, STAT, GUIDANCE, ERROR

### Example API Call
```typescript
const result = await searchAPI({ 
  name: "moon \\context astronomy" 
});
```

## 🎮 Usage Examples

### Basic Search
- Type "moon" to search for moon-related content
- Use "recipe fish soup" for multi-word searches

### Special Commands
- `\\help` - Get search guidance
- `\\chapters` - Browse all chapters  
- `\\stats` - View statistics
- `\\notes brain` - Get notes for a topic

### Advanced Queries
- `moon \\context astronomy` - Search with context
- `\\story Mary` - Get story sequences
- `\\from start \\to target` - Find paths between nodes

## 📂 Project Structure

```
frontend/
├── src/
│   ├── components/          # React components
│   │   ├── SearchBar.tsx    # Smart search interface
│   │   ├── Visualization.tsx # 3D graph visualization
│   │   └── Navigation.tsx   # Mobile-friendly navigation
│   ├── services/            # API and utilities
│   │   └── api.ts          # Backend API integration
│   ├── types/              # TypeScript definitions
│   │   └── api.ts          # API response types
│   ├── App.tsx             # Main application component
│   ├── main.tsx            # Application entry point
│   └── index.css           # Global styles and Tailwind
├── public/                 # Static assets
├── package.json           # Dependencies and scripts
├── vite.config.ts         # Vite configuration
├── tailwind.config.js     # Tailwind CSS configuration
└── tsconfig.json          # TypeScript configuration
```

## 🎨 Design System

### Colors
- **Primary**: Blue gradient (#0ea5e9)
- **Secondary**: Purple gradient (#d946ef)  
- **Background**: Dark slate (#0f172a, #020617)
- **Glass**: Semi-transparent with backdrop blur

### Typography
- **Primary Font**: Inter (sans-serif)
- **Code Font**: JetBrains Mono (monospace)

### Components
- **Glass Morphism**: Translucent panels with blur effects
- **Smooth Animations**: Framer Motion powered transitions
- **Responsive Grid**: Mobile-first responsive layouts

## 🚀 Performance Features

- **Code Splitting**: Automatic route-based code splitting
- **Tree Shaking**: Eliminates unused code
- **Asset Optimization**: Optimized images and fonts
- **Bundle Analysis**: Built-in bundle size analysis
- **Service Worker**: PWA caching strategies

## 🔧 Configuration

### Environment Variables
Create `.env.local` for custom configuration:
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_TITLE=SSTorytime
```

### Vite Configuration
- **Proxy Setup**: API calls proxied to backend
- **PWA Plugin**: Service worker and manifest generation
- **TypeScript**: Full type checking during build

## 📱 Mobile Features

### Touch Interactions
- **Drag & Drop**: Move visualization nodes
- **Pinch to Zoom**: Scale visualization
- **Swipe Navigation**: Mobile menu interactions

### Responsive Design
- **Breakpoints**: Mobile, tablet, desktop optimized
- **Touch Targets**: 44px minimum touch targets
- **Readable Text**: Optimal font sizes for mobile

### PWA Capabilities
- **Installable**: Add to home screen
- **Offline Ready**: Service worker caching
- **Native Feel**: Splash screen and app icons

## 🎯 Future Enhancements

### Planned Features
- **Real-time Collaboration**: Multi-user exploration
- **Voice Search**: Speech-to-text search input
- **AR Visualization**: Augmented reality node graphs
- **Advanced Filtering**: Complex query builders
- **Export Options**: Save visualizations and results

### Performance Improvements
- **Virtual Scrolling**: Handle large datasets
- **WebGL Rendering**: Hardware-accelerated graphics
- **Edge Caching**: CDN-based content delivery
- **Bundle Optimization**: Further size reductions

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

Same license as the main SSTorytime project.

---

**SSTorytime Frontend** - Bringing knowledge exploration to the modern web! 🚀