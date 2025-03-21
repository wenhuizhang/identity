import {z} from 'zod';

export const ConnectIdentityNetworkSchema = z.object({
  nodeUrl: z.string().url()
});
export type ConnectIdentityNetworFormValues = z.infer<typeof ConnectIdentityNetworkSchema>;
