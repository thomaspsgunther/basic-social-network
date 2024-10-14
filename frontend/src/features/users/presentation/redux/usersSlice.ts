import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';

import { setError } from '@/src/core/errors/errorsSlice';
import { User } from '@/src/features/shared/data/models/User';

import { UserRepositoryImpl } from '../../data/repositories/UserRepositoryImpl';
import { UserUsecaseImpl } from '../../domain/usecases/UserUsecase';

const userRepository = new UserRepositoryImpl();
const userUsecase = new UserUsecaseImpl(userRepository);

interface UsersState {
  users: User[];
  searchUsers: User[];
  followers: User[];
  loadingCurrentUser: boolean;
  loadingUser: boolean;
  loadingFollow: boolean;
  follows: boolean | null;
}

const initialState: UsersState = {
  users: [],
  searchUsers: [],
  followers: [],
  loadingCurrentUser: false,
  loadingUser: false,
  loadingFollow: false,
  follows: null,
};

export const getUsersById = createAsyncThunk(
  'users/getUsersById',
  async (idList: string, { dispatch }) => {
    try {
      const users = await userUsecase.getUsersById(idList);

      return users;
    } catch (error) {
      if (error instanceof Error) {
        dispatch(setError(error.message));
      } else {
        dispatch(setError('unknown error'));
      }
      throw error;
    }
  },
);

export const getUsersBySearch = createAsyncThunk(
  'users/getUsersBySearch',
  async (searchTerm: string, { dispatch }) => {
    try {
      const users = await userUsecase.getUsersBySearch(searchTerm);

      return users;
    } catch (error) {
      if (error instanceof Error) {
        dispatch(setError(error.message));
      } else {
        dispatch(setError('unknown error'));
      }
      throw error;
    }
  },
);

export const updateUser = createAsyncThunk(
  'users/updateUser',
  async (user: User, { dispatch }) => {
    try {
      const didUpdate = await userUsecase.updateUser(user);

      return didUpdate;
    } catch (error) {
      if (error instanceof Error) {
        dispatch(setError(error.message));
      } else {
        dispatch(setError('unknown error'));
      }
      throw error;
    }
  },
);

export const deleteUser = createAsyncThunk(
  'users/deleteUser',
  async (id: string, { dispatch }) => {
    try {
      const didDelete = await userUsecase.deleteUser(id);

      return didDelete;
    } catch (error) {
      if (error instanceof Error) {
        dispatch(setError(error.message));
      } else {
        dispatch(setError('unknown error'));
      }
      throw error;
    }
  },
);

export const followUser = createAsyncThunk(
  'users/followUser',
  async (args: [string, string], { dispatch }) => {
    try {
      const [followerId, followedId] = args;
      const didFollow = await userUsecase.followUser(followerId, followedId);

      return didFollow;
    } catch (error) {
      if (error instanceof Error) {
        dispatch(setError(error.message));
      } else {
        dispatch(setError('unknown error'));
      }
      throw error;
    }
  },
);

export const unfollowUser = createAsyncThunk(
  'users/unfollowUser',
  async (args: [string, string], { dispatch }) => {
    try {
      const [followerId, followedId] = args;
      const didUnfollow = await userUsecase.unfollowUser(
        followerId,
        followedId,
      );

      return didUnfollow;
    } catch (error) {
      if (error instanceof Error) {
        dispatch(setError(error.message));
      } else {
        dispatch(setError('unknown error'));
      }
      throw error;
    }
  },
);

export const userFollowsUser = createAsyncThunk(
  'users/userFollowsUser',
  async (args: [string, string], { dispatch }) => {
    try {
      const [followerId, followedId] = args;
      const follows = await userUsecase.userFollowsUser(followerId, followedId);

      return follows;
    } catch (error) {
      if (error instanceof Error) {
        dispatch(setError(error.message));
      } else {
        dispatch(setError('unknown error'));
      }
      throw error;
    }
  },
);

export const getUserFollowers = createAsyncThunk(
  'users/getUserFollowers',
  async (id: string, { dispatch }) => {
    try {
      const users = await userUsecase.getUserFollowers(id);

      return users;
    } catch (error) {
      if (error instanceof Error) {
        dispatch(setError(error.message));
      } else {
        dispatch(setError('unknown error'));
      }
      throw error;
    }
  },
);

export const getUserFollowed = createAsyncThunk(
  'users/getUserFollowed',
  async (id: string, { dispatch }) => {
    try {
      const users = await userUsecase.getUserFollowed(id);

      return users;
    } catch (error) {
      if (error instanceof Error) {
        dispatch(setError(error.message));
      } else {
        dispatch(setError('unknown error'));
      }
      throw error;
    }
  },
);

const usersSlice = createSlice({
  name: 'users',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(getUsersById.pending, (state) => {
        state.loadingUser = true;
      })
      .addCase(getUsersById.fulfilled, (state, action) => {
        state.loadingUser = false;
        state.users = action.payload;
      })
      .addCase(getUsersById.rejected, (state) => {
        state.loadingUser = false;
      })
      .addCase(getUsersBySearch.pending, (state) => {
        state.loadingUser = true;
      })
      .addCase(getUsersBySearch.fulfilled, (state, action) => {
        state.loadingUser = false;
        state.searchUsers = action.payload;
      })
      .addCase(getUsersBySearch.rejected, (state) => {
        state.loadingUser = false;
      })
      .addCase(updateUser.pending, (state) => {
        state.loadingCurrentUser = true;
      })
      .addCase(updateUser.fulfilled, (state) => {
        state.loadingCurrentUser = false;
      })
      .addCase(updateUser.rejected, (state) => {
        state.loadingCurrentUser = false;
      })
      .addCase(deleteUser.pending, (state) => {
        state.loadingCurrentUser = true;
      })
      .addCase(deleteUser.fulfilled, (state) => {
        state.loadingCurrentUser = false;
      })
      .addCase(deleteUser.rejected, (state) => {
        state.loadingCurrentUser = false;
      })
      .addCase(followUser.pending, (state) => {
        state.loadingFollow = true;
      })
      .addCase(followUser.fulfilled, (state) => {
        state.loadingFollow = false;
      })
      .addCase(followUser.rejected, (state) => {
        state.loadingFollow = false;
      })
      .addCase(unfollowUser.pending, (state) => {
        state.loadingFollow = true;
      })
      .addCase(unfollowUser.fulfilled, (state) => {
        state.loadingFollow = false;
      })
      .addCase(unfollowUser.rejected, (state) => {
        state.loadingFollow = false;
      })
      .addCase(userFollowsUser.pending, (state) => {
        state.loadingFollow = true;
      })
      .addCase(userFollowsUser.fulfilled, (state, action) => {
        state.loadingFollow = false;
        state.follows = action.payload;
      })
      .addCase(userFollowsUser.rejected, (state) => {
        state.loadingFollow = false;
        state.follows = null;
      })
      .addCase(getUserFollowers.pending, (state) => {
        state.loadingFollow = true;
      })
      .addCase(getUserFollowers.fulfilled, (state, action) => {
        state.loadingFollow = false;
        state.followers = action.payload;
      })
      .addCase(getUserFollowers.rejected, (state) => {
        state.loadingFollow = false;
      })
      .addCase(getUserFollowed.pending, (state) => {
        state.loadingFollow = true;
      })
      .addCase(getUserFollowed.fulfilled, (state, action) => {
        state.loadingFollow = false;
        state.followers = action.payload;
      })
      .addCase(getUserFollowed.rejected, (state) => {
        state.loadingFollow = false;
      });
  },
});

export const usersReducer = usersSlice.reducer;
