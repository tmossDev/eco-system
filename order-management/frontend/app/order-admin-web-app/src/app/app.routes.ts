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
        appName: 'Order Admin Web App',
        appInitial: 'O',
        eyebrow: 'Order management',
        heading: 'Manage fulfillment with confidence.',
        description:
          'Securely access your admin tools, review customer orders, and keep your workspace organised.',
        features: [
          'Centralised order queue',
          'Status and fulfillment updates',
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
        path: 'orders',
        loadChildren: () =>
          import('@eco/admin-features/order').then((m) => m.ORDER_ROUTES),
      },
    ],
  },
  {
    path: '**',
    redirectTo: '',
  },
];
