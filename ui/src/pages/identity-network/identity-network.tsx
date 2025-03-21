import {CreateUpdateIdentityNetworkContent} from '@/components/indentity-network/create-update-indentity-network';
import {BasePage} from '@/components/layout/base-page';
import {EmptyWallet} from '@/components/shared/empty-wallet';
import {useStore} from '@/store';
import {useShallow} from 'zustand/react/shallow';

const IdentityNetwork: React.FC = () => {
  const {walletProvider} = useStore(
    useShallow((store) => ({
      walletProvider: store.walletProvider
    }))
  );

  return (
    <BasePage
      title="Identity Network"
      description="The agent passport data is stored in the decentralised ID Nodes"
      breadcrumbs={[{text: 'Identity Network'}]}
    >
      {walletProvider ? <CreateUpdateIdentityNetworkContent /> : <EmptyWallet />}
    </BasePage>
  );
};

export default IdentityNetwork;
