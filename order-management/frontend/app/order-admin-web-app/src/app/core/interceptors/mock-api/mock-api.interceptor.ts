import {
  HttpErrorResponse,
  HttpInterceptorFn,
  HttpResponse,
} from '@angular/common/http';
import { inject } from '@angular/core';
import { of, throwError } from 'rxjs';
import { delay } from 'rxjs/operators';
import {
  RuntimeConfigService,
  mockAuthApiResponse,
} from '@eco/auth-features';

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
