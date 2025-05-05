import {Route} from '@/types/router';
import {useCallback, useMemo} from 'react';
import {PATHS} from './paths';
import {NodeRoute} from '@/components/router/node-route';
import NotFound from '@/components/router/404';
import Layout from '@/components/layout/layout';
import {Navigate} from 'react-router-dom';
import React from 'react';

const WalletsKeys = React.lazy(() => import('@/pages/wallets-keys/wallets-keys'));
const IdentityNetwork = React.lazy(() => import('@/pages/identity-network/identity-network'));
const AgentLineages = React.lazy(() => import('@/pages/agent-lineages/agent-lineages'));
const VerifyAgentPassport = React.lazy(() => import('@/pages/verify-agent-passport/verify-agent-passport'));

export const generateRoutes = (routes: Route[]): Route[] => {
  return [
    {
      path: PATHS.basePath,
      element: (
        <NodeRoute>
          <Layout />
        </NodeRoute>
      ),
      children: [
        {
          index: true,
          element: <Navigate to={PATHS.walletsKeys} replace />
        },
        ...routes,
        {
          path: '*',
          element: <NotFound />
        }
      ]
    }
  ];
};

export const useRoutes = () => {
  const routes = useMemo<Route[]>(() => {
    return [
      {
        path: PATHS.walletsKeys,
        element: <WalletsKeys />
      },
      {
        path: PATHS.identityNetwork,
        element: <IdentityNetwork />
      },
      {
        path: PATHS.agentLineages,
        element: <AgentLineages />
      },
      {
        path: PATHS.verifyAgentPassport,
        element: <VerifyAgentPassport />
      }
    ];
  }, []);

  const removeDisabledRoutes = useCallback((routes: Route[]): Route[] => {
    return routes
      .filter((route) => !route.disabled)
      .map((route) => {
        if (route.children) {
          return {
            ...route,
            children: removeDisabledRoutes(route.children)
          };
        }
        return route;
      });
  }, []);

  const routesGenerated = generateRoutes(routes);

  return useMemo(() => {
    return removeDisabledRoutes(routesGenerated);
  }, [removeDisabledRoutes, routesGenerated]);
};
