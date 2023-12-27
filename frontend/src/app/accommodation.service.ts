import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { tap, catchError, map } from 'rxjs/operators';

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
  getReviewsByAccommodation(accommodationId: string): Observable<any> {
    return this.http.get<any>(
      `http://localhost:8010/reviews/accommodation/${accommodationId}`
    );
  }
  getAverageRating(accommodationId: string): Observable<number> {
    return this.http
      .get<{ averageRating: number }>(
        `http://localhost:8010/reviews/average-rating/${accommodationId}`
      )
      .pipe(
        tap((data) =>
          console.log(
            `Average Rating API Response for ID ${accommodationId}:`,
            data
          )
        ),
        map((data) => data.averageRating),
        catchError((error) => {
          console.error('Error fetching average rating:', error);
          return of(0);
        })
      );
  }
}
