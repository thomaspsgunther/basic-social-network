import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';

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
  userError: string | null;
  followError: string | null;
}

const initialState: UsersState = {
  users: [],
  searchUsers: [],
  followers: [],
  loadingCurrentUser: false,
  loadingUser: false,
  loadingFollow: false,
  follows: null,
  userError: null,
  followError: null,
};

export const getUsersById = createAsyncThunk(
  'users/getUsersById',
  async (idList: string) => {
    const users = await userUsecase.getUsersById(idList);
    return users;
  },
);

export const getUsersBySearch = createAsyncThunk(
  'users/getUsersBySearch',
  async (searchTerm: string) => {
    const users = await userUsecase.getUsersBySearch(searchTerm);
    return users;
  },
);

export const updateUser = createAsyncThunk(
  'users/updateUser',
  async (user: User) => {
    const didUpdate = await userUsecase.updateUser(user);
    if (didUpdate) {
      return didUpdate;
    } else {
      throw new Error('failed to update user');
    }
  },
);

export const deleteUser = createAsyncThunk(
  'users/deleteUser',
  async (id: string) => {
    const didDelete = await userUsecase.deleteUser(id);
    if (didDelete) {
      return didDelete;
    } else {
      throw new Error('failed to delete user');
    }
  },
);

export const followUser = createAsyncThunk(
  'users/followUser',
  async (args: [string, string]) => {
    const [followerId, followedId] = args;
    const didFollow = await userUsecase.followUser(followerId, followedId);
    if (didFollow) {
      return didFollow;
    } else {
      throw new Error('failed to follow user');
    }
  },
);

export const unfollowUser = createAsyncThunk(
  'users/unfollowUser',
  async (args: [string, string]) => {
    const [followerId, followedId] = args;
    const didUnfollow = await userUsecase.unfollowUser(followerId, followedId);
    if (didUnfollow) {
      return didUnfollow;
    } else {
      throw new Error('failed to unfollow user');
    }
  },
);

export const userFollowsUser = createAsyncThunk(
  'users/userFollowsUser',
  async (args: [string, string]) => {
    const [followerId, followedId] = args;
    const follows = await userUsecase.userFollowsUser(followerId, followedId);
    return follows;
  },
);

export const getUserFollowers = createAsyncThunk(
  'users/getUserFollowers',
  async (id: string) => {
    const users = await userUsecase.getUserFollowers(id);
    return users;
  },
);

export const getUserFollowed = createAsyncThunk(
  'users/getUserFollowed',
  async (id: string) => {
    const users = await userUsecase.getUserFollowed(id);
    return users;
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
        state.userError = null;
      })
      .addCase(getUsersById.fulfilled, (state, action) => {
        state.loadingUser = false;
        state.users = action.payload;
        state.userError = null;
      })
      .addCase(getUsersById.rejected, (state, action) => {
        state.loadingUser = false;
        state.userError = action.error.message || null;
      })
      .addCase(getUsersBySearch.pending, (state) => {
        state.loadingUser = true;
        state.userError = null;
      })
      .addCase(getUsersBySearch.fulfilled, (state, action) => {
        state.loadingUser = false;
        state.searchUsers = action.payload;
        state.userError = null;
      })
      .addCase(getUsersBySearch.rejected, (state, action) => {
        state.loadingUser = false;
        state.userError = action.error.message || null;
      })
      .addCase(updateUser.pending, (state) => {
        state.loadingCurrentUser = true;
        state.userError = null;
      })
      .addCase(updateUser.fulfilled, (state) => {
        state.loadingCurrentUser = false;
        state.userError = null;
      })
      .addCase(updateUser.rejected, (state, action) => {
        state.loadingCurrentUser = false;
        state.userError = action.error.message || null;
      })
      .addCase(deleteUser.pending, (state) => {
        state.loadingCurrentUser = true;
        state.userError = null;
      })
      .addCase(deleteUser.fulfilled, (state) => {
        state.loadingCurrentUser = false;
        state.userError = null;
      })
      .addCase(deleteUser.rejected, (state, action) => {
        state.loadingCurrentUser = false;
        state.userError = action.error.message || null;
      })
      .addCase(followUser.pending, (state) => {
        state.loadingFollow = true;
        state.followError = null;
      })
      .addCase(followUser.fulfilled, (state) => {
        state.loadingFollow = false;
        state.followError = null;
      })
      .addCase(followUser.rejected, (state, action) => {
        state.loadingFollow = false;
        state.followError = action.error.message || null;
      })
      .addCase(unfollowUser.pending, (state) => {
        state.loadingFollow = true;
        state.followError = null;
      })
      .addCase(unfollowUser.fulfilled, (state) => {
        state.loadingFollow = false;
        state.followError = null;
      })
      .addCase(unfollowUser.rejected, (state, action) => {
        state.loadingFollow = false;
        state.followError = action.error.message || null;
      })
      .addCase(userFollowsUser.pending, (state) => {
        state.loadingFollow = true;
        state.followError = null;
      })
      .addCase(userFollowsUser.fulfilled, (state, action) => {
        state.loadingFollow = false;
        state.follows = action.payload;
        state.followError = null;
      })
      .addCase(userFollowsUser.rejected, (state, action) => {
        state.loadingFollow = false;
        state.follows = null;
        state.followError = action.error.message || null;
      })
      .addCase(getUserFollowers.pending, (state) => {
        state.loadingFollow = true;
        state.followError = null;
      })
      .addCase(getUserFollowers.fulfilled, (state, action) => {
        state.loadingFollow = false;
        state.followers = action.payload;
        state.followError = null;
      })
      .addCase(getUserFollowers.rejected, (state, action) => {
        state.loadingFollow = false;
        state.followError = action.error.message || null;
      })
      .addCase(getUserFollowed.pending, (state) => {
        state.loadingFollow = true;
        state.followError = null;
      })
      .addCase(getUserFollowed.fulfilled, (state, action) => {
        state.loadingFollow = false;
        state.followers = action.payload;
        state.followError = null;
      })
      .addCase(getUserFollowed.rejected, (state, action) => {
        state.loadingFollow = false;
        state.followError = action.error.message || null;
      });
  },
});

export const usersReducer = usersSlice.reducer;
