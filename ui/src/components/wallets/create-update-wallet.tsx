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
import {LoaderRelative} from '../ui/loading';
import {Button} from '../ui/button';
import {Instructions} from '../ui/instructions';
import {SelectWallet} from './steps/select-wallet';
import {GenerateStoreKeys} from './steps/generate-store-keys';
import {validateForm} from '@/lib/utils';
import {WalletProviderFormValues, WalletProviderSchema} from '@/schemas/wallet-schemas';
import {useShallow} from 'zustand/react/shallow';
import {useStore} from '@/store';
import {WalletProviders} from '@/types/wallet-providers';
import {toast} from 'sonner';
import {useNavigate} from 'react-router-dom';
import {PATHS} from '@/router/paths';

export const CreateUpdateWalletContent: React.FC = () => {
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

  const {walletProvider, setWalletProvider} = useStore(
    useShallow((store) => ({
      walletProvider: store.walletProvider,
      setWalletProvider: store.setWalletProvider
    }))
  );

  const form = useForm<z.infer<typeof methods.current.schema>>({
    // @ts-expect-error zodResolver expects a zod schema
    resolver: zodResolver(methods.current.schema),
    mode: 'all'
  });

  const instructions = useMemo(() => {
    return [
      <div key={1}>
        The <strong>AGNTCY</strong> Agent Identity Management tool does not store or share any keys that are used to provide identity to your agents.
      </div>,
      <div key={2}>The tool connects to popoular password management applications or crypto wallets to handle the keys.</div>,
      <div key={3}>
        The keys are generated via this tool use quantum-resistant algorithms, and you can find more information on these{' '}
        <a href="#" target="_blank" className="inline-link">
          here
        </a>
        .
      </div>,
      <div key={4}>
        In order for others users to use and verify the indentity of agents you publish, you will have to also publish your{' '}
        <strong>public key</strong> in one of the supported <strong>Trust Anchors</strong>.
      </div>,
      <div key={5}>
        You can find out more about trust anchors and how to publish the public key{' '}
        <a href="#" target="_blank" className="inline-link">
          here
        </a>
        .
      </div>
    ];
  }, []);

  const handleSelectProvider = useCallback(() => {
    const values = form.getValues() as WalletProviderFormValues;
    const validationResult = validateForm(WalletProviderSchema, values);
    if (!validationResult.success) {
      validationResult.errors?.forEach((error) => {
        const fieldName = error.path[0] as keyof z.infer<typeof WalletProviderSchema>;
        form.setError(fieldName, {type: 'manual', ...error});
      });
      return;
    }
    methods.setMetadata('connectWallet', {...methods.getMetadata('connectWallet'), provider: values.provider});
    methods.next();
  }, [form, methods]);

  const handleSave = useCallback(() => {
    setIsLoading(true);
    const walletProvider = methods.getMetadata('connectWallet')?.provider as WalletProviders;
    setTimeout(() => {
      setWalletProvider(walletProvider);
      toast.success('Wallet provider connected successfully');
      void navigate(PATHS.identityNetwork, {replace: true});
      setIsLoading(false);
    }, 2500);
  }, [methods, navigate, setWalletProvider]);

  const onSubmit = () => {
    if (methods.current.id === 'connectWallet') {
      return handleSelectProvider();
    }
    if (methods.current.id === 'generateAndStoreKeys') {
      return handleSave();
    }
  };

  useEffect(() => {
    if (walletProvider) {
      methods.setMetadata('connectWallet', {...methods.getMetadata('connectWallet'), provider: walletProvider});
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [walletProvider]);

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
                {isLoading ? (
                  <ContainerStepper>
                    <div className="flex items-center	h-full">
                      <LoaderRelative />
                    </div>
                  </ContainerStepper>
                ) : (
                  methods.switch({
                    connectWallet: () => (
                      <ContainerStepper>
                        <SelectWallet />
                      </ContainerStepper>
                    ),
                    generateAndStoreKeys: () => (
                      <ContainerStepper>
                        <GenerateStoreKeys />
                      </ContainerStepper>
                    )
                  })
                )}
              </StepperPanel>
            </div>
            <StepperControls className="pt-4">
              <Button variant="secondary" onClick={methods.prev} disabled={methods.isFirst || isLoading} className="cursor-pointer">
                Previous
              </Button>
              <Button type="submit" disabled={isLoading || !form.formState.isValid} className="cursor-pointer">
                {methods.isLast ? 'Generate Quantum Safe Keys and add to Wallet' : 'Next'}
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
