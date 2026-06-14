import { Injectable } from '@angular/core';

export interface RuntimeConfig {
  BACKEND_API_URL?: string;
  ENABLE_MOCK_API?: boolean;
}

@Injectable({
  providedIn: 'root',
})
export class RuntimeConfigService {
  private config: RuntimeConfig = {
    ENABLE_MOCK_API: true,
  };

  async load(): Promise<void> {
    try {
      const response = await fetch('/assets/config.json', {
        cache: 'no-store',
      });

      if (!response.ok) {
        return;
      }

      const runtimeConfig = (await response.json()) as RuntimeConfig;

      this.config = {
        ...this.config,
        ...runtimeConfig,
      };
    } catch {
      this.config = {
        ...this.config,
        ENABLE_MOCK_API: true,
      };
    }
  }

  get backendApiUrl(): string {
    return this.config.BACKEND_API_URL?.replace(/\/$/, '') ?? '';
  }

  get mockApiEnabled(): boolean {
    return this.config.ENABLE_MOCK_API ?? true;
  }
}
