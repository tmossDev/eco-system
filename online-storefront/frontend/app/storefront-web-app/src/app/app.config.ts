import {
  ApplicationConfig,
  inject,
  provideAppInitializer,
  provideBrowserGlobalErrorListeners,
} from '@angular/core';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { RuntimeConfigService } from './core/config/runtime-config.service';
import { apiBaseUrlInterceptor } from './core/interceptors/api-base-url.interceptor';
import { authInterceptor } from './core/interceptors/auth.interceptor';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideAppInitializer(() => inject(RuntimeConfigService).load()),
    provideRouter(routes),
    provideHttpClient(
      withInterceptors([authInterceptor, apiBaseUrlInterceptor]),
    ),
  ],
};
