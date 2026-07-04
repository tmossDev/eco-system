import {
  AfterContentInit,
  ChangeDetectionStrategy,
  Component,
  ContentChild,
  Input,
  OnChanges,
  SimpleChanges,
  booleanAttribute,
} from '@angular/core';

import { DsCheckboxDirective } from './checkbox.directive';
import { DsFormControlDirective, DsInputDirective, DsSelectDirective } from './form-control.directive';

let nextFieldId = 0;

@Component({
  selector: 'primitive-form-field',
  standalone: true,
  imports: [],
  template: `
    <label class="primitive-form-field__label" [attr.for]="controlId">
      <span>{{ label }}</span>
      @if (required) {
        <span class="primitive-form-field__required" aria-hidden="true">*</span>
      }
    </label>

    <ng-content />

    @if (hint) {
      <p class="primitive-form-field__hint" [id]="hintId">{{ hint }}</p>
    }

    @if (error) {
      <p class="primitive-form-field__error" [id]="errorId" role="alert">{{ error }}</p>
    }
  `,
  styleUrl: './form.scss',
  host: {
    class: 'primitive-form-field',
    '[class.primitive-form-field--invalid]': 'isInvalid',
  },
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DsFormField implements AfterContentInit, OnChanges {
  @Input({ required: true })
  public label = '';

  @Input()
  public hint = '';

  @Input()
  public error = '';

  @Input({ transform: booleanAttribute })
  public required = false;

  @ContentChild(DsInputDirective)
  private readonly inputControl?: DsInputDirective;

  @ContentChild(DsSelectDirective)
  private readonly selectControl?: DsSelectDirective;

  @ContentChild(DsCheckboxDirective)
  private readonly checkboxControl?: DsCheckboxDirective;

  protected readonly fieldId = `primitive-form-field-${nextFieldId++}`;
  protected readonly hintId = `${this.fieldId}-hint`;
  protected readonly errorId = `${this.fieldId}-error`;

  protected get controlId(): string | null {
    return this.control?.id ?? null;
  }

  protected get isInvalid(): boolean {
    return Boolean(this.error || this.control?.isInvalid);
  }

  public ngAfterContentInit(): void {
    this.syncControlAria();
  }

  public ngOnChanges(_changes: SimpleChanges): void {
    this.syncControlAria();
  }

  private get control(): DsFormControlDirective | DsCheckboxDirective | undefined {
    return this.inputControl ?? this.selectControl ?? this.checkboxControl;
  }

  private syncControlAria(): void {
    const control = this.control;

    if (!control) {
      return;
    }

    control.describedBy = [this.hint ? this.hintId : '', this.error ? this.errorId : '']
      .filter(Boolean)
      .join(' ') || null;

    control.invalid = this.error ? true : null;
    control.ariaRequired = this.required;
  }
}
