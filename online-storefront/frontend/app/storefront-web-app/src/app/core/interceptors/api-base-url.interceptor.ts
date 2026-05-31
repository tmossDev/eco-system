import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';

import { RuntimeConfigService } from '../config/runtime-config.service';

export const apiBaseUrlInterceptor: HttpInterceptorFn = (request, next) => {
  if (!request.url.startsWith('/api')) {
    return next(request);
  }

  const backendApiUrl = inject(RuntimeConfigService).backendApiUrl;
  return next(
    backendApiUrl
      ? request.clone({ url: `${backendApiUrl}${request.url}` })
      : request,
  );
};
