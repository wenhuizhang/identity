import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {useStepper} from '../stepper';
import {useFormContext} from 'react-hook-form';
import {GenerateStoreFormValues} from '@/schemas/wallet-schemas';
import {FormControl, FormDescription, FormField, FormItem, FormLabel} from '@/components/ui/form';
import {Checkbox} from '@/components/ui/checkbox';

export const GenerateStoreKeys = () => {
  const {control} = useFormContext<GenerateStoreFormValues>();
  const methods = useStepper();

  return (
    <Card className="text-start" variant="secondary">
      <CardHeader className="p-4">
        <CardTitle>{methods.get('generateAndStoreKeys').title}</CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <FormField
          control={control}
          name="storeKeys"
          render={({field}) => (
            <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4 shadow">
              <FormControl>
                <Checkbox checked={field.value} onCheckedChange={field.onChange} />
              </FormControl>
              <div className="space-y-1 leading-none">
                <FormLabel>Public and Private keys created and stored</FormLabel>
                <FormDescription>Your public and private keys will be created and stored in your wallet.</FormDescription>
              </div>
            </FormItem>
          )}
        />
      </CardContent>
    </Card>
  );
};
