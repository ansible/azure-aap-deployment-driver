import * as React from 'react';
import { Installer } from './Installer/Installer';
import { Documentation } from './Documentation/Documentation';

export interface IAppRoute {
  label?: string; // Excluding the label will exclude the route from the nav sidebar in AppLayout
  /* eslint-disable @typescript-eslint/no-explicit-any */
  component: React.FunctionComponent;
  /* eslint-enable @typescript-eslint/no-explicit-any */
  exact?: boolean;
  path: string;
  title: string;
  isAsync?: boolean;
  routes?: undefined;
}

export interface IAppRouteGroup {
  label: string;
  routes: IAppRoute[];
}

export type AppRouteConfig = IAppRoute | IAppRouteGroup;

const routesConfig: AppRouteConfig[] = [
  {
    component: Installer,
    path: '/',
    label: 'Installer',
    title: 'Installer'
  },
  {
    component: Documentation,
    path: '/documentation',
    label: 'Documentation',
    title: 'Documentation'
  }
];

export { routesConfig };
