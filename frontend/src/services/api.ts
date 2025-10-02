import axios from 'axios';
import { APIResponse, SearchParams } from '../types/api';

// API Configuration
const API_BASE_URL = (import.meta as any).env.PROD ? '' : 'http://localhost:8080';
const SEARCH_ENDPOINT = '/searchN4L';

// Create axios instance with default config
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'multipart/form-data',
  },
});

// Request interceptor for logging
api.interceptors.request.use(
  (config) => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`, config.data);
    return config;
  },
  (error) => {
    console.error('API Request Error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    console.log('API Response:', response.data);
    return response;
  },
  (error) => {
    console.error('API Response Error:', error);
    
    // Handle specific error cases
    if (error.code === 'ECONNABORTED') {
      throw new Error('Request timeout - please try again');
    }
    
    if (error.response?.status === 404) {
      throw new Error('API endpoint not found');
    }
    
    if (error.response?.status >= 500) {
      throw new Error('Server error - please try again later');
    }
    
    throw new Error(error.message || 'Unknown API error');
  }
);

/**
 * Build search query string from parameters
 */
function buildSearchQuery(params: SearchParams): string {
  const parts: string[] = [];
  
  // Handle main search terms
  if (params.name) {
    parts.push(params.name);
  }
  
  // Handle special commands
  if (params.notes) {
    parts.unshift('\\notes');
  }
  
  if (params.story) {
    parts.unshift('\\story');
  }
  
  if (params.sequence) {
    parts.unshift('\\sequence');
  }
  
  if (params.stats) {
    parts.unshift('\\stats');
  }
  
  // Handle modifiers
  if (params.chapter) {
    parts.push(`\\chapter "${params.chapter}"`);
  }
  
  if (params.context) {
    parts.push(`\\context "${params.context}"`);
  }
  
  if (params.from) {
    parts.push(`\\from "${params.from}"`);
  }
  
  if (params.to) {
    parts.push(`\\to "${params.to}"`);
  }
  
  if (params.limit) {
    parts.push(`\\limit ${params.limit}`);
  }
  
  if (params.depth) {
    parts.push(`\\depth ${params.depth}`);
  }
  
  if (params.arrows) {
    parts.push(`\\arrow "${params.arrows}"`);
  }
  
  return parts.join(' ').trim();
}

/**
 * Perform search API call
 */
export async function searchAPI(params: SearchParams): Promise<APIResponse> {
  try {
    const query = buildSearchQuery(params);
    
    // Create FormData as expected by the API
    const formData = new FormData();
    formData.append('name', query);
    
    const response = await api.post<APIResponse>(SEARCH_ENDPOINT, formData);
    
    return response.data;
  } catch (error) {
    console.error('Search API Error:', error);
    throw error;
  }
}

/**
 * Quick search for simple queries
 */
export async function quickSearch(query: string): Promise<APIResponse> {
  return searchAPI({ name: query });
}

/**
 * Get table of contents
 */
export async function getTableOfContents(limit?: number): Promise<APIResponse> {
  const params: SearchParams = { name: '\\chapters' };
  if (limit) {
    params.limit = limit;
  }
  return searchAPI(params);
}

/**
 * Get help information
 */
export async function getHelp(): Promise<APIResponse> {
  return searchAPI({ name: '\\help' });
}

/**
 * Get statistics
 */
export async function getStats(): Promise<APIResponse> {
  return searchAPI({ stats: true });
}

/**
 * Get notes for a topic
 */
export async function getNotes(topic: string, chapter?: string): Promise<APIResponse> {
  const params: SearchParams = { 
    name: topic, 
    notes: true 
  };
  
  if (chapter) {
    params.chapter = chapter;
  }
  
  return searchAPI(params);
}

/**
 * Search with context
 */
export async function searchWithContext(
  query: string, 
  context: string, 
  limit?: number
): Promise<APIResponse> {
  return searchAPI({
    name: query,
    context,
    limit
  });
}

/**
 * Find path between two points
 */
export async function findPath(
  from: string, 
  to: string, 
  depth?: number
): Promise<APIResponse> {
  return searchAPI({
    from,
    to,
    depth
  });
}

/**
 * Get arrows/relationships
 */
export async function getArrows(pattern?: string): Promise<APIResponse> {
  return searchAPI({
    arrows: pattern || '1'
  });
}

/**
 * Get story/sequence
 */
export async function getStory(query: string, chapter?: string): Promise<APIResponse> {
  const params: SearchParams = {
    name: query,
    story: true
  };
  
  if (chapter) {
    params.chapter = chapter;
  }
  
  return searchAPI(params);
}

/**
 * Search suggestions based on common patterns
 */
export function getSearchSuggestions(partial: string): string[] {
  const suggestions: string[] = [];
  
  if (!partial.trim()) {
    return [
      '\\help - Get search guidance',
      '\\chapters - Browse all chapters',
      '\\stats - View statistics',
      'moon - Search for "moon"',
      'recipe fish soup - Multi-word search',
    ];
  }
  
  const lower = partial.toLowerCase();
  
  // Command suggestions
  if (lower.startsWith('\\')) {
    const commands = [
      '\\notes topic - Get notes for a topic',
      '\\chapter "chapter name" - Search within chapter',
      '\\context "context" - Search with context',
      '\\story topic - Get story sequence',
      '\\from start \\to end - Find path',
      '\\limit 5 - Limit results',
      '\\arrow pattern - Find relationships',
    ];
    
    suggestions.push(...commands.filter(cmd => 
      cmd.toLowerCase().includes(lower)
    ));
  } else {
    // Topic suggestions based on common examples
    const topics = [
      'moon astronomy',
      'recipe fish soup',
      'brain neuroscience',
      'chinese directions',
      'mathematics integration',
      'Mary story',
    ];
    
    suggestions.push(...topics.filter(topic => 
      topic.toLowerCase().includes(lower)
    ));
  }
  
  return suggestions.slice(0, 5);
}

export default api;