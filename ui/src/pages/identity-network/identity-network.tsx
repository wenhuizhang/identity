import {BasePage} from '@/components/layout/base-page';
import PlaceholderPageContent from '@/components/ui/placeholder-page-content';

const IdentityNetwork: React.FC = () => {
  return (
    <BasePage title="Identity Network" description={'Agent passport data is stored in the decentralised ID Nodes'} useBreadcrumbs={false}>
      <PlaceholderPageContent />
    </BasePage>
  );
};

export default IdentityNetwork;
