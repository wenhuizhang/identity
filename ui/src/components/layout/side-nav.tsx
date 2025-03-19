import {ReactNode, useMemo} from 'react';
import {Link, useLocation} from 'react-router-dom';
import {Tooltip, TooltipContent, TooltipTrigger} from '@/components/ui/tooltip';
import {BookIcon, FileLock2Icon, GitForkIcon, WalletMinimalIcon} from 'lucide-react';
import {cn} from '@/lib/utils';
import {PATHS} from '@/router/paths';
import AgntcyLogo from '@/assets/agntcy-logo.svg?react';
import SmallAgntcyLogo from '@/assets/favicon.svg?react';
import {GlobeIcon} from '@radix-ui/react-icons';

export const SideNav: React.FC<{isCollapsed?: boolean}> = ({isCollapsed}) => {
  const sideNavLinks: {
    href: string;
    label: ReactNode;
    icon: ReactNode;
    description: ReactNode;
    onClick?: () => void;
  }[] = useMemo(() => {
    return [
      {
        href: PATHS.wallets,
        label: 'Wallets',
        icon: <WalletMinimalIcon className="w-4 h-4" />,
        description: 'View your wallets.'
      },
      {
        href: PATHS.identityNetwork,
        label: 'Identity Network',
        icon: <GlobeIcon className="w-4 h-4" />,
        description: 'View your identity network.'
      },
      {
        href: PATHS.agentLineages,
        label: 'Agent Lineages',
        icon: <GitForkIcon className="w-4 h-4" />,
        description: 'View your agent lineages.'
      },
      {
        href: PATHS.verifyAgentPassport,
        label: 'Verify Agent Passport',
        icon: <FileLock2Icon className="w-4 h-4" />,
        description: 'View your agent passport.'
      }
    ];
  }, []);

  const location = useLocation();
  const currentPathName = location.pathname;

  const active = sideNavLinks.find((link) => {
    return currentPathName.startsWith(link.href);
  });

  return (
    <nav className="flex relative flex-col justify-between gap-1 px-2 group-[[data-collapsed=true]]:justify-center group-[[data-collapsed=true]]:px-2 bg-side-nav-background h-full text-white">
      <div className="flex flex-col gap-1">
        <Link to={PATHS.basePath}>
          <button className={cn('flex justify-center items-center py-6 rounded w-full px-3', isCollapsed && 'px-1')}>
            {!isCollapsed ? (
              <div className="space-y-2">
                <AgntcyLogo className="h-12 w-full" />
                <span className="text-[16px] font-bold text-center text-[#187ADC]">Agent Identity Managment</span>
              </div>
            ) : (
              <SmallAgntcyLogo className="h-full w-12" />
            )}
          </button>
        </Link>
        <div className="flex flex-col gap-1">
          {sideNavLinks.map((link) => {
            return <SideNavLink {...link} isCollapsed={isCollapsed} key={`side-nav-link-${link.href}`} isActive={active?.href === link.href} />;
          })}
        </div>
      </div>
      <div className="mb-4">
        <SideNavLink
          label="Documentation"
          href={'#'}
          icon={<BookIcon className={cn('mr-1', isCollapsed ? '!w-4 !min-w-4 !h-4 !min-h-4' : '!w-4 !min-w-4 !h-4 !min-h-4')} />}
          description="View the documentation."
          isCollapsed={isCollapsed}
          isActive={false}
          isExternal={true}
        />
      </div>
    </nav>
  );
};

const SideNavLink: React.FC<{
  label: ReactNode;
  href: string;
  icon: ReactNode;
  isActive?: boolean;
  isCollapsed?: boolean;
  isExternal?: boolean;
  description: ReactNode;
  onClick?: () => void;
}> = ({label, href, icon, isActive, isCollapsed, description, isExternal, onClick}) => {
  const ThisLink = (
    <Link to={href} target={isExternal ? '_blank' : undefined} rel={isExternal ? 'noopener noreferrer' : undefined}>
      <button
        className={cn(
          'flex items-center px-3 py-3 gap-4 hover:bg-side-nav-hover w-full rounded-md text-white overflow-hidden text-sm cursor-pointer font-medium transition-colors',
          isActive && 'bg-side-nav-selected hover:bg-side-nav-selected-hover',
          isCollapsed && 'justify-center',
          !isCollapsed && 'pl-4'
        )}
        onClick={onClick}
      >
        <div
          className={cn(
            'object-cover flex-shrink-0 flex-grow-0 [&>svg]:transition-all',
            isCollapsed
              ? '[&>svg]:min-w-4 [&>svg]:min-h-4 [&>svg]:max-w-4 [&>svg]:max-h-4 stroke-[1.3]'
              : '[&>svg]:min-w-5 [&>svg]:min-h-5 [&>svg]:max-w-5 [&>svg]:max-h-5 stroke-[1.4]',
            isActive && '[&>svg]:text-[#F5A035]'
          )}
        >
          {icon}
        </div>
        {!isCollapsed && (
          <span
            className={cn(
              'text-white text-left whitespace-nowrap overflow-ellipsis overflow-hidden text-[0.9rem] pb-[0.125rem] font-semibold',
              isActive && 'font-bold text-[#F5A035]'
            )}
          >
            {label}
          </span>
        )}
      </button>
    </Link>
  );

  return isCollapsed ? (
    <SideNavTooltip label={label} description={description}>
      {ThisLink}
    </SideNavTooltip>
  ) : (
    ThisLink
  );
};

const SideNavTooltip: React.FC<{
  children: ReactNode;
  label: ReactNode;
  description: ReactNode;
}> = ({children, label, description}) => {
  return (
    <Tooltip>
      <TooltipTrigger asChild>{children}</TooltipTrigger>
      <TooltipContent side="right" sideOffset={12} className="font-semibold" align="start">
        <div className="flex flex-col gap-1 p-1 text-sm">
          <div className="font-semibold">{label}</div>
          <div className="text-muted-foreground font-normal">{description}</div>
        </div>
      </TooltipContent>
    </Tooltip>
  );
};
