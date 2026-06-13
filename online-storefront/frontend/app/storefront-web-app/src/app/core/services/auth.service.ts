import { HttpClient } from '@angular/common/http';
import { Injectable, computed, inject, signal } from '@angular/core';
import { Observable, finalize, tap } from 'rxjs';

import {
  AuthUser,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
} from './auth.models';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly tokenKey = 'storefront_access_token';
  private readonly userKey = 'storefront_user';
  private readonly handoffTokenParam = 'eco_access_token';
  private readonly handoffUserParam = 'eco_user';

  private readonly tokenSignal = signal<string | null>(
    localStorage.getItem(this.tokenKey),
  );
  private readonly userSignal = signal<AuthUser | null>(this.storedUser());

  readonly accessToken = computed(() => this.tokenSignal());
  readonly currentUser = computed(() => this.userSignal());
  readonly isAuthenticated = computed(() => Boolean(this.tokenSignal()));

  constructor() {
    this.consumeSessionHandoff();
  }

  login(request: LoginRequest): Observable<LoginResponse> {
    return this.http
      .post<LoginResponse>('/api/auth/login', request)
      .pipe(tap((response) => this.storeSession(response)));
  }

  register(request: RegisterRequest): Observable<LoginResponse> {
    return this.http
      .post<LoginResponse>('/api/register', request)
      .pipe(tap((response) => this.storeSession(response)));
  }

  logout(): Observable<unknown> {
    return this.http
      .post('/api/auth/logout', {})
      .pipe(finalize(() => this.clearSession()));
  }

  private storeSession(response: LoginResponse): void {
    localStorage.setItem(this.tokenKey, response.accessToken);
    localStorage.setItem(this.userKey, JSON.stringify(response.user));
    this.tokenSignal.set(response.accessToken);
    this.userSignal.set(response.user);
  }

  private consumeSessionHandoff(): void {
    const params = new URLSearchParams(window.location.hash.replace(/^#/, ''));
    const accessToken = params.get(this.handoffTokenParam);
    const userValue = params.get(this.handoffUserParam);

    if (!accessToken || !userValue) {
      return;
    }

    try {
      const user = JSON.parse(userValue) as AuthUser;
      this.storeSession({ accessToken, user });
      params.delete(this.handoffTokenParam);
      params.delete(this.handoffUserParam);
      const nextHash = params.toString();
      const nextUrl = `${window.location.pathname}${window.location.search}${nextHash ? `#${nextHash}` : ''}`;
      window.history.replaceState(null, '', nextUrl);
    } catch {
      return;
    }
  }

  private clearSession(): void {
    localStorage.removeItem(this.tokenKey);
    localStorage.removeItem(this.userKey);
    this.tokenSignal.set(null);
    this.userSignal.set(null);
  }

  private storedUser(): AuthUser | null {
    const value = localStorage.getItem(this.userKey);
    if (!value) {
      return null;
    }

    try {
      return JSON.parse(value) as AuthUser;
    } catch {
      return null;
    }
  }
}
