import {BasePage} from '@/components/layout/base-page';
import PlaceholderPageContent from '@/components/ui/placeholder-page-content';

const Dashboard: React.FC = () => {
  return (
    <BasePage title="Dashboard" description="View your dashboard.">
      <PlaceholderPageContent />
    </BasePage>
  );
};

export default Dashboard;
