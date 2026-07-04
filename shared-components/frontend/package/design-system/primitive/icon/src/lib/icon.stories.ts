import type { Meta, StoryObj } from '@storybook/angular';

import { Icon } from './icon';
import { iconNames } from './icon.registry';

const meta: Meta<Icon> = {
  title: 'Design System/Primitives/Icon',
  component: Icon,
  tags: ['autodocs'],
  argTypes: {
    name: {
      control: 'select',
      options: iconNames,
    },
    size: {
      control: 'select',
      options: ['xs', 'small', 'medium', 'large', 'xl'],
    },
    color: {
      control: 'color',
    },
    ariaLabel: {
      control: 'text',
    },
    rotate: {
      control: 'select',
      options: [0, 90, 180, 270],
    },
  },
  args: {
    name: 'check',
    size: 'medium',
    color: '#333333',
    ariaLabel: 'Icon',
    rotate: 0,
  },
};

export default meta;
type Story = StoryObj<Icon>;

export const Playground: Story = {};

export const IconSet: Story = {
  render: () => ({
    props: {
      iconNames,
    },
    template: `
      <div style="display: grid; grid-template-columns: repeat(auto-fill, minmax(120px, 1fr)); gap: 16px;">
        <div
          *ngForOf="let iconName of iconNames"
          style="display: flex; flex-direction: column; align-items: center; gap: 8px; padding: 16px; border: 1px solid #e5e7eb; border-radius: 8px;"
        >
          <primitive-icon [name]="iconName" size="large"></primitive-icon>
          <span style="font-size: 12px;">{{ iconName }}</span>
        </div>
      </div>
    `,
  }),
};

export const Sizes: Story = {
  render: () => ({
    template: `
      <div style="display: flex; align-items: center; gap: 16px;">
        <primitive-icon name="check" size="xs"></primitive-icon>
        <primitive-icon name="check" size="small"></primitive-icon>
        <primitive-icon name="check" size="medium"></primitive-icon>
        <primitive-icon name="check" size="large"></primitive-icon>
        <primitive-icon name="check" size="xl"></primitive-icon>
      </div>
    `,
  }),
};

export const Colors: Story = {
  render: () => ({
    template: `
      <div style="display: flex; align-items: center; gap: 16px;">
        <primitive-icon name="info" size="large" color="#2563eb"></primitive-icon>
        <primitive-icon name="check" size="large" color="#16a34a"></primitive-icon>
        <primitive-icon name="warning" size="large" color="#f59e0b"></primitive-icon>
        <primitive-icon name="close" size="large" color="#dc2626"></primitive-icon>
      </div>
    `,
  }),
};
