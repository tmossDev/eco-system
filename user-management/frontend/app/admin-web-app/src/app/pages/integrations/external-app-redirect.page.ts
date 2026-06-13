import { Component, OnInit, computed, inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { CrossAppNavigationService } from '../../core/services/navigation/cross-app-navigation.service';

@Component({
  selector: 'app-external-app-redirect-page',
  template: `
    <section class="redirect-page">
      <p class="eyebrow">Opening workspace</p>
      <h1>{{ title() }}</h1>
      <p>Redirecting to the existing {{ title().toLowerCase() }} implementation.</p>
      <a [href]="targetUrl()">Open {{ title() }}</a>
    </section>
  `,
  styles: `
    .redirect-page {
      display: grid;
      align-content: center;
      min-height: 60dvh;
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
      letter-spacing: 0;
    }

    p {
      margin: 0.75rem 0 1rem;
      color: #56657f;
    }

    a {
      width: fit-content;
      border-radius: 999px;
      background: #2563eb;
      color: #ffffff;
      padding: 0.75rem 1rem;
      font-weight: 700;
      text-decoration: none;
    }
  `,
})
export class ExternalAppRedirectPage implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly crossAppNavigation = inject(CrossAppNavigationService);

  protected readonly title = computed(
    () => this.route.snapshot.data['title'] ?? 'workspace',
  );
  protected readonly targetUrl = computed(() =>
    this.crossAppNavigation.buildUrl(
      this.route.snapshot.data['appHostSegment'],
      this.route.snapshot.data['path'] ?? '/',
    ),
  );

  public ngOnInit(): void {
    window.location.assign(this.targetUrl());
  }
}
