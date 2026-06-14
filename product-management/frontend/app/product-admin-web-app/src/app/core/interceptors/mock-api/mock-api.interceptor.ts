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
import { MOCK_DASHBOARD_SUMMARY } from '../../services/dashboard/dashboard.mocks';
import { MOCK_SETTINGS } from '../../services/settings/settings.mocks';
import { ApplicationSettings } from '../../services/settings/settings.models';
import {
  createMockDiscount,
  createMockProduct,
  deleteMockDiscount,
  deleteMockProduct,
  getMockDiscountById,
  getMockDiscounts,
  getMockPromotionSettings,
  getMockProductById,
  getMockProducts,
  updateMockDiscount,
  updateMockPromotionSettings,
  updateMockProduct,
  uploadMockProductPhoto,
} from '../../services/product/product.mocks';
import {
  CreateDiscountRequest,
  CreateProductRequest,
  UpdateDiscountRequest,
  UpdatePromotionSettingsRequest,
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

  if (method === 'GET' && url === '/api/products') {
    return mockResponse(getMockProducts());
  }

  if (method === 'POST' && url === '/api/products') {
    return mockResponse(
      createMockProduct(request.body as CreateProductRequest),
      201,
    );
  }

  if (method === 'GET' && url === '/api/discounts') {
    return mockResponse(getMockDiscounts());
  }

  if (method === 'POST' && url === '/api/discounts') {
    return mockResponse(
      createMockDiscount(request.body as CreateDiscountRequest),
      201,
    );
  }

  if (method === 'GET' && url === '/api/promotions/settings') {
    return mockResponse(getMockPromotionSettings());
  }

  if (method === 'PUT' && url === '/api/promotions/settings') {
    const body = request.body as UpdatePromotionSettingsRequest;

    return mockResponse(updateMockPromotionSettings(body.promotions_enabled));
  }

  const discountId = getDiscountIdFromUrl(url);

  if (discountId && method === 'GET') {
    const discount = getMockDiscountById(discountId);

    if (!discount) {
      return mockError(404, `Discount with id "${discountId}" was not found`);
    }

    return mockResponse(discount);
  }

  if (discountId && method === 'PUT') {
    return mockResponse(
      updateMockDiscount(discountId, request.body as UpdateDiscountRequest),
    );
  }

  if (discountId && method === 'DELETE') {
    deleteMockDiscount(discountId);

    return mockResponse(undefined);
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

function getDiscountIdFromUrl(url: string): string | null {
  const match = url.match(/^\/api\/discounts\/([^/]+)$/);

  return match?.[1] ?? null;
}
