import { Directive, HostBinding, Input, Optional, Self } from '@angular/core';
import { NgControl } from '@angular/forms';

let nextControlId = 0;

@Directive({
  selector: 'input[primitiveInput], textarea[primitiveInput]',
  standalone: true,
})
export class DsInputDirective {
  @Input()
  @HostBinding('attr.id')
  public id = `primitive-input-${nextControlId++}`;

  @Input()
  @HostBinding('attr.aria-describedby')
  public describedBy: string | null = null;

  @Input()
  public invalid: boolean | null = null;

  public ariaRequired = false;

  @HostBinding('class.primitive-form-control')
  protected readonly hasFormControlClass = true;

  @HostBinding('attr.aria-invalid')
  protected get ariaInvalid(): 'true' | null {
    return this.isInvalid ? 'true' : null;
  }

  @HostBinding('attr.aria-required')
  protected get requiredState(): 'true' | null {
    return this.ariaRequired ? 'true' : null;
  }

  public get isInvalid(): boolean {
    if (this.invalid !== null) {
      return this.invalid;
    }

    const control = this.ngControl?.control;

    return Boolean(control?.invalid && (control.touched || control.dirty));
  }

  public constructor(@Optional() @Self() private readonly ngControl: NgControl | null) {}
}

@Directive({
  selector: 'select[primitiveSelect]',
  standalone: true,
})
export class DsSelectDirective {
  @Input()
  @HostBinding('attr.id')
  public id = `primitive-select-${nextControlId++}`;

  @Input()
  @HostBinding('attr.aria-describedby')
  public describedBy: string | null = null;

  @Input()
  public invalid: boolean | null = null;

  public ariaRequired = false;

  @HostBinding('class.primitive-form-control')
  protected readonly hasFormControlClass = true;

  @HostBinding('class.primitive-form-control--select')
  protected readonly hasSelectClass = true;

  @HostBinding('attr.aria-invalid')
  protected get ariaInvalid(): 'true' | null {
    return this.isInvalid ? 'true' : null;
  }

  @HostBinding('attr.aria-required')
  protected get requiredState(): 'true' | null {
    return this.ariaRequired ? 'true' : null;
  }

  public get isInvalid(): boolean {
    if (this.invalid !== null) {
      return this.invalid;
    }

    const control = this.ngControl?.control;

    return Boolean(control?.invalid && (control.touched || control.dirty));
  }

  public constructor(@Optional() @Self() private readonly ngControl: NgControl | null) {}
}

export type DsFormControlDirective = DsInputDirective | DsSelectDirective;
