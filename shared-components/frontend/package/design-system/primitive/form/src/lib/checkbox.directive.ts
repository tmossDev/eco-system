import { Directive, HostBinding, Input, Optional, Self } from '@angular/core';
import { NgControl } from '@angular/forms';

let nextCheckboxId = 0;

@Directive({
  selector: 'input[type=checkbox][primitiveCheckbox]',
  standalone: true,
})
export class DsCheckboxDirective {
  @Input()
  @HostBinding('attr.id')
  public id = `primitive-checkbox-${nextCheckboxId++}`;

  @Input()
  @HostBinding('attr.aria-describedby')
  public describedBy: string | null = null;

  @Input()
  public invalid: boolean | null = null;

  public ariaRequired = false;

  @HostBinding('class.primitive-checkbox')
  protected readonly hasCheckboxClass = true;

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
