import { Routes } from '@angular/router';

import { StorefrontLayout } from './core/layout/storefront-layout';

export const routes: Routes = [
  {
    path: '',
    component: StorefrontLayout,
    children: [
      {
        path: '',
        loadComponent: () =>
          import('./pages/catalog/catalog.page').then((module) => module.CatalogPage),
      },
      {
        path: 'products/:id',
        loadComponent: () =>
          import('./pages/product-detail/product-detail.page').then(
            (module) => module.ProductDetailPage,
          ),
      },
      {
        path: 'cart',
        loadComponent: () =>
          import('./pages/cart/cart.page').then((module) => module.CartPage),
      },
      {
        path: 'auth/login',
        loadComponent: () =>
          import('./pages/auth/login.page').then((module) => module.LoginPage),
      },
      {
        path: 'auth/register',
        loadComponent: () =>
          import('./pages/auth/register.page').then(
            (module) => module.RegisterPage,
          ),
      },
    ],
  },
  { path: '**', redirectTo: '' },
];
