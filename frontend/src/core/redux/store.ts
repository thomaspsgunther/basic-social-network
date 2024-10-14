import { configureStore } from '@reduxjs/toolkit';

import { usersReducer } from '@/src/features/users/presentation/redux/usersSlice';

import { errorsReducer } from '../errors/errorsSlice';

export const store = configureStore({
  reducer: {
    errors: errorsReducer,
    users: usersReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
