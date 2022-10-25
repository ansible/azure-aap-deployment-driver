import * as React from 'react';
import { AppRouteConfig, IAppRouteGroup, routesConfig } from '../routesconfig';
import { RouteObject, useRoutes } from 'react-router-dom';


// a helper function to map config to what useRoutes expects
function getRoutes(config: AppRouteConfig[]) :RouteObject[] {
  return config.map((aRoute) => (
    ('routes' in aRoute)
      ? { children: [...getRoutes((aRoute as IAppRouteGroup).routes) ] }
      : { element: React.createElement(aRoute.component), path: aRoute.path }
  ))
}

function AppRoutes () {
  //let routesElement = useRoutes(getRoutes(routesConfig));
  return (
      <>
        {useRoutes(getRoutes(routesConfig))}
      </>
  )
}

export { AppRoutes };
