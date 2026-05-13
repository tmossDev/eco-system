import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

import { iconRegistry, type IconName } from './icon.registry';

export type IconSize = 'xs' | 'small' | 'medium' | 'large' | 'xl';

@Component({
  selector: 'ds-icon',
  standalone: true,
  imports: [CommonModule],
  template: `
    <span
      class="ds-icon"
      [ngClass]="classes"
      [style.color]="color"
      [style.transform]="transform"
      [attr.aria-label]="ariaLabel"
      [attr.aria-hidden]="ariaLabel ? null : 'true'"
      role="img"
    >
      <svg [attr.viewBox]="viewBox" focusable="false">
        <path [attr.d]="path"></path>
      </svg>
    </span>
  `,
  styleUrls: ['./icon.css'],
})
export class Icon {
  @Input()
  name: IconName = 'check';

  @Input()
  size: IconSize = 'medium';

  @Input()
  color?: string;

  @Input()
  ariaLabel?: string;

  @Input()
  rotate: 0 | 90 | 180 | 270 = 0;

  protected get classes(): string[] {
    return ['ds-icon', `ds-icon--${this.size}`];
  }

  protected get path(): string {
    return iconRegistry[this.name].path;
  }

  protected get viewBox(): string {
    return iconRegistry[this.name].viewBox ?? '0 0 24 24';
  }

  protected get transform(): string | null {
    return this.rotate ? `rotate(${this.rotate}deg)` : null;
  }
}
