import { create } from "zustand";

type Data = {
    modalStack: string[];
    modalName: string | null;
    OpenModal: (name: string) => void;
    CloseModal: () => void;
    reset: () => void;
};

export const useModal = create<Data>((set) => ({
    modalStack: [],
    modalName: null,

    OpenModal: (name: string) =>
        set((state) => {
            const modalStack = [...state.modalStack, name];
            return {
                modalStack,
                modalName: modalStack.at(-1) ?? null,
            };
        }),

    CloseModal: () =>
        set((state) => {
            const newStack = state.modalStack.slice(0, -1);

            return {
                modalStack: newStack,
                modalName: newStack.length > 0 ? newStack[newStack.length - 1] : null,
            };
        }),

    reset: () => set({ modalStack: [], modalName: null }),
}));