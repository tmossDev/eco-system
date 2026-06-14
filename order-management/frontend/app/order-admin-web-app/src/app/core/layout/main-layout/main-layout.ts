import { Component, computed, inject } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import {
  NavigationBar,
  NavigationBarVariant,
  NavigationItem,
} from '@ds/navigation-bar';
import { AuthService } from '@eco/auth-features';

@Component({
  selector: 'app-main-layout',
  imports: [RouterOutlet, NavigationBar],
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.scss',
})
export class MainLayout {
  protected readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  public navigationVariant: NavigationBarVariant = 'sidebar';
  protected readonly currentUser = this.authService.currentUser;

  protected readonly profileRoute = computed(() => {
    const user = this.currentUser();

    return user ? ['/'] : null;
  });

  public readonly mainOptions: NavigationItem[] = [
    {
      text: 'Dashboard',
      route: '/',
      icon: 'home',
    },
    {
      text: 'Orders',
      route: '/orders',
      icon: 'box',
    },
  ];

  protected logout(): void {
    this.authService.logout();
    void this.router.navigate(['/auth/login']);
  }
}
