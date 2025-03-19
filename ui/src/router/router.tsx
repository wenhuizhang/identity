import {useMemo} from 'react';
import {createBrowserRouter, RouterProvider} from 'react-router-dom';
import {useRoutes} from './routes';

export const Router = () => {
  const routes = useRoutes();
  const router = useMemo(() => createBrowserRouter(routes), [routes]);
  return <RouterProvider router={router} />;
};
