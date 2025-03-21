import {useCallback, useEffect, useMemo, useState} from 'react';
import {
  StepperControls,
  StepperDescription,
  StepperNavigation,
  StepperPanel,
  StepperProvider,
  StepperStep,
  StepperTitle,
  useStepper
} from './stepper';
import {useForm} from 'react-hook-form';
import {zodResolver} from '@hookform/resolvers/zod';
import {Card, CardContent} from '../ui/card';
import {Separator} from '../ui/separator';
import {z} from 'zod';
import {Form} from '../ui/form';
import {Button} from '../ui/button';
import {Instructions} from '../ui/instructions';
import {ConnectIdentiyURL} from './steps/connect-identity-url';
import {VerifyConnection} from './steps/verify-connection';
import {ConnectIdentityNetworFormValues, ConnectIdentityNetworkSchema} from '@/schemas/identity-network-schema';
import {validateForm} from '@/lib/utils';
import {useStore} from '@/store';
import {useShallow} from 'zustand/react/shallow';
import {toast} from 'sonner';
import {PATHS} from '@/router/paths';
import {useNavigate} from 'react-router-dom';

export const CreateUpdateIdentityNetworkContent: React.FC = () => {
  return (
    <StepperProvider variant="vertical" className="space-y-4">
      <FormStepperComponent />
    </StepperProvider>
  );
};

const FormStepperComponent: React.FC = () => {
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const methods = useStepper();
  const navigate = useNavigate();

  const {nodeUrl, setNodeUrl} = useStore(
    useShallow((store) => ({
      nodeUrl: store.nodeUrl,
      setNodeUrl: store.setNodeUrl
    }))
  );

  const form = useForm<z.infer<typeof methods.current.schema>>({
    resolver: zodResolver(methods.current.schema),
    mode: 'all'
  });

  const instructions = useMemo(() => {
    return [
      <div key={1}>
        Agent Identities are stored in the decentralised <strong>AGNCY</strong> Identity Network.
      </div>,
      <div key={2}>Identities stored in the network are synced every few minutes, and o connecting to any node will work.</div>,
      <div key={3}>
        You can find the available nodes and alternative networks{' '}
        <a href="#" target="_blank" className="inline-link">
          here
        </a>
        .
      </div>
    ];
  }, []);

  const handleConnect = useCallback(() => {
    const values = form.getValues() as ConnectIdentityNetworFormValues;
    const validationResult = validateForm(ConnectIdentityNetworkSchema, values);
    if (!validationResult.success) {
      validationResult.errors?.forEach((error) => {
        const fieldName = error.path[0] as keyof z.infer<typeof ConnectIdentityNetworkSchema>;
        form.setError(fieldName, {type: 'manual', ...error});
      });
      return;
    }
    methods.setMetadata('connectIdentityNetwork', {...methods.getMetadata('connectIdentityNetwork'), nodeUrl: values.nodeUrl});
    methods.next();
  }, [form, methods]);

  const handleVerify = useCallback(() => {
    setIsLoading(true);
    const nodeUrl = methods.getMetadata('connectIdentityNetwork')?.nodeUrl as string;
    setTimeout(() => {
      setNodeUrl(nodeUrl);
      toast.success('Successfully connected to the Identity Network');
      void navigate(PATHS.agentLineages, {replace: true});
      setIsLoading(false);
    }, 2500);
  }, [methods, navigate, setNodeUrl]);

  const onSubmit = () => {
    if (methods.current.id === 'connectIdentityNetwork') {
      return handleConnect();
    }
    if (methods.current.id === 'verifyTheConnection') {
      return handleVerify();
    }
  };

  useEffect(() => {
    if (nodeUrl) {
      methods.setMetadata('connectIdentityNetwork', {...methods.getMetadata('connectIdentityNetwork'), nodeUrl: nodeUrl});
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [nodeUrl]);

  return (
    <>
      <Card>
        <CardContent>{instructions && <Instructions instructions={instructions} />}</CardContent>
      </Card>
      <Card className="p-6">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="flex gap-4">
              <StepperNavigation className="min-w-[100px]">
                {methods.all.map((step) => (
                  <StepperStep key={step.id} of={step.id}>
                    <StepperTitle>{step.title}</StepperTitle>
                    <StepperDescription>{step.description}</StepperDescription>
                  </StepperStep>
                ))}
              </StepperNavigation>
              <StepperPanel className="w-full">
                {methods.switch({
                  connectIdentityNetwork: () => (
                    <ContainerStepper>
                      <ConnectIdentiyURL />
                    </ContainerStepper>
                  ),
                  verifyTheConnection: () => (
                    <ContainerStepper>
                      <VerifyConnection />
                    </ContainerStepper>
                  )
                })}
              </StepperPanel>
            </div>
            <StepperControls className="pt-4">
              <Button variant="secondary" onClick={methods.prev} disabled={methods.isFirst || isLoading} className="cursor-pointer">
                Previous
              </Button>
              <Button type="submit" disabled={isLoading || !form.formState.isValid} isLoading={isLoading} className="cursor-pointer">
                {methods.isLast ? 'Test the ID Node Connection' : 'Next'}
              </Button>
            </StepperControls>
          </form>
        </Form>
      </Card>
    </>
  );
};

const ContainerStepper: React.FC<React.PropsWithChildren> = ({children}) => {
  return (
    <div className="w-full flex h-full gap-4">
      <Separator orientation="vertical" className="h-full w-[1px]" />
      <div className="w-full">{children}</div>
    </div>
  );
};
