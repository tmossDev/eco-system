import { Routes } from '@angular/router';

export const PRODUCTS_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./product-list/product-list.page').then(
        (m) => m.ProductListPage,
      ),
  },
  {
    path: ':id',
    loadComponent: () =>
      import('./product-detail/product-details.page').then(
        (m) => m.ProductDetailPage,
      ),
  },
  {
    path: ':id/edit',
    loadComponent: () =>
      import('./product-edit/product-edit.page').then(
        (m) => m.ProductEditPage,
      ),
  },
];
