import { CreateUserRequest, UpdateUserRequest, UserDetails, UserSummary } from './user.model';

export const MOCK_USERS: UserSummary[] = [
  {
    id: '1',
    name: 'Alex Morgan',
    email: 'alex.morgan@example.com',
    role: 'Admin',
    status: 'Active',
  },
  {
    id: '2',
    name: 'Priya Shah',
    email: 'priya.shah@example.com',
    role: 'Manager',
    status: 'Active',
  },
  {
    id: '3',
    name: 'Jordan Lee',
    email: 'jordan.lee@example.com',
    role: 'User',
    status: 'Pending',
  },
  {
    id: '4',
    name: 'Sam Taylor',
    email: 'sam.taylor@example.com',
    role: 'User',
    status: 'Suspended',
  },
];

let mockUsers = [...MOCK_USERS];

export function getMockUsers(): UserSummary[] {
  return [...mockUsers];
}

export function getMockUserById(id: string): UserDetails | undefined {
  return mockUsers.find((user) => user.id === id);
}

export function createMockUser(request: CreateUserRequest): UserDetails {
  const user: UserDetails = {
    id: createMockUserId(),
    ...request,
  };

  mockUsers = [user, ...mockUsers];

  return user;
}

export function updateMockUser(id: string, request: UpdateUserRequest): UserDetails {
  const updatedUser: UserDetails = {
    id,
    ...request,
  };

  mockUsers = mockUsers.map((user) => (user.id === id ? updatedUser : user));

  return updatedUser;
}

export function deleteMockUser(id: string): void {
  mockUsers = mockUsers.filter((user) => user.id !== id);
}

export function createMockUserId(): string {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    return crypto.randomUUID();
  }

  return String(Date.now());
}
