import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ReviewService {
  private apiUrl = 'https://localhost';

  constructor(private http: HttpClient) {}

  submitAccommodationReview(
    accommodationId: string,
    review: any
  ): Observable<any> {
    const url = `${this.apiUrl}/reviews/accommodation/${accommodationId}`;
    return this.http.post(url, review);
  }
}
