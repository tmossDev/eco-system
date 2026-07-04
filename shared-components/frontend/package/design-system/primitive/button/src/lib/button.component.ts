import { CommonModule } from '@angular/common';
import { Component, Input, Output, EventEmitter } from '@angular/core';

export type ButtonVariant = 'primary' | 'secondary';
export type ButtonSize = 'small' | 'medium' | 'large';
export type ButtonType = 'button' | 'submit' | 'reset';

@Component({
  selector: 'primitive-button',
  standalone: true,
  imports: [CommonModule],
  template: ` <button
  [type]="type"
  [disabled]="disabled"
  (click)="onClick.emit($event)"
  [ngClass]="classes"
>
  <ng-content />
  @if (label) {
    {{ label }}
  }
</button>`,
  styleUrls: ['./button.component.scss'],
})
export class ButtonComponent {
  @Input()
  primary = false;

  @Input()
  variant: ButtonVariant = 'secondary';

  @Input()
  size: ButtonSize = 'medium';

  @Input()
  type: ButtonType = 'button';

  @Input()
  disabled = false;

  @Input()
  label = '';

  @Output()
  onClick = new EventEmitter<Event>();

  public get classes(): string[] {
    const mode = this.primary ? 'primary' : this.variant;

    return ['primitive-button', `primitive-button--${this.size}`, `primitive-button--${mode}`];
  }
}
