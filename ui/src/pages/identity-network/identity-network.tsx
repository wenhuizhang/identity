import {CreateUpdateIdentityNetworkContent} from '@/components/indentity-network/create-update-indentity-network';
import {BasePage} from '@/components/layout/base-page';

const IdentityNetwork: React.FC = () => {
  return (
    <BasePage
      title="Identity Network"
      description="The agent passport data is stored in the decentralised ID Nodes"
      breadcrumbs={[{text: 'Identity Network'}]}
    >
      <CreateUpdateIdentityNetworkContent />
    </BasePage>
  );
};

export default IdentityNetwork;
