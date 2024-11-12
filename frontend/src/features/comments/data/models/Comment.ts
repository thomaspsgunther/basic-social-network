import { User } from '@/src/features/shared/data/models/User';

export interface Comment {
  id: string;
  postId?: string;
  user?: User;
  message?: string;
  createdAt?: Date;
}
