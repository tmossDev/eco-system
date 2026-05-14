import type { Meta, StoryObj } from '@storybook/angular';

import { LoginForm } from './login-form';

const meta: Meta<LoginForm> = {
  title: 'Design System/Login Form',
  component: LoginForm,
  tags: ['autodocs'],
  argTypes: {
    eyebrow: {
      control: 'text',
    },
    heading: {
      control: 'text',
    },
    description: {
      control: 'text',
    },
    emailPlaceholder: {
      control: 'text',
    },
    passwordPlaceholder: {
      control: 'text',
    },
    submitText: {
      control: 'text',
    },
    submittingText: {
      control: 'text',
    },
    isSubmitting: {
      control: 'boolean',
    },
    errorMessage: {
      control: 'text',
    },
    showForgotPasswordLink: {
      control: 'boolean',
    },
    forgotPasswordHref: {
      control: 'text',
    },
    loginSubmit: {
      action: 'loginSubmit',
    },
  },
  args: {
    eyebrow: 'Welcome back',
    heading: 'Sign in to your account',
    description: 'Enter your details below to access the admin dashboard.',
    emailPlaceholder: 'you@example.com',
    passwordPlaceholder: 'Enter your password',
    submitText: 'Sign in',
    submittingText: 'Signing in...',
    isSubmitting: false,
    errorMessage: '',
    showForgotPasswordLink: true,
    forgotPasswordHref: '/auth/forgot-password',
  },
  parameters: {
    layout: 'centered',
  },
  decorators: [
    (story) => ({
      ...story(),
      template: `
        <div style="width: min(28rem, calc(100vw - 2rem));">
          ${story().template}
        </div>
      `,
    }),
  ],
};

export default meta;
type Story = StoryObj<LoginForm>;

export const Default: Story = {};

export const WithError: Story = {
  args: {
    errorMessage: 'Unable to sign in. Please check your details and try again.',
  },
};

export const Submitting: Story = {
  args: {
    isSubmitting: true,
  },
};

export const WithoutForgotPasswordLink: Story = {
  args: {
    showForgotPasswordLink: false,
  },
};

export const CustomCopy: Story = {
  args: {
    eyebrow: 'Admin portal',
    heading: 'Access your workspace',
    description: 'Use your administrator credentials to continue.',
    submitText: 'Continue',
    submittingText: 'Checking...',
  },
};
