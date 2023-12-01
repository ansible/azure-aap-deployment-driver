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
    path: "/deployment",
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
		element: <PlainLayout />,
		children: [
			{
				path: "/",
				element: <Login />
			}
		]
	},
	{
		element: <AppLayout navigation={appNavigation} />,
		children: [
			{
				path: "/welcome",
				element: <Deployment showLoginDialog={true} />
			}
		]
	},
	{
		element: <AppLayout navigation={appNavigation} />,
		children: [
			{
				path: "/deployment",
				element: <Deployment showLoginDialog={false}/>
			},
			{
				path: "/documentation",
				element: <Documentation />
			}
		]
	}

]


const App = () => (
  <RouterProvider router={createBrowserRouter(routes)} />
)

export default App;
