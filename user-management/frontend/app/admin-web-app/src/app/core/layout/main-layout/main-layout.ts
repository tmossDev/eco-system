import { Component, computed, inject } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { NavigationBar, NavigationBarVariant, NavigationItem } from '@shared/navigation-bar';
import { AuthService } from '@eco/auth-features';

@Component({
  selector: 'app-main-layout',
  imports: [RouterOutlet, NavigationBar],
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.scss'
})
export class MainLayout {
  protected readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  public navigationVariant: NavigationBarVariant = 'sidebar';
  protected readonly currentUser = this.authService.currentUser;

  protected readonly profileRoute = computed(() => {
    const user = this.currentUser();

    return user ? ['/users', user.id] : null;
  });

  public readonly mainOptions: NavigationItem[] = [
    {
      text: 'Dashboard',
      route: '/',
      icon: 'home',
    },
    {
      text: 'Users',
      route: '/users',
      icon: 'user',
    },
    {
      text: 'Products',
      route: '/products',
      icon: 'box',
    },
    {
      text: 'Orders',
      route: '/orders',
      icon: 'calendar',
    },
    {
      text: 'Settings',
      route: '/settings',
      icon: 'settings',
    },
  ];

  protected logout(): void {
    this.authService.logout();
    void this.router.navigate(['/auth/login']);
  }
}
