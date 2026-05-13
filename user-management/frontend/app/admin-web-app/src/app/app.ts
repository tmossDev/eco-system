import { Component, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import {Header, type HeaderItem} from '@ds/header';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, Header],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class App {
  protected readonly title = signal('admin-web-app');

  public defaultItems: HeaderItem[] = [
    {
      text: 'Home',
      route: '/',
      icon: 'home',
    },
    {
      text: 'Search',
      route: '/search',
      icon: 'search',
    },
    {
      text: 'Calendar',
      route: '/calendar',
      icon: 'calendar',
    },
    {
      text: 'Settings',
      route: '/settings',
      icon: 'settings',
    },
    {
      text: 'Profile',
      route: '/profile',
      icon: 'user',
    },
  ];

}
