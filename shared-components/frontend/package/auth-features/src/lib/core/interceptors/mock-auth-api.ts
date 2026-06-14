import {
  HttpErrorResponse,
  HttpEvent,
  HttpRequest,
  HttpResponse,
} from '@angular/common/http';
import { Observable, of, throwError } from 'rxjs';
import { delay } from 'rxjs/operators';

import {
  MOCK_LOGIN_RESPONSE,
  isMockLoginValid,
} from '../services/auth/auth.mocks';
import { LoginRequest } from '../services/auth/auth.models';

const MOCK_AUTH_API_DELAY_MS = 350;

export function mockAuthApiResponse(
  request: HttpRequest<unknown>,
  url: string,
): Observable<HttpEvent<unknown>> | null {
  const method = request.method.toUpperCase();

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

  return null;
}

function mockResponse<T>(body: T, status = 200) {
  return of(
    new HttpResponse<T>({
      status,
      body,
    }),
  ).pipe(delay(MOCK_AUTH_API_DELAY_MS));
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
  ).pipe(delay(MOCK_AUTH_API_DELAY_MS));
}
