import {ReactNode} from 'react';
import {Card, CardDescription, CardHeader, CardTitle} from '../ui/card';
import {cn} from '@/lib/utils';
import {WalletProviders} from '@/types/wallet-providers';

export const WalletProvider = ({imgURI, walletName, walletDetails, isDisabled, isSelected, type, onSelect}: WalletProviderProps) => {
  return (
    <Card
      onClick={() => (!isDisabled ? onSelect?.(type) : undefined)}
      className={cn(
        'bg-[rgba(229,231,235,.2)] opacity-80 p-[0.75rem] text-[0.875rem] flex justify-between max-w-[225px] w-full flex-row items-center rounded-lg gap-[0.625rem] px-2 py-2 cursor-pointer hover:outline-1 hover:outline hover:opacity-100',
        isSelected && 'opacity-100 outline border-solid',
        isDisabled && 'opacity-50 cursor-no-drop',
        'min-w-[300px]',
        'card-flex-group'
      )}
    >
      <CardHeader className="px-2 py-0">
        <CardTitle>{walletName}</CardTitle>
        {walletDetails && <CardDescription>{walletDetails}</CardDescription>}
      </CardHeader>
      {imgURI}
    </Card>
  );
};

export interface WalletProviderProps {
  type: WalletProviders;
  imgURI?: ReactNode;
  walletName: string;
  walletDetails?: string;
  isDisabled?: boolean;
  isSelected?: boolean;
  onSelect?: (provider: WalletProviders) => void;
}
