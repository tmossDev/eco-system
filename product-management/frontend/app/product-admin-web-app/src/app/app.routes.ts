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
        appName: 'Product Admin Web App',
        appInitial: 'A',
        eyebrow: 'Product management',
        heading: 'Manage products with confidence.',
        description:
          'Securely access your admin tools, review catalog items, and keep your workspace organised.',
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
        data: {
          loginEmailPlaceholder: 'admin@test.com',
        },
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
        path: 'products',
        loadChildren: () =>
          import('@eco/admin-features/product').then(
            (m) => m.PRODUCTS_ROUTES,
          ),
      },
      {
        path: 'promotions',
        loadComponent: () =>
          import('./pages/promotions/promotions-page/promotions-page').then(
            (m) => m.PromotionsPage,
          ),
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
