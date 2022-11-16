import React from 'react';
import { BrowserRouter } from "react-router-dom";
import '@patternfly/react-core/dist/styles/base.css';
import { AppLayout } from './AppLayout/AppLayout';
import { AppRoutes } from './AppRoutes/index';
import './app.css'

const App: React.FunctionComponent = () => (
  <BrowserRouter>
    <AppLayout>
      <AppRoutes/>
    </AppLayout>
  </BrowserRouter>
);

export default App;
