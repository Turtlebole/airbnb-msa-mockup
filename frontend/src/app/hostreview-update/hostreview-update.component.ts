import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup } from '@angular/forms';
import { AccommodationService } from '../accommodation.service';

@Component({
  selector: 'app-hostreview-update',
  templateUrl: './hostreview-update.component.html',
  styleUrls: ['./hostreview-update.component.css'],
})
export class HostReviewUpdateComponent implements OnInit {
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
      host_rating: [''], // add validators if needed
      host_comment: [''], // add validators if needed
      // Add other form controls as needed
    });

    this.route.params.subscribe((params) => {
      this.reviewId = params['reviewId'];

      // Fetch the review details using this.reviewId
      this.accommodationService.getReviewById(this.reviewId).subscribe(
        (reviewDetails) => {
          // Log the fetched data to the console
          console.log('Fetched review details:', reviewDetails);

          // Update form values with the fetched data
          this.form.patchValue({
            host_rating: reviewDetails.host_rating,
            host_comment: reviewDetails.host_comment,
            // Update other form controls as needed
          });
        },
        (error) => {
          console.error('Error fetching review details:', error);
        }
      );
    });
  }

  updateHostReview(): void {
    const hostRatingControl = this.form.get('host_rating');
    const hostCommentControl = this.form.get('host_comment');

    // Check if controls are not null or undefined
    if (hostRatingControl && hostCommentControl) {
      const updatedReview = {
        host_rating: hostRatingControl.value,
        host_comment: hostCommentControl.value,
        // Add other fields as needed
      };

      console.log('Updated Host Review:', updatedReview);

      this.accommodationService
        .updateReview(this.reviewId, updatedReview)
        .subscribe(
          () => {
            console.log('Successfully updated host review.');
            this.router.navigate(['/profile']);
          },
          (error) => {
            console.error('Error updating host review:', error);
          }
        );
    } else {
      console.error(
        'Host rating or comment form controls are null or undefined.'
      );
    }
  }
}
