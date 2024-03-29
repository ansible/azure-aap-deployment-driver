import React from 'react';
import { PlainLayout } from './AppLayout/PlainLayout';
import { AppLayout } from './AppLayout/AppLayout';
import { Login } from './Login/Login';
import { Deployment } from './Deployment/Deployment';
import { Documentation } from './Documentation/Documentation';
import appNavigation from './app.navigation';

const appRoutes = [
	{
		element: <PlainLayout />,
		children: [
			{
				path: "/login",
				element: <Login />
			}
		]
	},
	{
		element: <AppLayout navigation={appNavigation} />,
		children: [
			{
				path: "/rhlogin",
				element: <Deployment showLoginDialog={true} />
			}
		]
	},
	{
		element: <AppLayout navigation={appNavigation} />,
		children: [
			{
				path: "/",
				element: <Deployment showLoginDialog={false}/>
			},
			{
				path: "/documentation",
				element: <Documentation />
			}
		]
	}
]

export default appRoutes;
