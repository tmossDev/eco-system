import { Routes } from '@angular/router';

export const ORDER_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./order-list.page').then((m) => m.OrderListPage),
  },
  {
    path: ':id',
    loadComponent: () =>
      import('./order-detail.page').then((m) => m.OrderDetailPage),
  },
];
