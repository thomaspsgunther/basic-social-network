import { configureStore } from '@reduxjs/toolkit';

import { usersReducer } from '@/src/features/users/presentation/redux/usersSlice';

export const store = configureStore({
  reducer: {
    users: usersReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
