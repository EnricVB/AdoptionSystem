const express = require('express');
const path = require('path');
const app = express();

// Serve static files from the dist directory
app.use(express.static(path.join(__dirname, 'dist/frontend/browser')));

// Handle Angular routing, return all requests to Angular app
app.get('*', (req, res) => {
  res.sendFile(path.join(__dirname, 'dist/frontend/browser/index.html'));
});

const port = process.env.PORT || 3000;
app.listen(port, () => {
  console.log(`Server running on port ${port}`);
});
