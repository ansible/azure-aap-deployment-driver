import React from 'react';
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import '@patternfly/react-core/dist/styles/base.css';
import appRoutes from './app.routes';
import './app.css'

const App = () => (
  <RouterProvider router={createBrowserRouter(appRoutes)} />
)

export default App;
