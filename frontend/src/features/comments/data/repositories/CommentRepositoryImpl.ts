import { ICommentRepository } from '../../domain/repositories/CommentRepository';
import { commentApi } from '../api/commentApi';
import { Comment } from '../models/Comment';

export class CommentRepositoryImpl implements ICommentRepository {
  async createComment(comment: Omit<Comment, 'id'>): Promise<Comment> {
    const response = await commentApi.create(comment);
    const createdComment: Comment = response.data;

    return createdComment;
  }

  async getCommentsFromPost(postId: string): Promise<Comment[]> {
    const response = await commentApi.getFromPost(postId);
    const comments: Comment[] = response.data ? response.data : [];

    return comments;
  }

  async updateComment(comment: Comment): Promise<boolean> {
    await commentApi.update(comment);

    return true;
  }

  async deleteComment(id: string): Promise<boolean> {
    await commentApi.remove(id);

    return true;
  }
}
