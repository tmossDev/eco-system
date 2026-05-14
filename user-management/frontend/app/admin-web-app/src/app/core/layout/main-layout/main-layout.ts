import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { NavigationBar, NavigationBarVariant, NavigationItem } from '@ds/navigation-bar';

@Component({
  selector: 'app-main-layout',
  imports: [RouterOutlet, NavigationBar],
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.scss'
})
export class MainLayout {
  public navigationVariant: NavigationBarVariant = 'sidebar';

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
      text: 'Settings',
      route: '/settings',
      icon: 'settings',
    },
  ];
}
