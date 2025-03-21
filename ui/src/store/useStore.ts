import {WalletProviders} from '@/types/wallet-providers';
import {create, StateCreator} from 'zustand';
import {createJSONStorage, persist, PersistOptions} from 'zustand/middleware';

type Store = {
  walletProvider?: WalletProviders;
  setWalletProvider: (provider: WalletProviders) => void;
  nodeUrl?: string;
  setNodeUrl: (url: string) => void;
  cleanStore: () => void;
};

type PersistStore = (config: StateCreator<Store>, options: PersistOptions<Store>) => StateCreator<Store>;

export const useStore = create<Store>(
  (persist as PersistStore)(
    (set): Store => ({
      walletProvider: undefined,
      setWalletProvider: (provider: WalletProviders) => set(() => ({walletProvider: provider})),
      nodeUrl: undefined,
      setNodeUrl: (url: string) => set(() => ({nodeUrl: url})),
      cleanStore: () => set(() => ({walletProvider: undefined, nodeUrl: undefined}))
    }),
    {
      name: 'pyramd-ui-store',
      storage: createJSONStorage(() => localStorage)
    }
  )
);
