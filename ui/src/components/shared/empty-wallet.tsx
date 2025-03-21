import {PlusIcon, WalletMinimalIcon} from 'lucide-react';
import {useNavigate} from 'react-router-dom';
import {PATHS} from '@/router/paths';
import {Button} from '../ui/button';
import {Card} from '../ui/card';

export const EmptyWallet = () => {
  const navigate = useNavigate();

  const onClick = () => {
    void navigate(PATHS.wallets, {replace: true});
  };

  return (
    <Card>
      <div className="w-full h-full flex items-center	justify-center flex-col gap-6">
        <WalletMinimalIcon />
        <p className="italic text-[14px] text-white">Get started by adding a wallet</p>
        <Button size={'sm'} onClick={onClick} className="flex items-center justify-center gap-2">
          <PlusIcon className="h-4 w-4" />
          <p>Connect a Wallet</p>
        </Button>
      </div>
    </Card>
  );
};
