import {z} from 'zod';
import {WalletProviders} from '@/types/wallet-providers';

export const WalletProviderSchema = z.object({
  provider: z.nativeEnum(WalletProviders)
});
export type WalletProviderFormValues = z.infer<typeof WalletProviderSchema>;

export const GenerateStoreSchema = z.object({
  storeKeys: z.boolean().optional()
});
export type GenerateStoreFormValues = z.infer<typeof GenerateStoreSchema>;
