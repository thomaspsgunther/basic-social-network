import { Comment } from '../../data/models/Comment';

export interface ICommentRepository {
  createComment(comment: Omit<Comment, 'id'>): Promise<Comment>;
  getCommentsFromPost(postId: string): Promise<Comment[]>;
  updateComment(comment: Comment): Promise<boolean>;
  deleteComment(id: string): Promise<boolean>;
}
