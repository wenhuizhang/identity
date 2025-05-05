import {BasePage} from '@/components/layout/base-page';
import PlaceholderPageContent from '@/components/ui/placeholder-page-content';
import {Link} from 'react-router-dom';

const VerifyAgentPassport: React.FC = () => {
  return (
    <BasePage
      title="Verify Agent Passport"
      description={
        <div className="space-y-4">
          <p>You can copy an Agent Passport into the field below, and verify that is legitimate and linked to a trust anchor. </p>
          <p>
            You can read more about the Internet of Agents trust anchors{' '}
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

export default VerifyAgentPassport;
