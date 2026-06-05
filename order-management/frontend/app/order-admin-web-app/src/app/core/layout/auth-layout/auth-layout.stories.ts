import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { provideRouter, RouterLink, RouterOutlet } from '@angular/router';
import type { Meta, StoryObj } from '@storybook/angular';
import { applicationConfig, moduleMetadata } from '@storybook/angular';

import { AuthLayout } from './auth-layout';

@Component({
  selector: 'app-auth-layout-story-login',
  imports: [FormsModule, RouterLink],
  template: `
    <form class="login-form">
      <div class="form-header">
        <p class="eyebrow">Welcome back</p>
        <h2>Sign in to your account</h2>
        <p>Enter your details below to access the admin dashboard.</p>
      </div>

      <label>
        <span>Email address</span>
        <input
          type="email"
          name="email"
          autocomplete="email"
          placeholder="you@example.com"
          [(ngModel)]="form.email"
        />
      </label>

      <label>
        <span>Password</span>
        <input
          type="password"
          name="password"
          autocomplete="current-password"
          placeholder="Enter your password"
          [(ngModel)]="form.password"
        />
      </label>

      <div class="form-row">
        <label class="remember-me">
          <input
            type="checkbox"
            name="rememberMe"
            [(ngModel)]="form.rememberMe"
          />
          <span>Remember me</span>
        </label>

        <a routerLink="/forgot-password">Forgot password?</a>
      </div>

      <button type="button">Sign in</button>
    </form>
  `,
  styles: `
    .login-form {
      display: grid;
      gap: 1rem;
      border: 1px solid #dbe3ef;
      border-radius: 8px;
      background: #ffffff;
      padding: 1.5rem;
      box-shadow: 0 20px 60px rgb(15 23 42 / 10%);
    }

    .form-header {
      margin-bottom: 0.5rem;
    }

    .eyebrow {
      margin: 0 0 0.5rem;
      color: #2563eb;
      font-size: 0.8rem;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }

    h2 {
      margin: 0;
      color: #172033;
      font-size: 2rem;
      letter-spacing: 0;
    }

    .form-header p {
      margin: 0.75rem 0 0;
      color: #64748b;
      line-height: 1.6;
    }

    label {
      display: grid;
      gap: 0.4rem;
      color: #172033;
      font-weight: 700;
    }

    input {
      width: 100%;
      box-sizing: border-box;
      border: 1px solid #cbd5e1;
      border-radius: 0.75rem;
      padding: 0.8rem 0.9rem;
      color: #172033;
      font: inherit;
    }

    input:focus {
      border-color: #2563eb;
      outline: 3px solid #dbeafe;
    }

    .form-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 1rem;
    }

    .remember-me {
      display: flex;
      align-items: center;
      gap: 0.45rem;
      color: #64748b;
      font-weight: 600;
    }

    .remember-me input {
      width: 1rem;
      height: 1rem;
    }

    a {
      color: #2563eb;
      font-weight: 700;
      text-decoration: none;
      white-space: nowrap;
    }

    button {
      border: 0;
      border-radius: 999px;
      background: #2563eb;
      color: #ffffff;
      padding: 0.85rem 1rem;
      font: inherit;
      font-weight: 800;
      cursor: pointer;
    }

    @media (max-width: 420px) {
      .form-row {
        align-items: flex-start;
        flex-direction: column;
      }
    }
  `,
})
class AuthLayoutStoryLogin {
  protected readonly form = {
    email: '',
    password: '',
    rememberMe: true,
  };
}

@Component({
  selector: 'app-auth-layout-story-forgot-password',
  imports: [FormsModule, RouterLink],
  template: `
    <form class="forgot-password-form">
      <div class="form-header">
        <p class="eyebrow">Password recovery</p>
        <h2>Forgot your password?</h2>
        <p>
          Enter your email address and we’ll send reset instructions to your inbox.
        </p>
      </div>

      <label>
        <span>Email address</span>
        <input
          type="email"
          name="email"
          autocomplete="email"
          placeholder="you@example.com"
          [(ngModel)]="email"
        />
      </label>

      <button type="button">Send reset link</button>

      <a routerLink="/" class="back-link">Back to login</a>
    </form>
  `,
  styles: `
    .forgot-password-form {
      display: grid;
      gap: 1rem;
      border: 1px solid #dbe3ef;
      border-radius: 8px;
      background: #ffffff;
      padding: 1.5rem;
      box-shadow: 0 20px 60px rgb(15 23 42 / 10%);
    }

    .form-header {
      margin-bottom: 0.5rem;
    }

    .eyebrow {
      margin: 0 0 0.5rem;
      color: #2563eb;
      font-size: 0.8rem;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }

    h2 {
      margin: 0;
      color: #172033;
      font-size: 2rem;
      letter-spacing: 0;
    }

    .form-header p {
      margin: 0.75rem 0 0;
      color: #64748b;
      line-height: 1.6;
    }

    label {
      display: grid;
      gap: 0.4rem;
      color: #172033;
      font-weight: 700;
    }

    input {
      width: 100%;
      box-sizing: border-box;
      border: 1px solid #cbd5e1;
      border-radius: 0.75rem;
      padding: 0.8rem 0.9rem;
      color: #172033;
      font: inherit;
    }

    input:focus {
      border-color: #2563eb;
      outline: 3px solid #dbeafe;
    }

    button {
      border: 0;
      border-radius: 999px;
      background: #2563eb;
      color: #ffffff;
      padding: 0.85rem 1rem;
      font: inherit;
      font-weight: 800;
      cursor: pointer;
    }

    .back-link {
      justify-self: center;
      color: #2563eb;
      font-weight: 700;
      text-decoration: none;
    }
  `,
})
class AuthLayoutStoryForgotPassword {
  protected email = '';
}

@Component({
  selector: 'app-auth-layout-story-shell',
  imports: [RouterOutlet],
  template: `<router-outlet />`,
})
class AuthLayoutStoryShell {}

const meta: Meta<AuthLayoutStoryShell> = {
  title: 'App/Layout/Auth Layout',
  component: AuthLayoutStoryShell,
  decorators: [
    applicationConfig({
      providers: [
        provideRouter([
          {
            path: '',
            component: AuthLayout,
            children: [
              {
                path: '',
                component: AuthLayoutStoryLogin,
              },
              {
                path: 'forgot-password',
                component: AuthLayoutStoryForgotPassword,
              },
            ],
          },
        ]),
      ],
    }),
    moduleMetadata({
      imports: [
        AuthLayout,
        AuthLayoutStoryShell,
        AuthLayoutStoryLogin,
        AuthLayoutStoryForgotPassword,
      ],
    }),
  ],
  parameters: {
    layout: 'fullscreen',
  },
};

export default meta;
type Story = StoryObj<AuthLayoutStoryShell>;

export const Login: Story = {};

export const ForgotPassword: Story = {
  parameters: {
    docs: {
      description: {
        story:
          'Preview of the auth layout with a password recovery form in the outlet.',
      },
    },
  },
  render: () => ({
    template: `
      <main class="auth-layout-story">
        <section class="brand-panel-story" aria-label="Application information">
          <div class="brand-content-story">
            <a href="/" class="brand-mark-story" aria-label="Order Admin Web App home">
              <span class="brand-icon-story">O</span>
              <span>Order Admin Web App</span>
            </a>

            <div class="brand-message-story">
              <p class="eyebrow-story">Order management</p>
              <h1>Manage fulfillment with confidence.</h1>
              <p>
                Securely access your admin tools, review customer orders, and keep your
                workspace organised.
              </p>
            </div>

            <div class="feature-list-story">
              <div class="feature-item-story">
                <span>✓</span>
                <p>Centralised order queue</p>
              </div>

              <div class="feature-item-story">
                <span>✓</span>
                <p>Status and fulfillment updates</p>
              </div>

              <div class="feature-item-story">
                <span>✓</span>
                <p>Secure admin experience</p>
              </div>
            </div>
          </div>
        </section>

        <section class="auth-panel-story" aria-label="Authentication form">
          <div class="auth-card-story">
            <app-auth-layout-story-forgot-password />
          </div>
        </section>
      </main>
    `,
    styles: [
      `
        .auth-layout-story {
          display: grid;
          grid-template-columns: minmax(24rem, 0.9fr) minmax(0, 1.1fr);
          min-height: 100dvh;
          background:
            radial-gradient(circle at top right, rgb(37 99 235 / 12%), transparent 28rem),
            #f8fafc;
          color: #172033;
        }

        .brand-panel-story {
          display: flex;
          align-items: stretch;
          background:
            linear-gradient(135deg, rgb(15 23 42 / 92%), rgb(30 64 175 / 88%)),
            radial-gradient(circle at top left, rgb(96 165 250 / 40%), transparent 24rem);
          color: #ffffff;
          padding: 2rem;
        }

        .brand-content-story {
          display: flex;
          flex-direction: column;
          justify-content: space-between;
          width: 100%;
          border: 1px solid rgb(255 255 255 / 16%);
          border-radius: 8px;
          background: rgb(255 255 255 / 8%);
          padding: 2rem;
          backdrop-filter: blur(14px);
        }

        .brand-mark-story {
          display: inline-flex;
          align-items: center;
          gap: 0.75rem;
          width: fit-content;
          color: #ffffff;
          font-weight: 800;
          text-decoration: none;
        }

        .brand-icon-story {
          display: grid;
          place-items: center;
          width: 2.5rem;
          height: 2.5rem;
          border-radius: 0.8rem;
          background: #ffffff;
          color: #1d4ed8;
          font-weight: 900;
        }

        .brand-message-story {
          max-width: 34rem;
        }

        .eyebrow-story {
          margin: 0 0 0.75rem;
          color: #bfdbfe;
          font-size: 0.8rem;
          font-weight: 800;
          letter-spacing: 0.1em;
          text-transform: uppercase;
        }

        .brand-message-story h1 {
          margin: 0;
          font-size: clamp(2.5rem, 5vw, 4.5rem);
          line-height: 0.95;
          letter-spacing: 0;
        }

        .brand-message-story p:last-child {
          max-width: 30rem;
          margin: 1.25rem 0 0;
          color: #dbeafe;
          font-size: 1.05rem;
          line-height: 1.7;
        }

        .feature-list-story {
          display: grid;
          gap: 0.9rem;
        }

        .feature-item-story {
          display: flex;
          align-items: center;
          gap: 0.75rem;
        }

        .feature-item-story span {
          display: grid;
          place-items: center;
          width: 1.6rem;
          height: 1.6rem;
          border-radius: 999px;
          background: rgb(255 255 255 / 16%);
          color: #ffffff;
          font-weight: 900;
        }

        .feature-item-story p {
          margin: 0;
          color: #dbeafe;
          font-weight: 700;
        }

        .auth-panel-story {
          display: grid;
          place-items: center;
          padding: 2rem;
        }

        .auth-card-story {
          width: min(100%, 28rem);
        }

        @media (max-width: 960px) {
          .auth-layout-story {
            grid-template-columns: 1fr;
          }

          .brand-panel-story {
            display: none;
          }

          .auth-panel-story {
            min-height: 100dvh;
          }
        }

        @media (max-width: 520px) {
          .auth-panel-story {
            padding: 1rem;
          }
        }
      `,
    ],
  }),
};
