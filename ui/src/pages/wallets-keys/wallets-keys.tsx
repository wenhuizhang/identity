import {BasePage} from '@/components/layout/base-page';
import PlaceholderPageContent from '@/components/ui/placeholder-page-content';
import {Link} from 'react-router-dom';

const WalletsKeys: React.FC = () => {
  return (
    <BasePage
      title="Wallets and Keys"
      description={
        <div className="space-y-4 max-w-[70%]">
          <p>The keys used to sign your agent are stored in your local wallet.</p>
          <p>
            In order for other users to use and verify the identity of agents you publish, you will have to also publish your <b>public key</b> in one
            of the supported <b>Trust Anchors</b>. You can find out more about the trust anchors and how to publish the public key{' '}
            <Link to="#" className="inline-link">
              here
            </Link>
            .
          </p>
        </div>
      }
      useBreadcrumbs={false}
    >
      <PlaceholderPageContent />
    </BasePage>
  );
};

export default WalletsKeys;
