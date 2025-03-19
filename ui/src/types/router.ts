import {RouteObject} from 'react-router-dom';

export interface CustomRoute {
  disabled?: boolean;
}

export type Route = RouteObject & CustomRoute;
