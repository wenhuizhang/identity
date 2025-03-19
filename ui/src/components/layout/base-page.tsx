import type {ReactNode} from 'react';
import React from 'react';
import {Tabs, TabsList, TabsTrigger} from '@/components/ui/tabs';
import Breadcrumbs, {BreadcrumbsProps} from '@/components/ui/breadcrumbs';
import {Link} from 'react-router-dom';

export const BasePage: React.FC<BasePageProps> = ({children, breadcrumbs, parentTitle, title, description, rightSideItems, subNav}) => {
  const hideHeader = !title && !description && !rightSideItems;
  const showHeader = !hideHeader;

  return (
    <>
      <div className="flex justify-between px-5 py-2 items-center bg-background-secondary dark:bg-background max-w-screen overflow-hidden border-b sticky top-0 z-40">
        <Breadcrumbs breadcrumbs={breadcrumbs} />
        <div className="h-[30px]" />
      </div>
      <div>
        {showHeader && (
          <div className="mt-4 flex items-center justify-between gap-2 py-2 flex-wrap mx-5 pb-2 mb-2">
            <div className="flex items-center justify-between border-b w-full flex-wrap pb-2 gap-2">
              <div>
                <h1 className="text-2xl flex items-center gap-2 mb-1 font-semibold">{parentTitle || title}</h1>
                <div className="text-muted-foreground min-h-4">{description}</div>
              </div>
              <div className="flex items-center gap-2">{rightSideItems}</div>
            </div>
          </div>
        )}
        <div className="md:px-5 py-3 bg-background">
          {subNav && (
            <Tabs
              className="-mt-2 mb-4"
              value={
                subNav.find((item) => {
                  return window.location.pathname === item.href;
                })?.href || subNav[0].href
              }
            >
              <TabsList>
                {subNav.map((item) => (
                  <Link to={item.href} key={`subNavItem-${item.href}`}>
                    <TabsTrigger value={item.href}>{item.label}</TabsTrigger>
                  </Link>
                ))}
              </TabsList>
            </Tabs>
          )}
          {children}
        </div>
      </div>
    </>
  );
};

interface BasePageProps {
  children: ReactNode;
  breadcrumbs?: BreadcrumbsProps['breadcrumbs'];
  parentTitle?: ReactNode;
  title?: ReactNode;
  description?: ReactNode;
  rightSideItems?: ReactNode;
  subNav?: {href: string; label: ReactNode; active?: boolean}[];
}
