import { Routes } from '@angular/router';

export const USERS_ROUTES: Routes = [
{
path: '',
loadComponent: () =>
import('./user-list/user-list.page').then((m) => m.UserListPage),
},
{
path: ':id',
loadComponent: () =>
import('./user-detail/user-details.page').then((m) => m.UserDetailPage),
},
{
path: ':id/edit',
loadComponent: () =>
import('./user-edit/user-edit.page').then((m) => m.UserEditPage),
},
];
