import { Comment } from '../../data/models/Comment';
import { ICommentRepository } from '../repositories/CommentRepository';

interface ICommentUsecase {
  createComment(comment: Omit<Comment, 'id'>): Promise<Comment>;
  getCommentsFromPost(postId: string): Promise<Comment[]>;
  updateComment(comment: Comment): Promise<boolean>;
  deleteComment(id: string): Promise<boolean>;
}

export class CommentUsecaseImpl implements ICommentUsecase {
  private repository: ICommentRepository;

  constructor(repository: ICommentRepository) {
    this.repository = repository;
  }

  async createComment(comment: Omit<Comment, 'id'>): Promise<Comment> {
    if (!comment.message || !comment.user || !comment.postId) {
      throw new Error('comment message, user and postId are required');
    }

    const createdComment = await this.repository.createComment(comment);

    return createdComment;
  }

  async getCommentsFromPost(postId: string): Promise<Comment[]> {
    const comments: Comment[] =
      await this.repository.getCommentsFromPost(postId);

    return comments;
  }

  async updateComment(comment: Comment): Promise<boolean> {
    if (!comment.message || !comment.user || !comment.postId) {
      throw new Error('comment message, user and postId are required');
    }

    await this.repository.updateComment(comment);

    return true;
  }

  async deleteComment(id: string): Promise<boolean> {
    await this.repository.deleteComment(id);

    return true;
  }
}
