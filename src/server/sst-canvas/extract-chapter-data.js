/**
 * SSTorytime Chapter Data Extraction (Node.js)
 * 
 * This script extracts JSON data from the SSTorytime API and saves it for visualization.
 * Run with: node extract-chapter-data.js
 */

const fs = require('fs');
const path = require('path');

async function extractChapterData()
{
    console.log('🔍 Extracting SSTorytime Chapter Data...');

    // Check if we have fetch available (Node 18+)
    let fetch;
    try
    {
        fetch = globalThis.fetch;
    } catch (e)
    {
        console.log('⚠️  Global fetch not available, trying node-fetch...');
        try
        {
            fetch = require('node-fetch');
        } catch (e)
        {
            console.error('❌ Neither global fetch nor node-fetch available.');
            console.log('Please install node-fetch: npm install node-fetch');
            console.log('Or use Node.js 18+ which has built-in fetch');
            process.exit(1);
        }
    }

    // Test server connection
    try
    {
        const testResponse = await fetch('http://localhost:8080/status');
        if (!testResponse.ok)
        {
            throw new Error(`Server responded with status ${testResponse.status}`);
        }
        console.log('✅ Server is running');
    } catch (error)
    {
        console.error('❌ Error: SSTorytime server is not running on localhost:8080');
        console.log('Please start the server first with: make run');
        process.exit(1);
    }

    // Different search queries to try
    const queries = [
        '\\chapter \\limit 25',
        'any \\chapter \\limit 25',
        '\\notes \\chapter \\limit 25',
        'story \\limit 25',
        'knowledge \\limit 25'
    ];

    for (let i = 0; i < queries.length; i++)
    {
        const query = queries[i];
        const filename = `chapter-data-${i}.json`;

        console.log(`🔍 Trying query: ${query}`);

        try
        {
            // Create form data
            const formData = new FormData();
            formData.append('search', query);

            // Make the request
            const response = await fetch('http://localhost:8080/searchN4L', {
                method: 'POST',
                body: formData
            });

            if (!response.ok)
            {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();

            // Save the data
            fs.writeFileSync(filename, JSON.stringify(data, null, 2));
            console.log(`✅ Saved to ${filename}`);

            // Extract stats
            const events = data.Content ? data.Content.length : 0;
            console.log(`   📊 Events found: ${events}`);

            // If this is the first successful extraction, also save as default
            if (i === 0)
            {
                fs.writeFileSync('chapter-data.json', JSON.stringify(data, null, 2));
                console.log('   💾 Also saved as chapter-data.json (default)');
            }

        } catch (error)
        {
            console.error(`❌ Failed to get response: ${error.message}`);
        }

        console.log('');
    }

    console.log('🎯 Data extraction complete!');
    console.log('');
    console.log('📁 Files created:');

    try
    {
        const files = fs.readdirSync('.')
            .filter(file => file.startsWith('chapter-data') && file.endsWith('.json'))
            .map(file =>
            {
                const stats = fs.statSync(file);
                return `   ${file} (${stats.size} bytes)`;
            });

        if (files.length > 0)
        {
            files.forEach(file => console.log(file));
        } else
        {
            console.log('   No JSON files created');
        }
    } catch (error)
    {
        console.log('   Unable to list files');
    }

    console.log('');
    console.log('🌐 To view the visualization:');
    console.log('   1. Open chapter-visualization.html in a web browser');
    console.log('   2. Or serve via HTTP: python3 -m http.server 8000');
    console.log('   3. Then visit: http://localhost:8000/chapter-visualization.html');
}

// Check if running directly
if (require.main === module)
{
    extractChapterData().catch(error =>
    {
        console.error('💥 Extraction failed:', error);
        process.exit(1);
    });
}

module.exports = { extractChapterData };