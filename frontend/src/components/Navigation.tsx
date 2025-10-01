import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { 
  HomeIcon, 
  BookOpenIcon, 
  ChartBarIcon, 
  QuestionMarkCircleIcon,
  Bars3Icon,
  XMarkIcon 
} from '@heroicons/react/24/outline';

interface NavigationProps {
  onQuickSearch: (query: string) => void;
  className?: string;
}

interface NavItem {
  id: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  query: string;
  description: string;
}

const navItems: NavItem[] = [
  {
    id: 'home',
    label: 'Home',
    icon: HomeIcon,
    query: '',
    description: 'Return to main search'
  },
  {
    id: 'chapters',
    label: 'Chapters',
    icon: BookOpenIcon,
    query: '\\chapters',
    description: 'Browse all chapters'
  },
  {
    id: 'stats',
    label: 'Stats',
    icon: ChartBarIcon,
    query: '\\stats',
    description: 'View statistics'
  },
  {
    id: 'help',
    label: 'Help',
    icon: QuestionMarkCircleIcon,
    query: '\\help',
    description: 'Get search guidance'
  }
];

const quickSearches = [
  { label: 'Moon', query: 'moon', description: 'Astronomy and lunar topics' },
  { label: 'Recipes', query: 'recipe fish soup', description: 'Cooking and food topics' },
  { label: 'Brain', query: '\\notes brain', description: 'Neuroscience and cognition' },
  { label: 'Chinese', query: 'chinese directions', description: 'Language and culture' },
  { label: 'Stories', query: '\\story Mary', description: 'Narrative sequences' },
];

export const Navigation: React.FC<NavigationProps> = ({ onQuickSearch, className = "" }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [activeItem, setActiveItem] = useState<string | null>(null);

  const handleNavClick = (item: NavItem) => {
    setActiveItem(item.id);
    onQuickSearch(item.query);
    setIsOpen(false);
  };

  const handleQuickSearch = (query: string) => {
    onQuickSearch(query);
    setIsOpen(false);
  };

  return (
    <>
      {/* Mobile Menu Button */}
      <div className={`md:hidden ${className}`}>
        <button
          onClick={() => setIsOpen(true)}
          className="glass p-3 rounded-xl border border-gray-700/50 text-white hover:bg-gray-700/50 transition-colors"
        >
          <Bars3Icon className="w-6 h-6" />
        </button>
      </div>

      {/* Desktop Navigation */}
      <nav className={`hidden md:flex items-center space-x-2 ${className}`}>
        {navItems.map((item) => {
          const Icon = item.icon;
          const isActive = activeItem === item.id;
          
          return (
            <motion.button
              key={item.id}
              onClick={() => handleNavClick(item)}
              className={`flex items-center space-x-2 px-4 py-2 rounded-lg transition-all duration-200
                       ${isActive 
                         ? 'bg-primary-600 text-white' 
                         : 'glass text-gray-300 hover:text-white hover:bg-gray-700/50'
                       }`}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              title={item.description}
            >
              <Icon className="w-5 h-5" />
              <span className="text-sm font-medium">{item.label}</span>
            </motion.button>
          );
        })}
      </nav>

      {/* Mobile Navigation Overlay */}
      {isOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 z-50 md:hidden"
        >
          {/* Backdrop */}
          <div 
            className="absolute inset-0 bg-black/50 backdrop-blur-sm"
            onClick={() => setIsOpen(false)}
          />
          
          {/* Navigation Panel */}
          <motion.div
            initial={{ x: '-100%' }}
            animate={{ x: 0 }}
            exit={{ x: '-100%' }}
            transition={{ type: "spring", damping: 25, stiffness: 200 }}
            className="absolute left-0 top-0 bottom-0 w-80 max-w-[85vw] glass border-r border-gray-700/50 p-6"
          >
            {/* Header */}
            <div className="flex items-center justify-between mb-8">
              <h2 className="text-xl font-bold text-white">Navigation</h2>
              <button
                onClick={() => setIsOpen(false)}
                className="p-2 rounded-lg text-gray-400 hover:text-white hover:bg-gray-700/50 transition-colors"
              >
                <XMarkIcon className="w-6 h-6" />
              </button>
            </div>

            {/* Main Navigation */}
            <div className="space-y-3 mb-8">
              <h3 className="text-sm font-medium text-gray-400 uppercase tracking-wider">
                Main Sections
              </h3>
              {navItems.map((item) => {
                const Icon = item.icon;
                const isActive = activeItem === item.id;
                
                return (
                  <motion.button
                    key={item.id}
                    onClick={() => handleNavClick(item)}
                    className={`w-full flex items-center space-x-3 px-4 py-3 rounded-lg text-left transition-all duration-200
                             ${isActive 
                               ? 'bg-primary-600 text-white' 
                               : 'text-gray-300 hover:text-white hover:bg-gray-700/50'
                             }`}
                    whileHover={{ x: 4 }}
                    whileTap={{ scale: 0.98 }}
                  >
                    <Icon className="w-5 h-5" />
                    <div>
                      <div className="font-medium">{item.label}</div>
                      <div className="text-xs opacity-70">{item.description}</div>
                    </div>
                  </motion.button>
                );
              })}
            </div>

            {/* Quick Searches */}
            <div className="space-y-3">
              <h3 className="text-sm font-medium text-gray-400 uppercase tracking-wider">
                Quick Searches
              </h3>
              {quickSearches.map((search) => (
                <motion.button
                  key={search.query}
                  onClick={() => handleQuickSearch(search.query)}
                  className="w-full flex items-center justify-between px-4 py-3 rounded-lg text-left
                           text-gray-300 hover:text-white hover:bg-gray-700/50 transition-all duration-200"
                  whileHover={{ x: 4 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <div>
                    <div className="font-medium">{search.label}</div>
                    <div className="text-xs opacity-70">{search.description}</div>
                  </div>
                  <code className="text-xs bg-gray-800 px-2 py-1 rounded">
                    {search.query.length > 10 ? `${search.query.substring(0, 10)}...` : search.query}
                  </code>
                </motion.button>
              ))}
            </div>

            {/* Footer */}
            <div className="absolute bottom-6 left-6 right-6">
              <div className="text-center text-xs text-gray-500">
                SSTorytime v2.0 • Interactive Frontend
              </div>
            </div>
          </motion.div>
        </motion.div>
      )}
    </>
  );
};

export default Navigation;