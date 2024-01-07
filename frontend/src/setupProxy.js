const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://backend:5000',
      changeOrigin: true,
      pathRewrite: {'^/api': ''},
      timeout: 200000, // Timeout in milliseconds
      proxyTimeout: 200000, // Timeout in milliseconds
    })
  );
};
