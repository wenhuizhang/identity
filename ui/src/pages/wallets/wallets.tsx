import {BasePage} from '@/components/layout/base-page';
import {CreateUpdateWalletContent} from '@/components/wallets/create-update-wallet';

const Wallets: React.FC = () => {
  return (
    <BasePage title="Wallets" description="The keys used to sign your agent are stored in your local wallet." breadcrumbs={[{text: 'Wallets'}]}>
      <CreateUpdateWalletContent />
    </BasePage>
  );
};

export default Wallets;
