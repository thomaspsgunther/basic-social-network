import { User } from '@/src/features/shared/data/models/User';

export interface Comment {
  id: string;
  user: User;
  message: string;
  likeCount?: number;
  createdAt: Date;
}
