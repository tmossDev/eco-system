import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

import { Icon, type IconName } from '@ds/icon';

export interface HeaderItem {
  text: string;
  route: string;
  icon: IconName;
}

@Component({
  selector: 'ds-header',
  standalone: true,
  imports: [CommonModule, Icon],
  template: `
    <header class="ds-header">
      <a class="ds-header__brand" [href]="brandRoute" [attr.aria-label]="brandAriaLabel">
        <ng-container *ngIf="brandIcon">
          <ds-icon [name]="brandIcon" size="medium"></ds-icon>
        </ng-container>
        <span>{{ brandText }}</span>
      </a>

      <nav class="ds-header__nav" aria-label="Main navigation">
        <a
          *ngFor="let item of items"
          class="ds-header__nav-link"
          [href]="item.route"
          (click)="itemSelected.emit(item)"
        >
          <ds-icon [name]="item.icon" size="small"></ds-icon>
          <span>{{ item.text }}</span>
        </a>
      </nav>

      <button
        type="button"
        class="ds-header__menu-button"
        [attr.aria-expanded]="menuOpen"
        aria-controls="ds-header-mobile-menu"
        aria-label="Toggle navigation menu"
        (click)="toggleMenu()"
      >
        <ds-icon [name]="menuOpen ? 'close' : 'menu'" size="medium"></ds-icon>
      </button>
    </header>

    <nav
      id="ds-header-mobile-menu"
      class="ds-header__mobile-menu"
      [class.ds-header__mobile-menu--open]="menuOpen"
      aria-label="Mobile navigation"
    >
      <a
        *ngFor="let item of items"
        class="ds-header__mobile-link"
        [href]="item.route"
        (click)="selectMobileItem(item)"
      >
        <ds-icon [name]="item.icon" size="small"></ds-icon>
        <span>{{ item.text }}</span>
      </a>
    </nav>
  `,
  styleUrls: ['./header.css'],
})
export class Header {
  @Input()
  brandText = 'Design System';

  @Input()
  brandRoute = '/';

  @Input()
  brandIcon?: IconName;

  @Input()
  items: HeaderItem[] = [];

  @Output()
  itemSelected = new EventEmitter<HeaderItem>();

  protected menuOpen = false;

  protected get brandAriaLabel(): string {
    return `${this.brandText} home`;
  }

  protected toggleMenu(): void {
    this.menuOpen = !this.menuOpen;
  }

  protected selectMobileItem(item: HeaderItem): void {
    this.itemSelected.emit(item);
    this.menuOpen = false;
  }
}
