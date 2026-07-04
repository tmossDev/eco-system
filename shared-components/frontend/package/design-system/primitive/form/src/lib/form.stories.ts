import type { Meta, StoryObj } from '@storybook/angular';
import { moduleMetadata } from '@storybook/angular';
import { ReactiveFormsModule, FormControl, Validators } from '@angular/forms';

import { DsCheckboxDirective } from './checkbox.directive';
import { DsFormField } from './form-field';
import { DsFormGroup } from './form-group';
import { DsInputDirective, DsSelectDirective } from './form-control.directive';

const meta: Meta = {
  title: 'Design System/Primitives/Form',
  decorators: [
    moduleMetadata({
      imports: [
        ReactiveFormsModule,
        DsFormField,
        DsFormGroup,
        DsInputDirective,
        DsSelectDirective,
        DsCheckboxDirective,
      ],
    }),
    (story) => ({
      ...story(),
      template: `
        <div style="width: min(32rem, calc(100vw - 2rem));">
          ${story().template}
        </div>
      `,
    }),
  ],
  parameters: {
    layout: 'centered',
  },
};

export default meta;
type Story = StoryObj;

export const TextInput: Story = {
  render: () => ({
    template: `
      <primitive-form-field label="Email address" hint="Use your work email." required>
        <input primitiveInput type="email" autocomplete="email" placeholder="you@example.com" />
      </primitive-form-field>
    `,
  }),
};

export const WithError: Story = {
  render: () => ({
    props: {
      email: new FormControl('', {
        nonNullable: true,
        validators: [Validators.required, Validators.email],
      }),
    },
    template: `
      <primitive-form-field label="Email address" error="Enter a valid email address." required>
        <input primitiveInput type="email" [formControl]="email" autocomplete="email" />
      </primitive-form-field>
    `,
  }),
};

export const SelectAndTextarea: Story = {
  render: () => ({
    template: `
      <primitive-form-group legend="Product details" description="Fields share the same accessible styling.">
        <primitive-form-field label="Category">
          <select primitiveSelect>
            <option>Apparel</option>
            <option>Accessories</option>
            <option>Home</option>
          </select>
        </primitive-form-field>

        <primitive-form-field label="Description" hint="Keep it concise for storefront cards.">
          <textarea primitiveInput placeholder="Describe the product"></textarea>
        </primitive-form-field>
      </primitive-form-group>
    `,
  }),
};

export const Checkbox: Story = {
  render: () => ({
    template: `
      <primitive-form-field label="Publish product" hint="Published products are visible in the storefront.">
        <input primitiveCheckbox type="checkbox" />
      </primitive-form-field>
    `,
  }),
};
