import React from 'react';
import {Link} from 'react-router-dom';
import {Tooltip, TooltipContent, TooltipTrigger} from '../ui/tooltip';
import {Button} from '../ui/button';
import {TooltipArrow} from '@radix-ui/react-tooltip';
import Logo from '@/assets/logo-app-bar.svg';
import UnionLogo from '@/assets/union.svg?react';
import GitLogo from '@/assets/git.svg?react';
import {PATHS} from '@/router/paths';

export const AppBar: React.FC = () => {
  return (
    <header className="flex justify-between px-7 py-2 items-center max-w-screen overflow-hidden border-b sticky top-0 z-40 app-bar">
      <div className="flex gap-3 items-center">
        <Link to={PATHS.basePath}>
          <img src={Logo} alt="PyramdID" />
        </Link>
        <p className="product-name">Agent Identity</p>
      </div>
      <div className="flex items-center gap-1 flex-shrink-0">
        <Tooltip>
          <TooltipTrigger asChild>
            <Link to={'https://laughing-couscous-lrmjo5e.pages.github.io/'} target="_blank">
              <Button variant={'link'} size="icon" className="px-2 relative">
                <UnionLogo className="w-6 h-6" />
              </Button>
            </Link>
          </TooltipTrigger>
          <TooltipContent side="bottom">
            <TooltipArrow />
            Documentation
          </TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger asChild>
            <Link to={'https://github.com/agntcy/identity-spec'} target="_blank">
              <Button variant={'link'} size="icon" className="px-2 relative">
                <GitLogo className="w-6 h-6" />
              </Button>
            </Link>
          </TooltipTrigger>
          <TooltipContent side="bottom">
            <TooltipArrow />
            GitHub
          </TooltipContent>
        </Tooltip>
      </div>
    </header>
  );
};
