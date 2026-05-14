export type UserStatus = 'Active' | 'Pending' | 'Suspended';

export interface UserSummary {
  id: string;
  name: string;
  email: string;
  role: string;
  status: UserStatus;
}

export interface UserDetails extends UserSummary {}

export interface CreateUserRequest {
  name: string;
  email: string;
  role: string;
  status: UserStatus;
}

export interface UpdateUserRequest {
  name: string;
  email: string;
  role: string;
  status: UserStatus;
}
