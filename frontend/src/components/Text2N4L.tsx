import React, { useState } from 'react';

interface ProcessResult {
  n4l_content: string;
  stats: {
    total_sentences: number;
    selected_sentences: number;
    final_fraction: number;
    requested_fraction: number;
  };
}

const Text2N4L: React.FC = () => {
  const [inputText, setInputText] = useState('');
  const [outputText, setOutputText] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);
  const [percentage, setPercentage] = useState(10);
  const [error, setError] = useState('');
  const [stats, setStats] = useState<ProcessResult['stats'] | null>(null);

  const processText = async () => {
    if (!inputText.trim()) {
      setError('Please enter some text to process');
      return;
    }

    setIsProcessing(true);
    setError('');
    
    try {
      const response = await fetch('http://localhost:3001/process', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          text: inputText,
          percentage: percentage,
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result: ProcessResult = await response.json();
      setOutputText(result.n4l_content);
      setStats(result.stats);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setIsProcessing(false);
    }
  };

  const clearAll = () => {
    setInputText('');
    setOutputText('');
    setStats(null);
    setError('');
  };

  const copyToClipboard = () => {
    navigator.clipboard.writeText(outputText);
  };

  return (
    <div className="min-h-screen bg-gray-50 p-4">
      <div className="max-w-7xl mx-auto">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Text to N4L Converter
          </h1>
          <p className="text-gray-600">
            Extract high-intentionality sentences from your text and convert them to N4L format for Promise Theory analysis.
          </p>
        </div>

        {/* Controls */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-6">
          <div className="flex flex-wrap items-center gap-4">
            <div className="flex items-center gap-2">
              <label htmlFor="percentage" className="text-sm font-medium text-gray-700">
                Selection Percentage:
              </label>
              <input
                id="percentage"
                type="number"
                min="1"
                max="100"
                value={percentage}
                onChange={(e) => setPercentage(Number(e.target.value))}
                className="w-20 px-2 py-1 border border-gray-300 rounded-md text-sm"
              />
              <span className="text-sm text-gray-500">%</span>
            </div>
            <button
              onClick={processText}
              disabled={isProcessing || !inputText.trim()}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-sm font-medium"
            >
              {isProcessing ? 'Processing...' : 'Process Text'}
            </button>
            <button
              onClick={clearAll}
              className="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 text-sm font-medium"
            >
              Clear All
            </button>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
            <div className="flex">
              <div className="text-red-600">
                <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                </svg>
              </div>
              <div className="ml-3">
                <p className="text-sm text-red-600">{error}</p>
              </div>
            </div>
          </div>
        )}

        {/* Stats */}
        {stats && (
          <div className="bg-blue-50 border border-blue-200 rounded-md p-4 mb-6">
            <h3 className="text-sm font-medium text-blue-900 mb-2">Processing Results</h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
              <div>
                <span className="text-blue-700 font-medium">Total Sentences:</span>
                <span className="ml-2 text-blue-900">{stats.total_sentences}</span>
              </div>
              <div>
                <span className="text-blue-700 font-medium">Selected:</span>
                <span className="ml-2 text-blue-900">{stats.selected_sentences}</span>
              </div>
              <div>
                <span className="text-blue-700 font-medium">Requested:</span>
                <span className="ml-2 text-blue-900">{stats.requested_fraction.toFixed(2)}%</span>
              </div>
              <div>
                <span className="text-blue-700 font-medium">Actual:</span>
                <span className="ml-2 text-blue-900">{stats.final_fraction.toFixed(2)}%</span>
              </div>
            </div>
          </div>
        )}

        {/* Main Content Area */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Input Area */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200">
            <div className="px-4 py-3 border-b border-gray-200">
              <h2 className="text-lg font-semibold text-gray-900">Input Text</h2>
              <p className="text-sm text-gray-600">Paste your text here for analysis</p>
            </div>
            <div className="p-4">
              <textarea
                value={inputText}
                onChange={(e) => setInputText(e.target.value)}
                placeholder="Paste your text here..."
                className="w-full h-96 p-3 border border-gray-300 rounded-md resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
              />
              <div className="mt-2 text-sm text-gray-500">
                Characters: {inputText.length}
              </div>
            </div>
          </div>

          {/* Output Area */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200">
            <div className="px-4 py-3 border-b border-gray-200 flex justify-between items-center">
              <div>
                <h2 className="text-lg font-semibold text-gray-900">N4L Output</h2>
                <p className="text-sm text-gray-600">Edit the generated N4L content as needed</p>
              </div>
              {outputText && (
                <button
                  onClick={copyToClipboard}
                  className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 text-sm"
                >
                  Copy
                </button>
              )}
            </div>
            <div className="p-4">
              <textarea
                value={outputText}
                onChange={(e) => setOutputText(e.target.value)}
                placeholder="Processed N4L content will appear here..."
                className="w-full h-96 p-3 border border-gray-300 rounded-md resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
                readOnly={!outputText}
              />
              {outputText && (
                <div className="mt-2 text-sm text-gray-500">
                  Lines: {outputText.split('\n').length}
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Usage Instructions */}
        <div className="mt-8 bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-3">How to Use</h3>
          <div className="prose prose-sm text-gray-600">
            <ol className="list-decimal list-inside space-y-2">
              <li>Paste your text in the left panel (Input Text area)</li>
              <li>Adjust the selection percentage if needed (default: 10%)</li>
              <li>Click "Process Text" to analyze the content</li>
              <li>Review the generated N4L content in the right panel</li>
              <li>Edit the N4L output as needed for your Promise Theory analysis</li>
              <li>Use the "Copy" button to copy the final result</li>
            </ol>
            <div className="mt-4">
              <p className="font-medium">About N4L Format:</p>
              <p>N4L (Narrative for Learning) format extracts high-intentionality sentences that are significant for knowledge representation and Promise Theory analysis. The algorithm identifies sentences with both dynamic running assessment and static post-hoc assessment techniques.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Text2N4L;