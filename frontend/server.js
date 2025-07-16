const express = require('express');
const path = require('path');
const app = express();

console.log('Starting server...');
console.log('__dirname:', __dirname);
console.log('Static files from:', path.join(__dirname, 'dist/frontend/browser'));

// Serve static files from the dist directory
app.use(express.static(path.join(__dirname, 'dist/frontend/browser')));

// Handle Angular routing, return all requests to Angular app
app.get('*', (req, res) => {
  console.log('Request for:', req.url);
  res.sendFile(path.join(__dirname, 'dist/frontend/browser/index.html'));
});

const port = process.env.PORT || 10000;
app.listen(port, '0.0.0.0', () => {
  console.log(`Server running on port ${port}`);
});
