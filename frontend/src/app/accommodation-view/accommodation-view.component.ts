import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AccommodationService } from 'src/app/accommodation.service';

@Component({
  selector: 'app-accommodation-view',
  templateUrl: './accommodation-view.component.html',
  styleUrls: ['./accommodation-view.component.css'],
})
export class AccommodationViewComponent implements OnInit {
  accommodationId: string = '';
  accommodation: any;
  reviews: any[] = [];
  hostReviews: any[] = [];

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private accommodationService: AccommodationService
  ) {}

  ngOnInit(): void {
    // Retrieve the accommodation ID from the route parameters
    this.accommodationId = this.route.snapshot.params['id'];
    // Fetch accommodation details based on the obtained ID
    // this.getAccommodationDetails();
    // this.getReviewsByAccommodation();

    this.route.params.subscribe((params) => {
      this.accommodationId = params['id'];
      this.fetchAccommodationDetails();
      this.fetchReviews();
      this.fetchHostReviews();
    });
  }
  fetchHostReviews(): void {
    this.accommodationService
      .getHostReviewsByAccommodation(this.accommodationId)
      .subscribe(
        (hostReviews: any) => {
          console.log('Host Reviews:', hostReviews);
          this.hostReviews = hostReviews;
        },
        (error) => {
          console.error('Error fetching host reviews:', error);
        }
      );
  }

  fetchAccommodationDetails(): void {
    this.accommodationService
      .getAccommodationById(this.accommodationId)
      .subscribe(
        (response: any) => {
          this.accommodation = response;
          console.log(this.accommodation)
        },
        (error) => {
          console.error('Error fetching accommodation details:', error);
        }
      );
  }

  fetchReviews(): void {
    this.accommodationService
      .getAccommodationReviews(this.accommodationId)
      .subscribe(
        (reviews: any) => {
          this.reviews = reviews;
          // Filter host reviews based on conditions
          this.hostReviews = reviews.filter(
            (review: any) =>
              review.host_comment !== '' && review.host_rating !== 0
          );
        },
        (error) => {
          console.error('Error fetching reviews:', error);
        }
      );
  }
}
