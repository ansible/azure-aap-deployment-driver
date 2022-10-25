import * as React from 'react';
import { Dashboard } from './Dashboard/Dashboard';
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
    component: Dashboard,
    path: '/aap-dashboard',
    label: 'Dashboard',
    title: 'Dashboard'
  },
  {
    component: Support,
    path: '/aap-dashboard/support',
    label: 'Support',
    title: 'Support'
  },
  {
    label: 'Settings',
    routes: [
      {
        component: GeneralSettings,
        path: '/aap-dashboard/settings/general',
        label: 'General',
        title: 'General settings'
      },
      {
        component: ProfileSettings,
        path: '/aap-dashboard/settings/profile',
        label: 'Profile',
        title: 'Profile settings'
      },
    ],
  }
];

export { routesConfig };