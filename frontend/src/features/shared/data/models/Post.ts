import { User } from './User';

export interface Post {
  id: string;
  user?: User;
  image?: string;
  description?: string;
  likeCount?: number;
  commentCount?: number;
  createdAt?: Date;
}
