import { HttpErrorResponse, HttpInterceptorFn, HttpResponse } from '@angular/common/http';
import { of, throwError } from 'rxjs';
import { delay } from 'rxjs/operators';

import {
  MOCK_LOGIN_RESPONSE,
  isMockLoginValid,
} from '../../services/auth/auth.mocks';
import { LoginRequest } from '../../services/auth/auth.models';
import { MOCK_DASHBOARD_SUMMARY } from '../../services/dashboard/dashboard.mocks';
import { MOCK_SETTINGS } from '../../services/settings/settings.mocks';
import { ApplicationSettings } from '../../services/settings/settings.models';
import {
  createMockUser,
  deleteMockUser,
  getMockUserById,
  getMockUsers,
  updateMockUser,
} from '../../services/user/user.mocks';
import {
  CreateUserRequest,
  UpdateUserRequest,
} from '../../services/user/user.model';

const MOCK_API_DELAY_MS = 350;
const ENABLE_MOCK_API = true;

export const mockApiInterceptor: HttpInterceptorFn = (request, next) => {
  const url = normaliseApiUrl(request.url);
  const method = request.method.toUpperCase();

  if (!ENABLE_MOCK_API || !url.startsWith('/api')) {
    return next(request);
  }

  if (method === 'POST' && url === '/api/auth/login') {
    const body = request.body as LoginRequest;

    if (!isMockLoginValid(body.email, body.password)) {
      return mockError(401, 'Invalid email or password');
    }

    return mockResponse(MOCK_LOGIN_RESPONSE);
  }

  if (method === 'POST' && url === '/api/auth/forgot-password') {
    return mockResponse(undefined);
  }

  if (method === 'GET' && url === '/api/dashboard/summary') {
    return mockResponse(MOCK_DASHBOARD_SUMMARY);
  }

  if (method === 'GET' && url === '/api/settings') {
    return mockResponse(MOCK_SETTINGS);
  }

  if (method === 'PUT' && url === '/api/settings') {
    return mockResponse(request.body as ApplicationSettings);
  }

  if (method === 'GET' && url === '/api/users') {
    return mockResponse(getMockUsers());
  }

  if (method === 'POST' && url === '/api/users') {
    return mockResponse(createMockUser(request.body as CreateUserRequest), 201);
  }

  const userId = getUserIdFromUrl(url);

  if (userId && method === 'GET') {
    const user = getMockUserById(userId);

    if (!user) {
      return mockError(404, `User with id "${userId}" was not found`);
    }

    return mockResponse(user);
  }

  if (userId && method === 'PUT') {
    return mockResponse(updateMockUser(userId, request.body as UpdateUserRequest));
  }

  if (userId && method === 'DELETE') {
    deleteMockUser(userId);

    return mockResponse(undefined);
  }

  return next(request);
};

function mockResponse<T>(body: T, status = 200) {
  return of(
    new HttpResponse<T>({
      status,
      body,
    }),
  ).pipe(delay(MOCK_API_DELAY_MS));
}

function mockError(status: number, message: string) {
  return throwError(
    () =>
      new HttpErrorResponse({
        status,
        error: {
          message,
        },
      }),
  ).pipe(delay(MOCK_API_DELAY_MS));
}

function normaliseApiUrl(url: string): string {
  try {
    return new URL(url, window.location.origin).pathname;
  } catch {
    return url.split('?')[0];
  }
}

function getUserIdFromUrl(url: string): string | null {
  const match = url.match(/^\/api\/users\/([^/]+)$/);

  return match?.[1] ?? null;
}
