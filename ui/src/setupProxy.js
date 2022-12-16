// THIS IS DEVELOPMENT HELPER FILE! NOT PART OF THE WEB UI APP!

const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://localhost:55080',
      changeOrigin: true,
			pathRewrite: {
				'^/api/': '/' // rewrite path
			}
    })
  );
};
