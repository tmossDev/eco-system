import { HttpErrorResponse, HttpInterceptorFn, HttpResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { of, throwError } from 'rxjs';
import { delay } from 'rxjs/operators';
import {
  RuntimeConfigService,
  mockAuthApiResponse,
} from '@eco/auth-features';
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

export const mockApiInterceptor: HttpInterceptorFn = (request, next) => {
  const runtimeConfig = inject(RuntimeConfigService);
  const url = normaliseApiUrl(request.url);
  const method = request.method.toUpperCase();

  if (!runtimeConfig.mockApiEnabled || !url.startsWith('/api')) {
    return next(request);
  }

  const authMockResponse = mockAuthApiResponse(request, url);
  if (authMockResponse) {
    return authMockResponse;
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
