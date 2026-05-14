import { Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { Icon, IconName } from '@ds/icon';

export type NavigationBarVariant = 'sidebar' | 'header';

export interface NavigationItem {
  text: string;
  route: string;
  icon: IconName;
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

  @Output()
  public variantChange = new EventEmitter<NavigationBarVariant>();

  protected transitioningTo?: NavigationBarVariant;

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
}
