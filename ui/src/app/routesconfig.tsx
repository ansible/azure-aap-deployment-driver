import * as React from 'react';
import { Installer } from './Installer/Installer';
import { Support } from './Support/Support';
import { GeneralSettings } from './Settings/General/GeneralSettings';
import { ProfileSettings } from './Settings/Profile/ProfileSettings';

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

const routesConfig : AppRouteConfig[] = [
  {
    component: Installer,
    path: '/',
    label: 'Installer',
    title: 'Installer'
  },
  {
    component: Support,
    path: '/support',
    label: 'Support',
    title: 'Support'
  },
  {
    label: 'Settings',
    routes: [
      {
        component: GeneralSettings,
        path: '/settings/general',
        label: 'General',
        title: 'General settings'
      },
      {
        component: ProfileSettings,
        path: '/settings/profile',
        label: 'Profile',
        title: 'Profile settings'
      },
    ],
  }
];

export { routesConfig };
