import { Component, signal } from '@angular/core';
import { provideRouter } from '@angular/router';
import type { Meta, StoryObj } from '@storybook/angular';
import { applicationConfig, moduleMetadata } from '@storybook/angular';

import {
  NavigationBar,
  NavigationBarVariant,
  type NavigationItem,
} from './navigation-bar';

const defaultItems: NavigationItem[] = [
  {
    text: 'Dashboard',
    route: '/',
    icon: 'home',
  },
  {
    text: 'Users',
    route: '/users',
    icon: 'user',
  },
  {
    text: 'Settings',
    route: '/settings',
    icon: 'settings',
  },
];

const manyItems: NavigationItem[] = [
  ...defaultItems,
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
    text: 'Downloads',
    route: '/downloads',
    icon: 'download',
  },
  {
    text: 'Upload',
    route: '/upload',
    icon: 'upload',
  },
];

@Component({
  selector: 'ds-navigation-bar-story-shell',
  imports: [NavigationBar],
  template: `
    <div class="story-shell" [class.story-shell--header]="variant() === 'header'">
      <ds-navigation-bar
        title="Admin"
        subtitle="Management"
        [variant]="variant()"
        [items]="items"
        [user]="user"
        [profileRoute]="['/users', user.id]"
        (variantChange)="variant.set($event)"
      />

      <main class="story-content">
        <p class="eyebrow">Interactive preview</p>
        <h1>{{ variant() === 'sidebar' ? 'Sidebar layout' : 'Header layout' }}</h1>
        <p>
          Use the arrow toggle in the navigation component to switch between
          sidebar and header navigation.
        </p>
      </main>
    </div>
  `,
  styles: [
    `
      .story-shell {
        display: grid;
        grid-template-columns: 17rem minmax(0, 1fr);
        min-height: 36rem;
        overflow: hidden;
        border: 1px solid #dbe3ef;
        border-radius: 1rem;
        background:
          radial-gradient(circle at top left, rgb(37 99 235 / 10%), transparent 28rem),
          #f8fafc;
      }

      .story-shell--header {
        display: block;
      }

      .story-content {
        min-width: 0;
        padding: 2rem;
        color: #172033;
      }

      .eyebrow {
        margin: 0 0 0.5rem;
        color: #56657f;
        font-size: 0.8rem;
        font-weight: 700;
        letter-spacing: 0.08em;
        text-transform: uppercase;
      }

      h1 {
        margin: 0;
        font-size: clamp(2rem, 4vw, 3rem);
        letter-spacing: -0.04em;
      }

      p {
        max-width: 42rem;
        color: #56657f;
        line-height: 1.6;
      }
    `,
  ],
})
class NavigationBarStoryShell {
  protected readonly variant = signal<NavigationBarVariant>('sidebar');
  protected readonly items = defaultItems;
  protected readonly user = {
    id: '1',
    name: 'Ada Lovelace',
    email: 'ada@example.com',
    role: 'Admin',
  };
}

const meta: Meta<NavigationBar> = {
  title: 'Design System/Navigation Bar',
  component: NavigationBar,
  decorators: [
    applicationConfig({
      providers: [provideRouter([])],
    }),
    moduleMetadata({
      imports: [NavigationBarStoryShell],
    }),
  ],
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'radio',
      options: ['sidebar', 'header'],
    },
    title: {
      control: 'text',
    },
    subtitle: {
      control: 'text',
    },
    items: {
      control: 'object',
    },
    user: {
      control: 'object',
    },
    profileRoute: {
      control: 'object',
    },
    variantChange: {
      action: 'variantChange',
    },
    logout: {
      action: 'logout',
    },
  },
  args: {
    variant: 'sidebar',
    title: 'Admin',
    subtitle: 'Management',
    items: defaultItems,
    user: {
      name: 'Ada Lovelace',
      email: 'ada@example.com',
      role: 'Admin',
    },
    profileRoute: '/users/1',
  },
};

export default meta;
type Story = StoryObj<NavigationBar>;

export const Sidebar: Story = {};

export const Header: Story = {
  args: {
    variant: 'header',
  },
};

export const FewItems: Story = {
  args: {
    items: [
      {
        text: 'Dashboard',
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
    items: manyItems,
  },
};

export const CustomBranding: Story = {
  args: {
    title: 'Eco System',
    subtitle: 'Workspace',
    variant: 'header',
    items: defaultItems,
  },
};

export const InteractiveToggle: StoryObj<NavigationBarStoryShell> = {
  render: () => ({
    props: {},
    template: `<ds-navigation-bar-story-shell />`,
  }),
};
