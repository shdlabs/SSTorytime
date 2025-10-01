// API Response Types based on the SSTorytime API

export type ResponseType = 
  | 'Orbits'
  | 'PageMap' 
  | 'TOC'
  | 'Sequence'
  | 'PathSolve'
  | 'Arrows'
  | 'STAT'
  | 'GUIDANCE'
  | 'ERROR';

export interface BaseResponse {
  Response: ResponseType;
  Time?: string;
  Intent?: string;
  Ambient?: string;
}

export interface ErrorResponse extends BaseResponse {
  Response: 'ERROR';
  ErrorType: string;
  Message: string;
  Query: string;
  Suggestions: string[];
}

export interface GuidanceResponse extends BaseResponse {
  Response: 'GUIDANCE';
  Content: {
    title: string;
    message: string;
    suggestions: string[];
    examples: string[];
  };
}

// Node and coordinate types
export interface XYZ {
  X: number;
  Y: number;
  Z: number;
  R?: number;
  Lat?: number;
  Lon?: number;
}

export interface NodeOrbit {
  Text: string;
  Reln?: any;
  XYZ: XYZ;
}

export interface Orbits extends BaseResponse {
  Response: 'Orbits';
  Content: NodeOrbit[];
}

// Page/Notes types
export interface PageMapContent {
  Title: string;
  Notes: any[];
  Navigation?: any;
}

export interface PageMap extends BaseResponse {
  Response: 'PageMap';
  Content: PageMapContent;
}

// Table of Contents types
export interface ChapterContext {
  Text: string;
  Reln?: any;
  XYZ: XYZ;
}

export interface TOCItem {
  Chapter: string;
  XYZ: XYZ;
  Context: ChapterContext[];
  Single: ChapterContext[];
  Common?: any;
}

export interface TOC extends BaseResponse {
  Response: 'TOC';
  Content: TOCItem[];
}

// Sequence/Story types
export interface Sequence extends BaseResponse {
  Response: 'Sequence';
  Content: any[];
}

// Path solving types
export interface PathSolve extends BaseResponse {
  Response: 'PathSolve';
  Content: any[];
}

// Arrows/Relations types
export interface ArrowData {
  ArrPtr: number;
  ASTtype: number;
  Short: string;
  Long: string;
  InvPtr: number;
  ISTtype: number;
  InvS: string;
  InvL: string;
}

export interface Arrows extends BaseResponse {
  Response: 'Arrows';
  Content: ArrowData[];
}

// Statistics types
export interface STAT extends BaseResponse {
  Response: 'STAT';
  Content: any;
}

// Union type for all possible responses
export type APIResponse = 
  | Orbits 
  | PageMap 
  | TOC 
  | Sequence 
  | PathSolve 
  | Arrows 
  | STAT 
  | GuidanceResponse 
  | ErrorResponse;

// Search parameters
export interface SearchParams {
  name?: string;
  chapter?: string;
  context?: string;
  limit?: number;
  from?: string;
  to?: string;
  depth?: number;
  story?: boolean;
  sequence?: boolean;
  notes?: boolean;
  stats?: boolean;
  arrows?: string;
}

// UI State types
export interface SearchState {
  query: string;
  isLoading: boolean;
  response: APIResponse | null;
  error: string | null;
  suggestions: string[];
}

export interface VisualizationNode {
  id: string;
  text: string;
  x: number;
  y: number;
  z: number;
  radius: number;
  color: string;
  connections: string[];
}

export interface AppState {
  search: SearchState;
  currentView: 'search' | 'chapters' | 'stats' | 'help';
  visualization: {
    nodes: VisualizationNode[];
    selectedNode: string | null;
    zoomLevel: number;
    centerPosition: { x: number; y: number };
  };
  ui: {
    isMobile: boolean;
    sidebarOpen: boolean;
    theme: 'dark' | 'light';
  };
}