import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {useStepper} from '../stepper';
import {useFormContext} from 'react-hook-form';
import {WalletProvider, WalletProviderProps} from '@/components/shared/wallet-provider';
import {WalletProviders} from '@/types/wallet-providers';
import {WalletProviderFormValues} from '@/schemas/wallet-schemas';
import {FormDescription, FormField, FormItem, FormMessage} from '@/components/ui/form';
import PasswordLogo from '@/assets/1password-logo.svg?react';
import BitWardenLogo from '@/assets/bitwarden-logo.svg?react';
import DropboxLogo from '@/assets/dropbox-logo.svg?react';
import ProtonLogo from '@/assets/proton-logo.svg?react';
import DashLogo from '@/assets/dash-logo.svg?react';
import ZohoLogo from '@/assets/zoho-logo.svg?react';
import KeeperLogo from '@/assets/keeper-logo.png';
import {useEffect} from 'react';

export const SelectWallet = () => {
  const {control, reset} = useFormContext<WalletProviderFormValues>();
  const methods = useStepper();

  const sourceProviders: WalletProviderProps[] = [
    {
      type: WalletProviders.ONE_PASSWORD,
      isDisabled: false,
      walletName: '1Password',
      walletDetails: 'Secure password manager',
      imgURI: <PasswordLogo className="w-12 h-12" />
    },
    {
      type: WalletProviders.BIT_WARDEN,
      isDisabled: false,
      walletName: 'Bitwarden',
      walletDetails: 'Open-source password manager',
      imgURI: <BitWardenLogo className="w-12 h-9" />
    },
    {
      type: WalletProviders.DROPBOX,
      isDisabled: false,
      walletName: 'Dropbox',
      walletDetails: 'File hosting service',
      imgURI: <DropboxLogo className="w-12 h-10" />
    },
    {
      type: WalletProviders.PROTON_PASSWORD,
      isDisabled: false,
      walletName: 'ProtonPass',
      walletDetails: 'Secure password manager',
      imgURI: <ProtonLogo className="w-12 h-10" />
    },
    {
      type: WalletProviders.DASH_LANE,
      isDisabled: false,
      walletName: 'Dashlane',
      walletDetails: 'Password manager',
      imgURI: <DashLogo className="w-12 h-9 fill-white" />
    },
    {
      type: WalletProviders.ZOHO,
      isDisabled: false,
      walletName: 'Zoho Vault',
      walletDetails: 'Password manager',
      imgURI: <ZohoLogo className="w-20 h-10" />
    },
    {
      type: WalletProviders.KEEPER,
      isDisabled: false,
      walletName: 'Keeper',
      walletDetails: 'Password manager',
      imgURI: <img src={KeeperLogo} className="w-12 h-10" />
    }
  ];

  const provider = methods.getMetadata('connectWallet')?.provider as WalletProviders;

  useEffect(() => {
    if (provider) {
      reset({
        provider: provider
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [provider]);

  return (
    <Card className="text-start" variant="secondary">
      <CardHeader className="p-4">
        <CardTitle>{methods.get('connectWallet').title}</CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <FormField
          control={control}
          name="provider"
          render={({field}) => (
            <FormItem>
              <div className="card-group">
                {sourceProviders.map((provider, index) => (
                  <WalletProvider key={index} {...provider} isSelected={field.value === provider.type} onSelect={field.onChange} />
                ))}
                <div className="card-flex-group min-w-[300px] hidden lg:block" />
                <div className="card-flex-group min-w-[300px] hidden lg:block" />
              </div>
              <FormDescription>Select your wallet provider.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      </CardContent>
    </Card>
  );
};
