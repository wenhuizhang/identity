import {Outlet} from 'react-router-dom';
import {useEffect, useState} from 'react';
import {cn} from '@/lib/utils';
import {ResizablePanel, ResizablePanelGroup} from '@/components/ui/resizable';
import {SideNav} from './side-nav';

const Layout = () => {
  const defaultLayout = [15, 85];
  const defaultCollapsedLayout = [3.5, 96.5];
  const [layout, setLayout] = useState<number[]>(window.innerWidth < 768 ? defaultCollapsedLayout : defaultLayout);
  const [isCollapsed, setIsCollapsed] = useState(window.innerWidth < 768 ? true : false);

  const handleResize = () => {
    if (window.innerWidth < 768) {
      setIsCollapsed(true);
    } else {
      setIsCollapsed(false);
    }
  };

  useEffect(() => {
    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  const onLayout = (sizes: number[]) => {
    setLayout(sizes);
  };

  useEffect(() => {
    if (isCollapsed) {
      onLayout(defaultCollapsedLayout);
    } else {
      onLayout(defaultLayout);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isCollapsed]);

  return (
    <ResizablePanelGroup
      direction="horizontal"
      style={{
        height: '100vh',
        minWidth: '700px'
      }}
      className="fixed "
      onLayout={onLayout}
    >
      <ResizablePanel
        defaultSize={layout[0]}
        collapsedSize={defaultCollapsedLayout[0]}
        minSize={10}
        maxSize={15}
        collapsible={true}
        onCollapse={() => {
          setIsCollapsed(true);
        }}
        onExpand={() => {
          setIsCollapsed(false);
        }}
        className={cn(
          'transition-all duration-300 ease-in-out',
          isCollapsed && 'min-w-[3.5rem] max-w-[3.5rem]',
          !isCollapsed && 'min-w-[14.5rem] max-w-[14.5rem]'
        )}
      >
        <SideNav isCollapsed={isCollapsed} />
      </ResizablePanel>
      <ResizablePanel defaultSize={defaultLayout[1]} collapsible={false} minSize={30} className="layout !overflow-auto">
        <main
          className="pb-10 h-full overflow-y-auto"
          style={{
            minHeight: 'calc(100vh)'
          }}
        >
          <Outlet />
        </main>
      </ResizablePanel>
    </ResizablePanelGroup>
  );
};

export default Layout;
