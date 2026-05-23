import {
  HttpErrorResponse,
  HttpInterceptorFn,
  HttpResponse,
} from '@angular/common/http';
import { inject } from '@angular/core';
import { of, throwError } from 'rxjs';
import { delay } from 'rxjs/operators';

import { RuntimeConfigService } from '../../config/runtime-config.service';
import {
  MOCK_LOGIN_RESPONSE,
  isMockLoginValid,
} from '../../services/auth/auth.mocks';
import { LoginRequest } from '../../services/auth/auth.models';
import { MOCK_DASHBOARD_SUMMARY } from '../../services/dashboard/dashboard.mocks';
import { MOCK_SETTINGS } from '../../services/settings/settings.mocks';
import { ApplicationSettings } from '../../services/settings/settings.models';
import {
  createMockProduct,
  deleteMockProduct,
  getMockProductById,
  getMockProducts,
  updateMockProduct,
  uploadMockProductPhoto,
} from '../../services/product/product.mocks';
import {
  CreateProductRequest,
  UpdateProductRequest,
} from '../../services/product/product.model';

const MOCK_API_DELAY_MS = 350;

export const mockApiInterceptor: HttpInterceptorFn = (request, next) => {
  const runtimeConfig = inject(RuntimeConfigService);
  const url = normaliseApiUrl(request.url);
  const method = request.method.toUpperCase();

  if (!runtimeConfig.mockApiEnabled || !url.startsWith('/api')) {
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

  if (method === 'GET' && url === '/api/products') {
    return mockResponse(getMockProducts());
  }

  if (method === 'POST' && url === '/api/products') {
    return mockResponse(
      createMockProduct(request.body as CreateProductRequest),
      201,
    );
  }

  const productId = getProductIdFromUrl(url);

  if (productId && method === 'POST' && url.endsWith('/photos')) {
    const file = (request.body as FormData).get('file') as File | null;

    return mockResponse(uploadMockProductPhoto(productId, file?.name ?? ''), 201);
  }

  if (productId && method === 'GET') {
    const product = getMockProductById(productId);

    if (!product) {
      return mockError(404, `Product with id "${productId}" was not found`);
    }

    return mockResponse(product);
  }

  if (productId && method === 'PUT') {
    return mockResponse(
      updateMockProduct(productId, request.body as UpdateProductRequest),
    );
  }

  if (productId && method === 'DELETE') {
    deleteMockProduct(productId);

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

function getProductIdFromUrl(url: string): string | null {
  const match = url.match(/^\/api\/products\/([^/]+)(?:\/photos)?$/);

  return match?.[1] ?? null;
}
