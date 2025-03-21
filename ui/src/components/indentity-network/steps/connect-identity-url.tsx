import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {useStepper} from '../stepper';
import {useFormContext} from 'react-hook-form';
import {FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage} from '@/components/ui/form';
import {useEffect} from 'react';
import {ConnectIdentityNetworFormValues} from '@/schemas/identity-network-schema';
import {Input} from '@/components/ui/input';

export const ConnectIdentiyURL = () => {
  const {control, reset} = useFormContext<ConnectIdentityNetworFormValues>();
  const methods = useStepper();

  const nodeUrl = methods.getMetadata('verifyTheConnection')?.nodeUrl as string;

  useEffect(() => {
    if (nodeUrl) {
      reset({
        nodeUrl: nodeUrl
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [nodeUrl]);

  return (
    <Card className="text-start" variant="secondary">
      <CardHeader className="p-4">
        <CardTitle>{methods.get('connectIdentityNetwork').title}</CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <FormField
          control={control}
          name="nodeUrl"
          render={({field}) => (
            <FormItem>
              <FormLabel>Identity Node Address</FormLabel>
              <FormControl>
                <Input placeholder="Node URL..." {...field} />
              </FormControl>
              <FormDescription>Enter the address of the identity node you want to connect to.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      </CardContent>
    </Card>
  );
};
