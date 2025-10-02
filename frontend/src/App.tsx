import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Toaster } from 'react-hot-toast';
import { PencilSquareIcon, HomeIcon } from '@heroicons/react/24/outline';
import SearchBar from './components/SearchBar';
import Navigation from './components/Navigation';
import Visualization from './components/Visualization';
import Text2N4L from './components/Text2N4L';
import { APIResponse } from './types/api';
import { searchAPI } from './services/api';
import './index.css';

type ViewType = 'search' | 'text2n4l';

function App() {
  const [searchResult, setSearchResult] = useState<APIResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isMobile, setIsMobile] = useState(false);
  const [currentView, setCurrentView] = useState<ViewType>('search');

  // Check if device is mobile
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };
    
    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  // Handle search results
  const handleSearchResult = (result: APIResponse) => {
    setSearchResult(result);
    setError(null);
  };

  // Handle search errors
  const handleSearchError = (errorMessage: string) => {
    setError(errorMessage);
    setSearchResult(null);
  };

  // Handle quick search from navigation
  const handleQuickSearch = async (query: string) => {
    if (!query) {
      setSearchResult(null);
      setError(null);
      return;
    }
    
    try {
      const result = await searchAPI({ name: query });
      handleSearchResult(result);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Search failed';
      handleSearchError(errorMessage);
    }
  };

  // Render result content based on type
  const renderResult = () => {
    if (!searchResult) return null;

    const { Response } = searchResult;

    // Handle different response types
    switch (Response) {
      case 'Orbits':
        return (
          <div className="space-y-6">
            <Visualization 
              data={searchResult} 
              onNodeClick={(searchQuery) => handleQuickSearch(searchQuery)}
            />
            <details className="glass rounded-xl p-4 border border-gray-700/50">
              <summary className="text-primary-400 font-medium cursor-pointer">
                Raw Data ({searchResult.Content?.length || 0} items)
              </summary>
              <pre className="text-gray-300 whitespace-pre-wrap font-mono text-sm mt-4 overflow-x-auto">
                {JSON.stringify(searchResult, null, 2)}
              </pre>
            </details>
          </div>
        );

      case 'ERROR':
        return (
          <div className="glass rounded-xl p-6 border border-red-500/30 bg-red-900/20">
            <div className="flex items-center space-x-3 mb-4">
              <div className="w-2 h-2 bg-red-500 rounded-full animate-pulse"></div>
              <h3 className="text-red-400 font-medium">
                {searchResult.ErrorType || 'Search Error'}
              </h3>
            </div>
            <p className="text-gray-300 mb-4">{searchResult.Message}</p>
            {searchResult.Suggestions && searchResult.Suggestions.length > 0 && (
              <div>
                <h4 className="text-red-300 font-medium mb-2">Suggestions:</h4>
                <ul className="space-y-1">
                  {searchResult.Suggestions.map((suggestion, index) => (
                    <li key={index} className="text-gray-400 text-sm">
                      • {suggestion}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        );

      case 'GUIDANCE':
        return (
          <div className="glass rounded-xl p-6 border border-primary-500/30 bg-primary-900/20">
            <div className="flex items-center space-x-3 mb-4">
              <div className="w-2 h-2 bg-primary-500 rounded-full animate-pulse"></div>
              <h3 className="text-primary-400 font-medium">
                {searchResult.Content?.title || 'Guidance'}
              </h3>
            </div>
            <p className="text-gray-300 mb-4">{searchResult.Content?.message}</p>
            
            {searchResult.Content?.suggestions && (
              <div className="mb-4">
                <h4 className="text-primary-300 font-medium mb-2">Suggestions:</h4>
                <ul className="space-y-1">
                  {searchResult.Content.suggestions.map((suggestion: string, index: number) => (
                    <li key={index} className="text-gray-400 text-sm">
                      • {suggestion}
                    </li>
                  ))}
                </ul>
              </div>
            )}
            
            {searchResult.Content?.examples && (
              <div>
                <h4 className="text-primary-300 font-medium mb-2">Examples:</h4>
                <ul className="space-y-1">
                  {searchResult.Content.examples.map((example: string, index: number) => (
                    <li key={index} className="text-gray-400 text-sm font-mono">
                      • {example}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        );

      default:
        return (
          <div className="glass rounded-xl p-6 border border-gray-700/50">
            <div className="flex items-center space-x-3 mb-4">
              <div className="w-2 h-2 bg-primary-500 rounded-full animate-pulse"></div>
              <h3 className="text-primary-400 font-medium">{Response}</h3>
            </div>
            <pre className="text-gray-300 whitespace-pre-wrap font-mono text-sm overflow-x-auto">
              {JSON.stringify(searchResult, null, 2)}
            </pre>
          </div>
        );
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background-dark via-background-darker to-background-dark">
      {/* Background Pattern */}
      <div className="fixed inset-0 opacity-20">
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,rgba(14,165,233,0.1),transparent)] animate-pulse-slow"></div>
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_80%_20%,rgba(217,70,239,0.05),transparent)]"></div>
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_20%_80%,rgba(14,165,233,0.05),transparent)]"></div>
      </div>

      {/* Main Content */}
      <div className="relative z-10">
        {/* Header */}
        <header className="px-4 py-6 md:py-8">
          <div className="max-w-7xl mx-auto">
            {/* Navigation and Title Row */}
            <div className="flex items-center justify-between mb-8">
              <Navigation onQuickSearch={handleQuickSearch} />
              
              <motion.div
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.8, ease: "easeOut" }}
                className="text-center flex-1"
              >
                <h1 className="text-3xl md:text-5xl font-bold text-white mb-2 text-shadow-lg">
                  <span className="bg-gradient-to-r from-primary-400 via-secondary-400 to-primary-400 bg-clip-text text-transparent">
                    SSTorytime
                  </span>
                </h1>
                <p className="text-sm md:text-lg text-gray-300 text-shadow">
                  Interactive knowledge exploration
                </p>
              </motion.div>
              
              <div className="w-24"> {/* Spacer for balance */}</div>
            </div>

            {/* Search Bar - only show in search view */}
            {currentView === 'search' && (
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.8, delay: 0.2, ease: "easeOut" }}
              >
                <SearchBar
                  onSearchResult={handleSearchResult}
                  onError={handleSearchError}
                  placeholder={isMobile ? "Search..." : "Search knowledge, ask questions, explore connections..."}
                />
              </motion.div>
            )}

            {/* Text2N4L Button - show when not in text2n4l view */}
            {currentView !== 'text2n4l' && (
              <div className="flex justify-center mt-4">
                <button
                  onClick={() => setCurrentView('text2n4l')}
                  className="glass px-6 py-3 rounded-xl border border-primary-500/50 text-primary-300 hover:bg-primary-900/30 transition-all duration-200 flex items-center space-x-2"
                >
                  <PencilSquareIcon className="w-5 h-5" />
                  <span>Convert Text to N4L</span>
                </button>
              </div>
            )}

            {/* Back to Search Button - show in text2n4l view */}
            {currentView === 'text2n4l' && (
              <div className="flex justify-center mb-4">
                <button
                  onClick={() => setCurrentView('search')}
                  className="glass px-6 py-3 rounded-xl border border-gray-700/50 text-white hover:bg-gray-700/50 transition-all duration-200 flex items-center space-x-2"
                >
                  <HomeIcon className="w-5 h-5" />
                  <span>Back to Search</span>
                </button>
              </div>
            )}
          </div>
        </header>

        {/* Results Area */}
        <main className="px-4 pb-8">
          <div className="max-w-7xl mx-auto">
            {currentView === 'text2n4l' ? (
              /* Show Text2N4L Component */
              <Text2N4L />
            ) : (
              /* Show Search Results */
              <AnimatePresence mode="wait">
              {error && (
                <motion.div
                  key="error"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                  className="glass rounded-xl p-6 border border-red-500/30 bg-red-900/20"
                >
                  <div className="flex items-center space-x-3">
                    <div className="w-2 h-2 bg-red-500 rounded-full animate-pulse"></div>
                    <h3 className="text-red-400 font-medium">Search Error</h3>
                  </div>
                  <p className="text-gray-300 mt-2">{error}</p>
                </motion.div>
              )}

              {searchResult && (
                <motion.div
                  key="result"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                  className="space-y-6"
                >
                  {/* Response Type Indicator */}
                  <div className="flex items-center space-x-4">
                    <div className="flex items-center space-x-2">
                      <div className="w-3 h-3 bg-primary-500 rounded-full animate-pulse"></div>
                      <span className="text-primary-400 font-medium">
                        {searchResult.Response}
                      </span>
                    </div>
                    {searchResult.Time && (
                      <span className="text-gray-500 text-sm">
                        {searchResult.Time}
                      </span>
                    )}
                  </div>

                  {/* Content Display */}
                  {renderResult()}
                </motion.div>
              )}

              {!searchResult && !error && (
                <motion.div
                  key="welcome"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="text-center py-16"
                >
                  <div className="glass rounded-xl p-8 max-w-2xl mx-auto border border-gray-700/30">
                    <motion.div
                      animate={{ rotate: [0, 360] }}
                      transition={{ duration: 20, repeat: Infinity, ease: "linear" }}
                      className="w-16 h-16 mx-auto mb-6 opacity-20"
                    >
                      <div className="w-full h-full rounded-full border-2 border-primary-500 border-t-transparent"></div>
                    </motion.div>
                    
                    <h2 className="text-2xl font-semibold text-white mb-4">
                      Welcome to SSTorytime
                    </h2>
                    
                    <p className="text-gray-400 mb-6">
                      Start exploring by searching for topics, asking questions, or using special commands:
                    </p>
                    
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-left">
                      <div className="space-y-2">
                        <h4 className="text-primary-400 font-medium">Quick Searches:</h4>
                        <ul className="text-sm text-gray-300 space-y-1">
                          <li><code className="bg-gray-800 px-2 py-1 rounded">moon</code> - Simple search</li>
                          <li><code className="bg-gray-800 px-2 py-1 rounded">recipe fish soup</code> - Multi-word</li>
                        </ul>
                      </div>
                      
                      <div className="space-y-2">
                        <h4 className="text-primary-400 font-medium">Special Commands:</h4>
                        <ul className="text-sm text-gray-300 space-y-1">
                          <li><code className="bg-gray-800 px-2 py-1 rounded">\help</code> - Get help</li>
                          <li><code className="bg-gray-800 px-2 py-1 rounded">\chapters</code> - Browse chapters</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </motion.div>
              )}
            </AnimatePresence>
            )}
          </div>
        </main>
      </div>

      {/* Toast Notifications */}
      <Toaster
        position={isMobile ? "bottom-center" : "top-right"}
        toastOptions={{
          duration: 4000,
          style: {
            background: 'rgba(31, 41, 55, 0.9)',
            color: '#fff',
            border: '1px solid rgba(55, 65, 81, 0.7)',
            backdropFilter: 'blur(8px)',
          },
        }}
      />
    </div>
  );
}

export default App;