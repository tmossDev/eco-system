import { Component, computed, inject } from '@angular/core';
import { ActivatedRoute, RouterOutlet } from '@angular/router';

export interface AuthLayoutContent {
  appName: string;
  appInitial: string;
  eyebrow: string;
  heading: string;
  description: string;
  features: string[];
}

const DEFAULT_AUTH_LAYOUT_CONTENT: AuthLayoutContent = {
  appName: 'Admin Web App',
  appInitial: 'A',
  eyebrow: 'Administration',
  heading: 'Manage your workspace with confidence.',
  description:
    'Securely access your admin tools, review operational data, and keep your workspace organised.',
  features: [
    'Centralised administration',
    'Role and access management',
    'Secure admin experience',
  ],
};

@Component({
  selector: 'app-auth-layout',
  imports: [RouterOutlet],
  templateUrl: './auth-layout.html',
  styleUrl: './auth-layout.scss',
})
export class AuthLayout {
  private readonly route = inject(ActivatedRoute);

  protected readonly content = computed<AuthLayoutContent>(() => ({
    ...DEFAULT_AUTH_LAYOUT_CONTENT,
    ...(this.route.snapshot.data['authLayout'] as
      | Partial<AuthLayoutContent>
      | undefined),
  }));
}
