import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private apiUrl = 'http://localhost:8000'; // Replace with your backend API URL

  constructor(private http: HttpClient) {}

  // Add a method to send a request to the backend to become a host
  becomeHost(userId: string): Observable<any> {
    const url = `${this.apiUrl}/users/become-host/${userId}`;
    const headers = new HttpHeaders().set(
      'Authorization',
      'Bearer ' + this.getToken()
    );

    return this.http.post(url, {}, { headers });
  }

  // Add a method to get the authentication token from local storage
  private getToken(): string | null {
    return localStorage.getItem('token');
  }
}
