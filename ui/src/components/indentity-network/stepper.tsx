import {defineStepper} from '@/components/ui/stepper';
import {ConnectIdentityNetworkSchema} from '@/schemas/identity-network-schema';
import {z} from 'zod';

export const {StepperProvider, StepperControls, StepperNavigation, StepperStep, StepperTitle, StepperDescription, StepperPanel, useStepper} =
  defineStepper(
    {
      id: 'connectIdentityNetwork',
      title: 'Connect Identity Network',
      description: 'Select the identity network you want to connect',
      schema: ConnectIdentityNetworkSchema
    },
    {
      id: 'verifyTheConnection',
      title: 'Verify the Connection',
      description: 'Verify the connection to the identity network',
      schema: z.object({})
    }
  );
