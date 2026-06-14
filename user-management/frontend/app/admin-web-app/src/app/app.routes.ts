import { Routes } from '@angular/router';
import { authGuard, guestGuard } from '@eco/auth-features';

export const routes: Routes = [
  {
    path: 'auth',
    canActivate: [guestGuard],
    loadComponent: () =>
      import('@eco/auth-features').then((m) => m.AuthLayout),
    data: {
      authLayout: {
        appName: 'Admin Web App',
        appInitial: 'A',
        eyebrow: 'User management',
        heading: 'Manage users with confidence.',
        description:
          'Securely access your admin tools, review user accounts, and keep your workspace organised.',
        features: [
          'Centralised user administration',
          'Role and access management',
          'Secure admin experience',
        ],
      },
    },
    children: [
      {
        path: 'login',
        loadComponent: () =>
          import('@eco/auth-features').then((m) => m.LoginPage),
      },
      {
        path: 'forgot-password',
        loadComponent: () =>
          import('@eco/auth-features').then((m) => m.ForgotPasswordPage),
      },
      {
        path: '',
        pathMatch: 'full',
        redirectTo: 'login',
      },
    ],
  },
  {
    path: '',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./core/layout/main-layout/main-layout').then(
        (m) => m.MainLayout,
      ),
    children: [
      {
        path: '',
        loadComponent: () =>
          import('./pages/dashboard/dashboard.page').then(
            (m) => m.DashboardPage,
          ),
      },
      {
        path: 'users',
        loadChildren: () =>
          import('./pages/users/user.routes').then((m) => m.USERS_ROUTES),
      },
      {
        path: 'products',
        loadChildren: () =>
          import('@eco/admin-features/product').then(
            (m) => m.PRODUCTS_ROUTES,
          ),
      },
      {
        path: 'orders',
        loadChildren: () =>
          import('@eco/admin-features/order').then((m) => m.ORDER_ROUTES),
      },
      {
        path: 'settings',
        loadComponent: () =>
          import('./pages/settings/settings.page').then((m) => m.SettingsPage),
      },
    ],
  },
  {
    path: '**',
    redirectTo: '',
  },
];
