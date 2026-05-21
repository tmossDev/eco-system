import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import {ApplicationSettings} from './settings.models';

@Injectable({
  providedIn: 'root',
})
export class SettingsService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/settings';

  public getSettings(): Observable<ApplicationSettings> {
    return this.http.get<ApplicationSettings>(this.baseUrl);
  }

  public updateSettings(settings: ApplicationSettings): Observable<ApplicationSettings> {
    return this.http.put<ApplicationSettings>(this.baseUrl, settings);
  }
}
