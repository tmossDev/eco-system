import { Injectable, computed, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';

import {
  AuthUser,
  ForgotPasswordRequest,
  LoginRequest,
  LoginResponse,
} from './auth.models';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly http = inject(HttpClient);

  private readonly accessTokenKey = 'admin_access_token';
  private readonly userKey = 'admin_user';
  private readonly handoffTokenParam = 'eco_access_token';
  private readonly handoffUserParam = 'eco_user';

  private readonly accessTokenSignal = signal<string | null>(
    localStorage.getItem(this.accessTokenKey) ?? sessionStorage.getItem(this.accessTokenKey),
  );

  private readonly userSignal = signal<AuthUser | null>(this.getStoredUser());

  public readonly accessToken = computed(() => this.accessTokenSignal());
  public readonly currentUser = computed(() => this.userSignal());
  public readonly isAuthenticated = computed(() => Boolean(this.accessTokenSignal()));

  public login(request: LoginRequest): Observable<LoginResponse> {
    return this.http.post<LoginResponse>('/api/auth/login', request).pipe(
      tap((response) => {
        this.storeSession(response, request.rememberMe);
      }),
    );
  }

  public forgotPassword(request: ForgotPasswordRequest): Observable<void> {
    return this.http.post<void>('/api/auth/forgot-password', request);
  }

  public logout(): void {
    localStorage.removeItem(this.accessTokenKey);
    localStorage.removeItem(this.userKey);
    sessionStorage.removeItem(this.accessTokenKey);
    sessionStorage.removeItem(this.userKey);

    this.accessTokenSignal.set(null);
    this.userSignal.set(null);
  }

  public appendSessionHandoff(url: URL): void {
    const accessToken = this.accessTokenSignal();
    const user = this.userSignal();

    if (!accessToken || !user) {
      return;
    }

    const params = new URLSearchParams(url.hash.replace(/^#/, ''));
    params.set(this.handoffTokenParam, accessToken);
    params.set(this.handoffUserParam, JSON.stringify(user));
    url.hash = params.toString();
  }

  private storeSession(response: LoginResponse, rememberMe: boolean): void {
    const storage = rememberMe ? localStorage : sessionStorage;

    localStorage.removeItem(this.accessTokenKey);
    localStorage.removeItem(this.userKey);
    sessionStorage.removeItem(this.accessTokenKey);
    sessionStorage.removeItem(this.userKey);

    storage.setItem(this.accessTokenKey, response.accessToken);
    storage.setItem(this.userKey, JSON.stringify(response.user));

    this.accessTokenSignal.set(response.accessToken);
    this.userSignal.set(response.user);
  }

  private getStoredUser(): AuthUser | null {
    const storedUser = localStorage.getItem(this.userKey) ?? sessionStorage.getItem(this.userKey);

    if (!storedUser) {
      return null;
    }

    try {
      return JSON.parse(storedUser) as AuthUser;
    } catch {
      return null;
    }
  }
}
