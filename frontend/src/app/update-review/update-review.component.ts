import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup } from '@angular/forms';
import { AccommodationService } from '../accommodation.service';

@Component({
  selector: 'app-update-review',
  templateUrl: './update-review.component.html',
  styleUrls: ['./update-review.component.css'],
})
export class UpdateReviewComponent implements OnInit {
  reviewId!: string;
  form!: FormGroup;

  constructor(
    private route: ActivatedRoute,
    private fb: FormBuilder,
    private accommodationService: AccommodationService,
    private router: Router
  ) {}

  ngOnInit(): void {
    // Initialize your form controls here
    this.form = this.fb.group({
      rating: [''], // add validators if needed
      comment: [''], // add validators if needed
      // Add other form controls as needed
    });

    this.route.params.subscribe((params) => {
      this.reviewId = params['reviewId'];
      // Fetch the review details using this.reviewId
      this.accommodationService.getReviewById(this.reviewId).subscribe(
        (reviewDetails) => {
          // Update form values with the fetched data
          this.form.patchValue({
            rating: reviewDetails.rating,
            comment: reviewDetails.comment,
            // Update other form controls as needed
          });
        },
        (error) => {
          console.error('Error fetching review details:', error);
        }
      );
    });
  }

  updateReview(): void {
    const ratingControl = this.form.get('rating');
    const commentControl = this.form.get('comment');

    // Check if controls are not null or undefined
    if (ratingControl && commentControl) {
      const updatedReview = {
        rating: ratingControl.value,
        comment: commentControl.value,
        // Add other fields as needed
      };

      this.accommodationService
        .updateReview(this.reviewId, updatedReview)
        .subscribe(
          () => {
            // Successfully updated, navigate back to the profile
            this.router.navigate(['/profile']); // Adjust the route as needed
          },
          (error) => {
            console.error('Error updating review:', error);
          }
        );
    } else {
      console.error('Rating or comment form controls are null or undefined.');
    }
  }
}
