import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {useStepper} from '../stepper';
import {Alert, AlertDescription, AlertTitle} from '@/components/ui/alert';

export const VerifyConnection = () => {
  const methods = useStepper();

  return (
    <Card className="text-start" variant="secondary">
      <CardHeader className="p-4">
        <CardTitle>{methods.get('verifyTheConnection').title}</CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <Alert variant="default">
          <AlertTitle>Verify the connection to the identity network.</AlertTitle>
          <AlertDescription>You can test the connection to the identity network by pressing the button below.</AlertDescription>
        </Alert>
      </CardContent>
    </Card>
  );
};
