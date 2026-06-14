import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';

import { AuthService } from '../services/auth/auth.service';

export const errorInterceptor: HttpInterceptorFn = (request, next) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  return next(request).pipe(
    catchError((error: HttpErrorResponse) => {
      const isLoginRequest = request.url.includes('/api/auth/login');

      if (error.status === 401 && !isLoginRequest) {
        authService.logout();
        void router.navigate(['/auth/login']);
      }

      if (error.status === 403) {
        void router.navigate(['/']);
      }

      return throwError(() => error);
    }),
  );
};
