import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of, throwError } from 'rxjs';
import { tap, catchError, map } from 'rxjs/operators';
import { HttpErrorResponse } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class AccommodationService {
  private baseUrl = 'https://localhost';
  constructor(private http: HttpClient) {}

  getAllAccommodations(): Observable<any[]> {
    return this.http
      .get<any[]>(`${this.baseUrl}/api/accommodations/accommodations`)
      .pipe(
        tap((accommodations) => {
          console.log('Accommodations:', accommodations);
        }),
        catchError((error) => {
          console.error('Error fetching accommodations:', error);
          return of([]);
        })
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
  updateReview(reviewId: string, updatedReview: any): Observable<any> {
    if (!reviewId) {
      console.error('Invalid reviewId:', reviewId);
      return throwError('Invalid reviewId');
    }

    const url = `http://localhost:8010/reviews/update/${reviewId}`;

    try {
      console.log('Update Review Payload:', JSON.stringify(updatedReview)); // Log the payload
    } catch (error) {
      console.error('Error converting payload to JSON:', error);
    }

    // Update the field names in the payload
    const updatedPayload = {
      host_comment: updatedReview.host_comment,
      host_rating: updatedReview.host_rating,
      // Add other fields as needed
    };

    return this.http.put(url, updatedPayload).pipe(
      tap((response) => {
        console.log('Update Review Response:', response);
        if (response === null) {
          console.error('Update review response is null. Check server logs.');
        }
      }),
      catchError((error) => {
        console.error('Error updating review:', error);
        return throwError(error);
      })
    );
  }

  deleteReview(reviewId: string): Observable<any> {
    const url = `http://localhost:8010/reviews/delete/${reviewId}`;
    return this.http.delete(url);
  }
  getUserAccommodationReviews(userId: string): Observable<any[]> {
    return this.http.get<any[]>(`http://localhost:8010/reviews/user/${userId}`);
  }

  getUserHostReviews(userId: string): Observable<any[]> {
    return this.http.get<any[]>(
      `http://localhost:8010/reviews/user/${userId}/hosts`
    );
  }

  getReviewsByUser(userId: string): Observable<any[]> {
    return this.http.get<any[]>(`http://localhost:8010/reviews/user/${userId}`);
  }

  getReviewsByUserAndHosts(userId: string): Observable<any[]> {
    return this.http.get<any[]>(
      `http://localhost:8010/reviews/user/${userId}/hosts`
    );
  }
  getHostReviewsByAccommodation(accommodationId: string): Observable<any[]> {
    // Assuming your existing endpoint is for all reviews
    const url = `http://localhost:8010/reviews/all`;
    return this.http.get<any[]>(url).pipe(
      tap((hostReviews) => {
        console.log('Host Reviews:', hostReviews);
      }),
      catchError((error) => {
        console.error('Error fetching host reviews:', error);
        return of([]);
      })
    );
  }

  getAccommodationReviews(accommodationID: string): Observable<any> {
    const url = `http://localhost:8010/reviews/accommodation/${accommodationID}`;
    return this.http.get(url).pipe(
      tap((data) => console.log('Received data from API:', data)),
      catchError((error) => {
        console.error('Error from API:', error);
        throw error;
      })
    );
  }
  getHostReviewsByHostId(hostId: string): Observable<any> {
    return this.http.get<any>(
      `https://localhost/api/accommodations/accommodations`
    );
  }
  getReviewById(reviewId: string): Observable<any> {
    const url = `http://localhost:8010/reviews/${reviewId}`;

    return this.http.get<any>(url);
  }
}
