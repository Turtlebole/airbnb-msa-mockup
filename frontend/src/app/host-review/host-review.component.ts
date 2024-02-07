import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-host-review',
  templateUrl: './host-review.component.html',
  styleUrls: ['./host-review.component.css'],
})
export class HostReviewComponent implements OnInit {
  user_id: string | null = null;
  host_id: string | null = null;
  token: string | null = null;
  form: FormGroup = this.formBuilder.group({
    hostRating: [
      1,
      [Validators.required, Validators.min(1), Validators.max(5)],
    ],
    hostComment: ['', [Validators.required]],
  });

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private formBuilder: FormBuilder,
    private http: HttpClient
  ) {}

  ngOnInit(): void {
    this.user_id = localStorage.getItem('user_id');
    this.host_id = this.route.snapshot.paramMap.get('id');
    this.token = localStorage.getItem('token');

    if (!this.user_id || !this.host_id || !this.token) {
      this.router.navigate(['/login']);
    }
  }

  submitHostReview() {
    // Retrieve host ID from route parameters
    const hostId = this.route.snapshot.paramMap.get('id');

    // Use the stored user information from local storage
    const user_id = localStorage.getItem('user_id');
    const token = localStorage.getItem('token');

    // Extract user_name from the local storage
    const user_name = localStorage.getItem('user_first_name');

    // Get the current time for created_at
    const created_at = new Date().toISOString();

    // Now you can use user_id, hostId, user_name, created_at as needed
    const payload = {
      User_ID: user_id,
      Host_ID: hostId,
      host_rating: this.form.value.hostRating,
      host_comment: this.form.value.hostComment,
      user_name: user_name,
      created_at: created_at,
    };

    const httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + token,
      }),
    };

    this.http
      .post('http://localhost:8010/reviews/create-host', payload, httpOptions)
      .subscribe(
        (response) => {
          console.log('Host review submitted successfully:', response);
          // You can navigate to the profile or perform any other actions after a successful review
          this.router.navigate(['/profile']);
        },
        (error) => {
          console.error('Error submitting host review:', error);
        }
      );
  }
}
