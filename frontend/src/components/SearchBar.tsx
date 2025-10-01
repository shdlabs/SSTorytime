import React, { useState, useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { MagnifyingGlassIcon, XMarkIcon, CommandLineIcon } from '@heroicons/react/24/outline';
import { searchAPI, getSearchSuggestions } from '../services/api';
import { APIResponse, SearchParams } from '../types/api';

interface SearchBarProps {
  onSearchResult: (result: APIResponse) => void;
  onError: (error: string) => void;
  placeholder?: string;
  className?: string;
}

export const SearchBar: React.FC<SearchBarProps> = ({
  onSearchResult,
  onError,
  placeholder = "Search knowledge, ask questions...",
  className = ""
}) => {
  const [query, setQuery] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [suggestions, setSuggestions] = useState<string[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedSuggestion, setSelectedSuggestion] = useState(-1);

  // Update suggestions when query changes
  useEffect(() => {
    if (query.length > 0) {
      const newSuggestions = getSearchSuggestions(query);
      setSuggestions(newSuggestions);
      setShowSuggestions(newSuggestions.length > 0);
    } else {
      setSuggestions([]);
      setShowSuggestions(false);
    }
    setSelectedSuggestion(-1);
  }, [query]);

  // Handle search execution
  const performSearch = useCallback(async (searchQuery: string) => {
    if (!searchQuery.trim()) return;

    setIsLoading(true);
    setShowSuggestions(false);
    
    try {
      const result = await searchAPI({ name: searchQuery.trim() });
      onSearchResult(result);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Search failed';
      onError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  }, [onSearchResult, onError]);

  // Handle form submission
  const handleSubmit = useCallback((e: React.FormEvent) => {
    e.preventDefault();
    if (selectedSuggestion >= 0 && selectedSuggestion < suggestions.length) {
      const selectedQuery = suggestions[selectedSuggestion].split(' - ')[0];
      setQuery(selectedQuery);
      performSearch(selectedQuery);
    } else {
      performSearch(query);
    }
  }, [query, suggestions, selectedSuggestion, performSearch]);

  // Handle keyboard navigation
  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (!showSuggestions || suggestions.length === 0) return;

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedSuggestion(prev => 
          prev < suggestions.length - 1 ? prev + 1 : prev
        );
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedSuggestion(prev => prev > 0 ? prev - 1 : -1);
        break;
      case 'Escape':
        setShowSuggestions(false);
        setSelectedSuggestion(-1);
        break;
    }
  }, [showSuggestions, suggestions.length]);

  // Handle suggestion selection
  const selectSuggestion = useCallback((suggestion: string, index: number) => {
    const selectedQuery = suggestion.split(' - ')[0];
    setQuery(selectedQuery);
    setSelectedSuggestion(index);
    setShowSuggestions(false);
    performSearch(selectedQuery);
  }, [performSearch]);

  // Clear search
  const clearSearch = useCallback(() => {
    setQuery('');
    setSuggestions([]);
    setShowSuggestions(false);
    setSelectedSuggestion(-1);
  }, []);

  return (
    <div className={`relative w-full max-w-2xl mx-auto ${className}`}>
      {/* Search Form */}
      <form onSubmit={handleSubmit} className="relative">
        <div className="relative group">
          <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
            {isLoading ? (
              <motion.div
                animate={{ rotate: 360 }}
                transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                className="w-5 h-5 text-primary-500"
              >
                <CommandLineIcon />
              </motion.div>
            ) : (
              <MagnifyingGlassIcon className="w-5 h-5 text-gray-400 group-focus-within:text-primary-500 transition-colors" />
            )}
          </div>
          
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={handleKeyDown}
            onFocus={() => query.length > 0 && setShowSuggestions(true)}
            placeholder={placeholder}
            disabled={isLoading}
            className="input-field pl-12 pr-12 py-4 text-lg
                     bg-gray-900/50 backdrop-blur-sm border-gray-700/50
                     focus:bg-gray-900/80 focus:border-primary-500/50
                     disabled:opacity-50 disabled:cursor-not-allowed
                     transition-all duration-300 ease-in-out"
          />
          
          {query && (
            <button
              type="button"
              onClick={clearSearch}
              className="absolute inset-y-0 right-0 pr-4 flex items-center
                       text-gray-400 hover:text-white transition-colors"
            >
              <XMarkIcon className="w-5 h-5" />
            </button>
          )}
        </div>
      </form>

      {/* Search Suggestions */}
      <AnimatePresence>
        {showSuggestions && suggestions.length > 0 && (
          <motion.div
            initial={{ opacity: 0, y: -10, scale: 0.95 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -10, scale: 0.95 }}
            transition={{ duration: 0.2, ease: "easeOut" }}
            className="absolute top-full left-0 right-0 mt-2 
                     glass rounded-xl border border-gray-700/50 
                     shadow-2xl z-50 overflow-hidden"
          >
            <div className="py-2">
              {suggestions.map((suggestion, index) => {
                const [command, description] = suggestion.split(' - ');
                
                return (
                  <motion.button
                    key={index}
                    type="button"
                    onClick={() => selectSuggestion(suggestion, index)}
                    className={`search-suggestion w-full text-left px-4 py-3
                             ${selectedSuggestion === index ? 'bg-primary-600/20 text-primary-300' : 'text-gray-300'}
                             hover:bg-gray-700/50 transition-all duration-150`}
                    whileHover={{ x: 4 }}
                    whileTap={{ scale: 0.98 }}
                  >
                    <div className="flex items-center space-x-3">
                      <div className="flex-shrink-0">
                        {command.startsWith('\\') ? (
                          <CommandLineIcon className="w-4 h-4 text-primary-400" />
                        ) : (
                          <MagnifyingGlassIcon className="w-4 h-4 text-gray-400" />
                        )}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="font-mono text-sm text-white font-medium">
                          {command}
                        </div>
                        {description && (
                          <div className="text-xs text-gray-400 mt-1 truncate">
                            {description}
                          </div>
                        )}
                      </div>
                    </div>
                  </motion.button>
                );
              })}
            </div>
            
            {/* Quick tips */}
            <div className="border-t border-gray-700/50 px-4 py-3 bg-gray-800/50">
              <div className="text-xs text-gray-400 flex items-center space-x-4">
                <span className="flex items-center space-x-1">
                  <kbd className="px-1.5 py-0.5 text-xs font-mono bg-gray-700 rounded">↑↓</kbd>
                  <span>Navigate</span>
                </span>
                <span className="flex items-center space-x-1">
                  <kbd className="px-1.5 py-0.5 text-xs font-mono bg-gray-700 rounded">Enter</kbd>
                  <span>Search</span>
                </span>
                <span className="flex items-center space-x-1">
                  <kbd className="px-1.5 py-0.5 text-xs font-mono bg-gray-700 rounded">Esc</kbd>
                  <span>Close</span>
                </span>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default SearchBar;