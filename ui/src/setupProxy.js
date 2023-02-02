// THIS IS DEVELOPMENT HELPER FILE! NOT PART OF THE WEB UI APP!

const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  const appFilesRegexp = /\.(json|ico|png|jpg|jpeg|txt|js|map|css|html)$/i;
  const cookieName = 'madd_session';
  const cookieValue = 'fake_cookie_for_development';
  const cookieOptions = {
    expires: 0,
    path: '/',
    httpOnly: true
  };
  app.use('/api/login', (req, res, next)=>{
    res.cookie(cookieName, cookieValue,cookieOptions)
    res.json({status: "success"})
  });
  app.use('/api/logout', (req, res, next)=>{
    res.clearCookie(cookieName)
    res.json({status: "success"})
  });
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://127.0.0.1:55080',
      changeOrigin: true,
			pathRewrite: {
				'^/api/': '/' // rewrite path
			}
    })
  );
  app.use(
    '/',
    (req, res, next) => {
      // keeping conditions for different paths separate to match what's configured in nginx
      if (req.path === '/login') {
        next();
      } else if (appFilesRegexp.test(req.path)) {
        next();
      } else {
        const cookies = req.headers['cookie']
        const hasSessionCookie = (cookies) ? cookies.split('; ').filter((aCookie)=>aCookie.startsWith(cookieName)).length: 0
        if (hasSessionCookie === 0) {
          res.redirect('/login')
        } else {
          next();
        }
      }
    }
  );
};
