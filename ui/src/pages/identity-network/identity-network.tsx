import {BasePage} from '@/components/layout/base-page';
import PlaceholderPageContent from '@/components/ui/placeholder-page-content';

const IdentityNetwork: React.FC = () => {
  return (
    <BasePage title="Identity Network" description="View your identity network." breadcrumbs={[{text: 'Identity Network'}]}>
      <PlaceholderPageContent />
    </BasePage>
  );
};

export default IdentityNetwork;
