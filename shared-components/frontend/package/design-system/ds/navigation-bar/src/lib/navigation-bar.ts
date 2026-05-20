import { Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { Icon, IconName } from '@ds/icon';

export type NavigationBarVariant = 'sidebar' | 'header';

export interface NavigationItem {
  text: string;
  route: string;
  icon: IconName;
}

export interface NavigationUser {
  name: string;
  email: string;
  role?: string;
}

@Component({
  selector: 'ds-navigation-bar',
  imports: [RouterLink, RouterLinkActive, Icon],
  templateUrl: './navigation-bar.html',
  styleUrl: './navigation-bar.scss',
})
export class NavigationBar {
  private readonly transitionDuration = 260;

  @Input()
  public variant: NavigationBarVariant = 'sidebar';

  @Input()
  public title = 'Admin';

  @Input()
  public subtitle = 'Management';

  @Input()
  public items: NavigationItem[] = [];

  @Input()
  public user: NavigationUser | null = null;

  @Input()
  public profileRoute: string | any[] | null = null;

  @Output()
  public variantChange = new EventEmitter<NavigationBarVariant>();

  @Output()
  public logout = new EventEmitter<void>();

  protected transitioningTo?: NavigationBarVariant;
  protected isUserMenuOpen = false;

  protected get nextVariant(): NavigationBarVariant {
    return this.variant === 'sidebar' ? 'header' : 'sidebar';
  }

  protected get isTransitioning(): boolean {
    return this.transitioningTo !== undefined;
  }

  protected get toggleAriaLabel(): string {
    return `Switch navigation to ${this.nextVariant}`;
  }

  protected get toggleIcon(): IconName {
    return this.variant === 'sidebar' ? 'arrow-left' : 'arrow-down';
  }

  protected get userInitials(): string {
    const name = this.user?.name?.trim();

    if (!name) {
      return '?';
    }

    return name
      .split(/\s+/)
      .slice(0, 2)
      .map((namePart) => namePart[0])
      .join('')
      .toUpperCase();
  }

  protected get userMenuLabel(): string {
    return this.isUserMenuOpen ? 'Close user menu' : 'Open user menu';
  }

  protected toggleVariant(): void {
    if (this.isTransitioning) {
      return;
    }

    const nextVariant = this.nextVariant;
    this.transitioningTo = nextVariant;

    window.setTimeout(() => {
      this.variantChange.emit(nextVariant);
      this.transitioningTo = undefined;
    }, this.transitionDuration);
  }

  protected toggleUserMenu(): void {
    this.isUserMenuOpen = !this.isUserMenuOpen;
  }

  protected closeUserMenu(): void {
    this.isUserMenuOpen = false;
  }

  protected onLogout(): void {
    this.closeUserMenu();
    this.logout.emit();
  }
}
