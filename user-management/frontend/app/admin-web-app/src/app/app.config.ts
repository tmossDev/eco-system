import { ApplicationConfig, inject, provideAppInitializer, provideBrowserGlobalErrorListeners } from '@angular/core';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideRouter } from '@angular/router';
import {
  RuntimeConfigService,
  apiBaseUrlInterceptor,
  authInterceptor,
  errorInterceptor,
} from '@eco/auth-features';
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
