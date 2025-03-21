import {defineStepper} from '@/components/ui/stepper';
import {GenerateStoreSchema, WalletProviderSchema} from '@/schemas/wallet-schemas';

export const {StepperProvider, StepperControls, StepperNavigation, StepperStep, StepperTitle, StepperDescription, StepperPanel, useStepper} =
  defineStepper(
    {
      id: 'connectWallet',
      title: 'Connect Wallet',
      description: 'Select the wallet you want to connect',
      schema: WalletProviderSchema
    },
    {
      id: 'generateAndStoreKeys',
      title: 'Generate and Store Keys',
      description: 'Generate and store the keys for the wallet',
      schema: GenerateStoreSchema
    }
  );
