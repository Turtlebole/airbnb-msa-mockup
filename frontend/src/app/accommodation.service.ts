import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AccommodationService {
  private baseUrl = 'https://localhost';
  constructor(private http: HttpClient) {}

  getAllAccommodations(): Observable<any[]> {
    return this.http.get<any[]>(
      `${this.baseUrl}/api/accommodations/accommodations`
    );
  }

  getAccommodationById(accommodationId: string): Observable<any> {
    return this.http.get<any>(
      `${this.baseUrl}/api/accommodations/accommodations/${accommodationId}`
    );
  }

  deleteAccommodation(accommodationId: string): Observable<any> {
    return this.http.delete<any>(
      `${this.baseUrl}/api/accommodations/accommodations/${accommodationId}`
    );
  }
}