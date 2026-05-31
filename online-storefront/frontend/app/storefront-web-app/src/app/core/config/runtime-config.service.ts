import { Injectable } from '@angular/core';

interface RuntimeConfig {
  BACKEND_API_URL?: string;
}

@Injectable({ providedIn: 'root' })
export class RuntimeConfigService {
  private config: RuntimeConfig = {};

  async load(): Promise<void> {
    try {
      const response = await fetch('/assets/config.json', { cache: 'no-store' });
      if (response.ok) {
        this.config = (await response.json()) as RuntimeConfig;
      }
    } catch {
      this.config = {};
    }
  }

  get backendApiUrl(): string {
    return this.config.BACKEND_API_URL?.replace(/\/$/, '') ?? '';
  }
}
