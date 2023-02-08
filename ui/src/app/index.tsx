import React from 'react';
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import '@patternfly/react-core/dist/styles/base.css';
import { AppLayout } from './AppLayout/AppLayout';
import { PlainLayout } from './AppLayout/PlainLayout';
import { Deployment } from './Deployment/Deployment';
import { Documentation } from './Documentation/Documentation';
import { Login } from './Login/Login';
import './app.css'

// Navigation links are built from this array
const appNavigation = [
  {
    path: "/",
    label: "Deployment"
  },
  {
    path: "/documentation",
    label: "Documentation"
  },
  {
    path: "/logout",
    label: "Logout"
  }
]

const routes = [
	{
		element: <AppLayout navigation={appNavigation} />,
		children: [
			{
				path: "/",
				element: <Deployment />
			},
			{
				path: "documentation",
				element: <Documentation />
			}
		]
	},
	{
		element: <PlainLayout />,
		children: [
			{
				path: "/login",
				element: <Login />
			}
		]
	},

]


const App = () => (
  <RouterProvider router={createBrowserRouter(routes)} />
)

export default App;
