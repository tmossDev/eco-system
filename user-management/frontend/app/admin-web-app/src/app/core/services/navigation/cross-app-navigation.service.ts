import { Injectable, inject } from '@angular/core';

import { AuthUser } from '../auth/auth.models';
import { AuthService } from '../auth/auth.service';

@Injectable({
  providedIn: 'root',
})
export class CrossAppNavigationService {
  private readonly authService = inject(AuthService);

  public buildUrl(
    appHostSegment: string,
    path = '/',
    userOverride?: AuthUser,
  ): string {
    const currentLocation = window.location;
    const host = currentLocation.hostname.includes('admin-web-app')
      ? currentLocation.hostname.replace('admin-web-app', appHostSegment)
      : `eco-test.${appHostSegment}.com`;
    const url = new URL(path, `${currentLocation.protocol}//${host}`);

    this.authService.appendSessionHandoff(url, userOverride);

    return url.toString();
  }
}
