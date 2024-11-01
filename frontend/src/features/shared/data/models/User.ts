export interface User {
  id: string;
  username: string;
  password?: string;
  email?: string;
  fullName?: string;
  description?: string;
  avatar?: string;
  postCount?: string;
  followerCount?: string;
  followedCount?: string;
}
