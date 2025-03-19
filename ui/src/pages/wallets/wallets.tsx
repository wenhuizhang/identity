import {BasePage} from '@/components/layout/base-page';
import PlaceholderPageContent from '@/components/ui/placeholder-page-content';

const Wallets: React.FC = () => {
  return (
    <BasePage title="Wallets" description="View your wallets." breadcrumbs={[{text: 'Wallets'}]}>
      <PlaceholderPageContent />
    </BasePage>
  );
};

export default Wallets;
