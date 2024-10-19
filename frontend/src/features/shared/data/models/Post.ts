export interface Post {
    id: string;
	user: string;
	image: string;
	description: string;
	likeCount: number;
	commentCount: number;
	createdAt: Date;
}  