export default interface User {
  id: string;
  username: string;
  password?: string;
  email?: string;
  fullName?: string;
  description?: string;
  avatar?: string;
  followerCount?: string;
}
