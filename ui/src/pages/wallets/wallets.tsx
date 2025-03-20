import {BasePage} from '@/components/layout/base-page';
import {CreateUpdateWalletContent} from '@/components/wallets/create-update-wallet';

const Wallets: React.FC = () => {
  return (
    <BasePage title="Wallets" description="Explore the wallet connections" breadcrumbs={[{text: 'Wallets'}]}>
      <CreateUpdateWalletContent />
    </BasePage>
  );
};

export default Wallets;
