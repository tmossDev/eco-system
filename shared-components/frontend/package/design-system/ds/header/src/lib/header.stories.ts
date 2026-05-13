import type { Meta, StoryObj } from '@storybook/angular';

import { Header, type HeaderItem } from './header';

const defaultItems: HeaderItem[] = [
  {
    text: 'Home',
    route: '/',
    icon: 'home',
  },
  {
    text: 'Search',
    route: '/search',
    icon: 'search',
  },
  {
    text: 'Calendar',
    route: '/calendar',
    icon: 'calendar',
  },
  {
    text: 'Settings',
    route: '/settings',
    icon: 'settings',
  },
  {
    text: 'Profile',
    route: '/profile',
    icon: 'user',
  },
];

const meta: Meta<Header> = {
  title: 'Design System/Header',
  component: Header,
  tags: ['autodocs'],
  argTypes: {
    brandText: {
      control: 'text',
    },
    brandRoute: {
      control: 'text',
    },
    brandIcon: {
      control: 'select',
      options: [undefined, 'home', 'star', 'heart', 'settings'],
    },
    items: {
      control: 'object',
    },
  },
  args: {
    brandText: 'Eco System',
    brandRoute: '/',
    brandIcon: 'home',
    items: defaultItems,
  },
};

export default meta;
type Story = StoryObj<Header>;

export const Default: Story = {};

export const WithoutBrandIcon: Story = {
  args: {
    brandIcon: undefined,
  },
};

export const FewItems: Story = {
  args: {
    items: [
      {
        text: 'Home',
        route: '/',
        icon: 'home',
      },
      {
        text: 'Settings',
        route: '/settings',
        icon: 'settings',
      },
    ],
  },
};

export const ManyItems: Story = {
  args: {
    items: [
      ...defaultItems,
      {
        text: 'Downloads',
        route: '/downloads',
        icon: 'download',
      },
      {
        text: 'Upload',
        route: '/upload',
        icon: 'upload',
      },
    ],
  },
};
