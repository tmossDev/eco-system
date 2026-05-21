import { AuthUser, LoginResponse } from './auth.models';

export const MOCK_AUTH_USER: AuthUser = {
  id: '1',
  name: 'Product Admin',
  email: 'admin@example.com',
  role: 'Admin',
};

export const MOCK_ACCESS_TOKEN = 'mock-admin-access-token';

export const MOCK_LOGIN_RESPONSE: LoginResponse = {
  accessToken: MOCK_ACCESS_TOKEN,
  user: MOCK_AUTH_USER,
};

export const MOCK_AUTH_CREDENTIALS = {
  email: 'admin@example.com',
  password: 'password',
};

export function isMockLoginValid(email: string, password: string): boolean {
  return (
    email.trim().toLowerCase() === MOCK_AUTH_CREDENTIALS.email &&
    password === MOCK_AUTH_CREDENTIALS.password
  );
}
