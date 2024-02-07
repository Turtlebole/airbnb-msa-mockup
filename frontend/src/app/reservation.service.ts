import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ReservationService {
  private baseUrl = 'http://localhost:8002';

  constructor(private http: HttpClient) {}

  getAccommodationIdForReservation(reservationId: string): Observable<string> {
    return this.http.get<string>(
      `${this.baseUrl}/reservations/${reservationId}/accommodation-id`
    );
  }
}
