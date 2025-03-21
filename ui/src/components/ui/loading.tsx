import {cn} from '@/lib/utils';
import {LoaderPinwheelIcon} from 'lucide-react';

export interface LoadingProps {
  style?: React.CSSProperties;
  classNameLoader?: string;
}

export const Loading = ({style, classNameLoader}: LoadingProps) => {
  return (
    <div
      className="layout flex flex-col justify-center items-center h-full w-full absolute top-[50%] left-[50%] z-[1000] -translate-x-[50%] -translate-y-[50%]"
      style={style}
    >
      <LoaderPinwheelIcon className={cn('w-20 h-20 animate-spin', classNameLoader)} />
    </div>
  );
};

export const LoaderRelative = ({className}: {className?: string}) => {
  return (
    <div className={cn('flex justify-center items-center w-full', className)}>
      <LoaderPinwheelIcon className={'w-10 h-10 animate-spin'} />
    </div>
  );
};
