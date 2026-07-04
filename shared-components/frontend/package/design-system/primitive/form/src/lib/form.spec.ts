import { Component } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DsInputDirective } from './form-control.directive';
import { DsFormField } from './form-field';
import { DsFormGroup } from './form-group';

@Component({
  standalone: true,
  imports: [DsFormField, DsInputDirective],
  template: `
    <primitive-form-field label="Email" hint="Use your work email" error="Email is required" required>
      <input primitiveInput type="email" />
    </primitive-form-field>
  `,
})
class FormFieldFixture {}

describe('DsFormField', () => {
  let fixture: ComponentFixture<FormFieldFixture>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FormFieldFixture],
    }).compileComponents();

    fixture = TestBed.createComponent(FormFieldFixture);
    fixture.detectChanges();
  });

  it('connects labels, descriptions, and invalid state to the projected control', () => {
    const nativeElement = fixture.nativeElement as HTMLElement;
    const label = nativeElement.querySelector('label');
    const input = nativeElement.querySelector('input');
    const hint = nativeElement.querySelector('.primitive-form-field__hint');
    const error = nativeElement.querySelector('.primitive-form-field__error');

    expect(input?.id).toBeTruthy();
    expect(label?.getAttribute('for')).toBe(input?.id);
    expect(input?.getAttribute('aria-invalid')).toBe('true');
    expect(input?.getAttribute('aria-required')).toBe('true');
    expect(input?.getAttribute('aria-describedby')).toBe(`${hint?.id} ${error?.id}`);
    expect(error?.getAttribute('role')).toBe('alert');
  });
});

describe('DsFormGroup', () => {
  it('should create', async () => {
    await TestBed.configureTestingModule({
      imports: [DsFormGroup],
    }).compileComponents();

    const group = TestBed.createComponent(DsFormGroup);

    expect(group.componentInstance).toBeTruthy();
  });
});
