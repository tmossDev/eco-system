import { ApplicationConfig, inject, provideAppInitializer, provideBrowserGlobalErrorListeners } from '@angular/core';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideRouter } from '@angular/router';

import { RuntimeConfigService } from './core/config/runtime-config.service';
import { apiBaseUrlInterceptor } from './core/interceptors/api-base-url/api-base-url.interceptor';
import { authInterceptor } from './core/interceptors/auth/auth.interceptor';
import { errorInterceptor } from './core/interceptors/error/error.interceptor';
import { mockApiInterceptor } from './core/interceptors/mock-api/mock-api.interceptor';
import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideAppInitializer(() => inject(RuntimeConfigService).load()),
    provideRouter(routes),
    provideHttpClient(
      withInterceptors([
        mockApiInterceptor,
        authInterceptor,
        apiBaseUrlInterceptor,
        errorInterceptor,
      ]),
    ),
  ],
};
