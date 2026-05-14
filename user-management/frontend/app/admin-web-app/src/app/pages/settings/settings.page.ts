import { Component, OnInit, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { finalize } from 'rxjs';

import { ApplicationSettings } from '../../core/services/settings/settings.models';
import { SettingsService } from '../../core/services/settings/settings.service';

@Component({
  selector: 'app-settings-page',
  imports: [FormsModule],
  template: `
    <section class="settings-page">
      <div class="page-header">
        <p class="eyebrow">Preferences</p>
        <h1>Settings</h1>
        <p class="description">
          Manage application-level preferences for the admin experience.
        </p>
      </div>

      @if (isLoading()) {
        <p class="state-message">Loading settings...</p>
      }

      @if (message()) {
        <p class="success-message">{{ message() }}</p>
      }

      @if (errorMessage()) {
        <p class="error-message">{{ errorMessage() }}</p>
      }

      <form class="settings-card" (ngSubmit)="saveSettings()">
        <label>
          <span>Application name</span>
          <input
            type="text"
            name="applicationName"
            [ngModel]="settings().applicationName"
            (ngModelChange)="updateSetting('applicationName', $event)"
          />
        </label>

        <label>
          <span>Default user role</span>
          <select
            name="defaultRole"
            [ngModel]="settings().defaultRole"
            (ngModelChange)="updateSetting('defaultRole', $event)"
          >
            <option>Admin</option>
            <option>Manager</option>
            <option>User</option>
          </select>
        </label>

        <div class="toggle-row">
          <div>
            <strong>Email notifications</strong>
            <p>Send admin updates and user activity alerts.</p>
          </div>

          <input
            type="checkbox"
            name="emailNotifications"
            [ngModel]="settings().emailNotifications"
            (ngModelChange)="updateSetting('emailNotifications', $event)"
          />
        </div>

        <div class="toggle-row">
          <div>
            <strong>Require approval for new users</strong>
            <p>Newly invited users must be approved before access.</p>
          </div>

          <input
            type="checkbox"
            name="requireApproval"
            [ngModel]="settings().requireApproval"
            (ngModelChange)="updateSetting('requireApproval', $event)"
          />
        </div>

        <button type="submit" class="primary-action" [disabled]="isSaving()">
          {{ isSaving() ? 'Saving...' : 'Save settings' }}
        </button>
      </form>
    </section>
  `,
  styles: `
    .settings-page {
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

    .state-message,
    .success-message,
    .error-message {
      margin: 0 0 1rem;
      border-radius: 0.75rem;
      padding: 0.75rem 0.9rem;
      font-weight: 700;
    }

    .state-message {
      background: #eff6ff;
      color: #1d4ed8;
    }

    .success-message {
      background: #dcfce7;
      color: #166534;
    }

    .error-message {
      background: #fee2e2;
      color: #991b1b;
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

    input:focus,
    select:focus {
      border-color: #2563eb;
      outline: 3px solid #dbeafe;
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

    .primary-action {
      width: fit-content;
      border: 0;
      border-radius: 999px;
      background: #2563eb;
      color: #ffffff;
      padding: 0.75rem 1rem;
      font: inherit;
      font-weight: 700;
      cursor: pointer;
    }

    .primary-action:disabled {
      cursor: default;
      opacity: 0.65;
    }

    @media (max-width: 640px) {
      .settings-page {
        padding: 1rem;
      }

      .toggle-row {
        align-items: flex-start;
      }
    }
  `,
})
export class SettingsPage implements OnInit {
  private readonly settingsService = inject(SettingsService);

  protected readonly isLoading = signal(true);
  protected readonly isSaving = signal(false);
  protected readonly message = signal('');
  protected readonly errorMessage = signal('');

  protected readonly settings = signal<ApplicationSettings>({
    applicationName: '',
    defaultRole: 'User',
    emailNotifications: false,
    requireApproval: false,
  });

  public ngOnInit(): void {
    this.settingsService
      .getSettings()
      .pipe(
        finalize(() => {
          this.isLoading.set(false);
        }),
      )
      .subscribe({
        next: (settings) => {
          this.settings.set({ ...settings });
        },
        error: () => {
          this.errorMessage.set('Unable to load settings.');
        },
      });
  }

  protected updateSetting<Key extends keyof ApplicationSettings>(
    key: Key,
    value: ApplicationSettings[Key],
  ): void {
    this.settings.update((settings) => ({
      ...settings,
      [key]: value,
    }));
  }

  protected saveSettings(): void {
    if (this.isSaving()) {
      return;
    }

    this.isSaving.set(true);
    this.message.set('');
    this.errorMessage.set('');

    this.settingsService
      .updateSettings(this.settings())
      .pipe(
        finalize(() => {
          this.isSaving.set(false);
        }),
      )
      .subscribe({
        next: (settings) => {
          this.settings.set({ ...settings });
          this.message.set('Settings saved successfully.');
        },
        error: () => {
          this.errorMessage.set('Unable to save settings.');
        },
      });
  }
}
