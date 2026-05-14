import { Component } from '@angular/core';
import { provideRouter, RouterLink, RouterOutlet } from '@angular/router';
import type { Meta, StoryObj } from '@storybook/angular';
import { applicationConfig, moduleMetadata } from '@storybook/angular';

import { MainLayout } from './main-layout';

@Component({
  selector: 'app-main-layout-story-dashboard',
  imports: [RouterLink],
  template: `
    <section class="story-page">
      <div class="page-header">
        <div>
          <p class="eyebrow">Overview</p>
          <h1>Dashboard</h1>
          <p class="description">
            Preview the main application layout with navigation and routed page
            content.
          </p>
        </div>

        <a routerLink="/users" class="primary-action">Manage users</a>
      </div>

      <div class="stats-grid">
        <article class="stat-card">
          <span>Total users</span>
          <strong>128</strong>
          <small>12 added this month</small>
        </article>

        <article class="stat-card">
          <span>Active users</span>
          <strong>96</strong>
          <small>75% active rate</small>
        </article>

        <article class="stat-card">
          <span>Pending invites</span>
          <strong>8</strong>
          <small>Awaiting response</small>
        </article>
      </div>

      <article class="panel">
        <h2>Recent activity</h2>
        <ul>
          <li>New user account created for Alex Morgan</li>
          <li>Priya Shah updated her profile details</li>
          <li>System settings were reviewed</li>
        </ul>
      </article>
    </section>
  `,
  styles: `
    .story-page {
      padding: 2rem;
      color: #172033;
    }

    .page-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 1rem;
      margin-bottom: 2rem;
    }

    .eyebrow {
      margin: 0 0 0.5rem;
      color: #56657f;
      font-size: 0.8rem;
      font-weight: 700;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }

    h1,
    h2 {
      margin: 0;
    }

    h1 {
      font-size: clamp(2rem, 4vw, 3rem);
      letter-spacing: -0.04em;
    }

    h2 {
      margin-bottom: 1rem;
      font-size: 1.15rem;
    }

    .description {
      max-width: 42rem;
      margin: 0.75rem 0 0;
      color: #56657f;
      line-height: 1.6;
    }

    .primary-action {
      display: inline-flex;
      align-items: center;
      border-radius: 999px;
      background: #2563eb;
      color: #ffffff;
      padding: 0.75rem 1rem;
      font-weight: 700;
      text-decoration: none;
      white-space: nowrap;
    }

    .stats-grid {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 1rem;
      margin-bottom: 1rem;
    }

    .stat-card,
    .panel {
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    .stat-card {
      padding: 1.25rem;
    }

    .stat-card span,
    .stat-card small,
    .panel li {
      color: #56657f;
    }

    .stat-card strong {
      display: block;
      margin: 0.5rem 0;
      font-size: 2rem;
    }

    .panel {
      padding: 1.25rem;
    }

    .panel ul {
      margin: 0;
      padding-left: 1.2rem;
      line-height: 1.8;
    }

    @media (max-width: 900px) {
      .page-header {
        display: grid;
      }

      .stats-grid {
        grid-template-columns: 1fr;
      }
    }
  `,
})
class MainLayoutStoryDashboard {}

@Component({
  selector: 'app-main-layout-story-users',
  template: `
    <section class="story-page">
      <div class="page-header">
        <div>
          <p class="eyebrow">User management</p>
          <h1>Users</h1>
          <p class="description">
            This route preview shows a simple users table inside the main layout.
          </p>
        </div>
      </div>

      <div class="table-card">
        <div class="table-header">
          <h2>All users</h2>
          <span>4 users</span>
        </div>

        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Email</th>
                <th>Role</th>
                <th>Status</th>
              </tr>
            </thead>

            <tbody>
              <tr>
                <td><strong>Alex Morgan</strong></td>
                <td>alex.morgan&#64;example.com</td>
                <td>Admin</td>
                <td><span class="status active">Active</span></td>
              </tr>
              <tr>
                <td><strong>Priya Shah</strong></td>
                <td>priya.shah&#64;example.com</td>
                <td>Manager</td>
                <td><span class="status active">Active</span></td>
              </tr>
              <tr>
                <td><strong>Jordan Lee</strong></td>
                <td>jordan.lee&#64;example.com</td>
                <td>User</td>
                <td><span class="status pending">Pending</span></td>
              </tr>
              <tr>
                <td><strong>Sam Taylor</strong></td>
                <td>sam.taylor&#64;example.com</td>
                <td>User</td>
                <td><span class="status suspended">Suspended</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </section>
  `,
  styles: `
    .story-page {
      padding: 2rem;
      color: #172033;
    }

    .page-header {
      margin-bottom: 2rem;
    }

    .eyebrow {
      margin: 0 0 0.5rem;
      color: #56657f;
      font-size: 0.8rem;
      font-weight: 700;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }

    h1,
    h2 {
      margin: 0;
    }

    h1 {
      font-size: clamp(2rem, 4vw, 3rem);
      letter-spacing: -0.04em;
    }

    h2 {
      font-size: 1.15rem;
    }

    .description {
      max-width: 42rem;
      margin: 0.75rem 0 0;
      color: #56657f;
      line-height: 1.6;
    }

    .table-card {
      overflow: hidden;
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    .table-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 1rem;
      padding: 1.25rem;
      border-bottom: 1px solid #dbe3ef;
    }

    .table-header span {
      color: #56657f;
    }

    .table-wrap {
      overflow-x: auto;
    }

    table {
      width: 100%;
      min-width: 720px;
      border-collapse: collapse;
    }

    th,
    td {
      padding: 1rem 1.25rem;
      border-bottom: 1px solid #eef2f7;
      text-align: left;
    }

    th {
      color: #56657f;
      font-size: 0.8rem;
      letter-spacing: 0.04em;
      text-transform: uppercase;
    }

    tbody tr:last-child td {
      border-bottom: 0;
    }

    .status {
      display: inline-flex;
      border-radius: 999px;
      padding: 0.3rem 0.65rem;
      font-size: 0.85rem;
      font-weight: 700;
    }

    .active {
      background: #dcfce7;
      color: #166534;
    }

    .pending {
      background: #fef3c7;
      color: #92400e;
    }

    .suspended {
      background: #fee2e2;
      color: #991b1b;
    }
  `,
})
class MainLayoutStoryUsers {}

@Component({
  selector: 'app-main-layout-story-settings',
  template: `
    <section class="story-page">
      <div class="page-header">
        <p class="eyebrow">Preferences</p>
        <h1>Settings</h1>
        <p class="description">
          Example settings page content rendered inside the main layout outlet.
        </p>
      </div>

      <form class="settings-card">
        <label>
          <span>Application name</span>
          <input type="text" value="Admin Web App" />
        </label>

        <label>
          <span>Default user role</span>
          <select>
            <option>User</option>
            <option>Manager</option>
            <option>Admin</option>
          </select>
        </label>

        <div class="toggle-row">
          <div>
            <strong>Email notifications</strong>
            <p>Send admin updates and user activity alerts.</p>
          </div>

          <input type="checkbox" checked />
        </div>
      </form>
    </section>
  `,
  styles: `
    .story-page {
      padding: 2rem;
      color: #172033;
    }

    .page-header {
      margin-bottom: 2rem;
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

    .description {
      max-width: 42rem;
      margin: 0.75rem 0 0;
      color: #56657f;
      line-height: 1.6;
    }

    .settings-card {
      display: grid;
      gap: 1.25rem;
      max-width: 44rem;
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      padding: 1.5rem;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    label {
      display: grid;
      gap: 0.4rem;
      font-weight: 700;
    }

    input,
    select {
      width: 100%;
      box-sizing: border-box;
      border: 1px solid #cbd5e1;
      border-radius: 0.75rem;
      padding: 0.8rem 0.9rem;
      color: #172033;
      font: inherit;
    }

    .toggle-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 1rem;
      border-top: 1px solid #eef2f7;
      padding-top: 1.25rem;
    }

    .toggle-row p {
      margin: 0.25rem 0 0;
      color: #56657f;
    }

    .toggle-row input {
      width: 1.25rem;
      height: 1.25rem;
      flex: 0 0 auto;
    }
  `,
})
class MainLayoutStorySettings {}

@Component({
  selector: 'app-main-layout-story-shell',
  imports: [RouterOutlet],
  template: `<router-outlet />`,
})
class MainLayoutStoryShell {}

const meta: Meta<MainLayoutStoryShell> = {
  title: 'App/Layout/Main Layout',
  component: MainLayoutStoryShell,
  decorators: [
    applicationConfig({
      providers: [
        provideRouter([
          {
            path: '',
            component: MainLayout,
            children: [
              {
                path: '',
                component: MainLayoutStoryDashboard,
              },
              {
                path: 'users',
                component: MainLayoutStoryUsers,
              },
              {
                path: 'settings',
                component: MainLayoutStorySettings,
              },
            ],
          },
        ]),
      ],
    }),
    moduleMetadata({
      imports: [
        MainLayout,
        MainLayoutStoryShell,
        MainLayoutStoryDashboard,
        MainLayoutStoryUsers,
        MainLayoutStorySettings,
      ],
    }),
  ],
  parameters: {
    layout: 'fullscreen',
  },
};

export default meta;
type Story = StoryObj<MainLayoutStoryShell>;

export const Dashboard: Story = {};

export const Users: Story = {
  render: () => ({
    template: `
      <app-main-layout>
        <app-main-layout-story-users />
      </app-main-layout>
    `,
  }),
};

export const Settings: Story = {
  render: () => ({
    template: `
      <app-main-layout>
        <app-main-layout-story-settings />
      </app-main-layout>
    `,
  }),
};
