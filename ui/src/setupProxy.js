// THIS IS DEVELOPMENT HELPER FILE! NOT PART OF THE WEB UI APP!
const express = require('express')
const { createProxyMiddleware, fixRequestBody} = require('http-proxy-middleware');


module.exports = function(app) {
  const appFilesRegexp = /\.(json|ico|png|jpg|jpeg|svg|txt|js|map|css|html|woff|woff2)$/i;
  const cookieName = 'madd_session'
  const cookieValue = 'fake_cookie_for_development'
  const cookieOptions = {
    expires: 0,
    path: '/',
    httpOnly: true
  }
  // this allows access to parsed request body
  app.use(express.json())

  app.use('/api/login', (req, res, next)=>{
    res.cookie(cookieName, cookieValue,cookieOptions)
    const {uid,pwd} = req.body
    const response = ((uid && uid === 'admin') && (pwd && typeof pwd === 'string' && pwd.length >= 12 )) ?
      {status: "success"} :
      {
        error: "Login Required",
        "status": 401
      }
    res.json(response)
  })
  app.use('/api/logout', (req, res, next)=>{
    res.clearCookie(cookieName)
    res.json({status: "success"})
  })
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://127.0.0.1:55080',
      changeOrigin: true,
      onProxyReq: fixRequestBody,
			pathRewrite: {
				'^/api/': '/' // rewrite path
			}
    })
  );
  app.use(
    '/',
    (req, res, next) => {
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
