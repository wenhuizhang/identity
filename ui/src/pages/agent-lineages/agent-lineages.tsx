import {BasePage} from '@/components/layout/base-page';
import PlaceholderPageContent from '@/components/ui/placeholder-page-content';

const AgentLineages: React.FC = () => {
  return (
    <BasePage
      title="Agent Lineages"
      description={
        <div className="space-y-4">
          <p>Create and manage Agent Lineages, and publish these to the Identity Network</p>
          <p>
            <b>Agents</b> can take many forms, including python packages and docker images. We call these different forms <b>Agent Artefacts</b>, and
            so an agent can be linked to multiple artefacts. Each of these artefacts can then also have different versions as the agent changes over
            time, which we call <b>Artefact Versions</b>.
          </p>
          <p>
            An <b>Agent Lineage</b> is the collection of all the Agent Artefacts and Artefact Versions that you&apos;ve created, and grows over time
            as you create and publish new artefacts and versions of your agent. If you are updating an existing agent, then either select that agent
            from the list below, or create a new agent!
          </p>
        </div>
      }
      useBreadcrumbs={false}
    >
      <PlaceholderPageContent />
    </BasePage>
  );
};

export default AgentLineages;
