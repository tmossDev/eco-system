import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import {CreateUserRequest, UpdateUserRequest, UserDetails, UserSummary} from './user.model';

@Injectable({
  providedIn: 'root',
})
export class UserService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/users';

  public getUsers(): Observable<UserSummary[]> {
    return this.http.get<UserSummary[]>(this.baseUrl);
  }

  public getUserById(id: string): Observable<UserDetails> {
    return this.http.get<UserDetails>(`${this.baseUrl}/${id}`);
  }

  public createUser(request: CreateUserRequest): Observable<UserDetails> {
    return this.http.post<UserDetails>(this.baseUrl, request);
  }

  public updateUser(id: string, request: UpdateUserRequest): Observable<UserDetails> {
    return this.http.put<UserDetails>(`${this.baseUrl}/${id}`, request);
  }

  public deleteUser(id: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${id}`);
  }
}
