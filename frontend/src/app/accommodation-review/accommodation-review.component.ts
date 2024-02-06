import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { ActivatedRoute, Router } from '@angular/router';
import { Emitters } from '../emitters/emitters';
import { UserService } from '../user.service';

@Component({
  selector: 'app-accommodation-review',
  templateUrl: './accommodation-review.component.html',
  styleUrls: ['./accommodation-review.component.css'],
})
export class AccommodationReviewComponent implements OnInit {
  user_id: string | null = null;
  //user_name: string | null = null;
  token: string | null = null;
  httpOptions: { headers: HttpHeaders } = { headers: new HttpHeaders() };
  form: FormGroup = this.formBuilder.group({
    rating: [1, [Validators.required, Validators.min(1), Validators.max(5)]],
    comment: ['', [Validators.required]],
  });

  review: { rating: number; comment: string } = { rating: 1, comment: '' };

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private formBuilder: FormBuilder,
    private http: HttpClient,
    private userService: UserService
  ) {}
  ngOnInit(): void {
    this.userService.user$.subscribe((user) => {
      if (user) {
        // Access user properties like user.value
        console.log(user.value);
      }
    });
    this.user_id = localStorage.getItem('user_id');
    //this.user_name = localStorage.getItem('user_name');
    this.token = localStorage.getItem('token');

    console.log('user_id:', this.user_id);
    //console.log('user_name:', this.user_name);
    console.log('token:', this.token);

    if (!this.user_id || !this.token) {
      this.router.navigate(['/login']);
    }

    if (this.token) {
      this.getUserData(this.token);
    } else {
      console.log('No token found. Please log in.');
      Emitters.authEmitter.emit(false);
    }
  }

  getUserData(token: string): void {
    this.httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + token,
      }),
    };
  }

  submitReview() {
    // Retrieve accommodation ID from route parameters
    const accommodationId = this.route.snapshot.paramMap.get('id');

    // Retrieve user information from localStorage
    const userId = localStorage.getItem('user_id');
    const userFirstName = localStorage.getItem('user_first_name');

    console.log('User ID:', userId);
    console.log('User First Name:', userFirstName);

    // Now you can use userFirstName as needed, for example, set it in the payload
    const payload = {
      User_ID: userId,
      user_name: userFirstName,
      accommodation_id: accommodationId,
      rating: this.form.value.rating,
      comment: this.form.value.comment,
    };

    console.log('Payload:', payload);

    // Make an HTTP POST request to the specified endpoint
    this.http
      .post('http://localhost:8010/reviews/create', payload, this.httpOptions)
      .subscribe(
        (response) => {
          // Handle success, e.g., show a success message
          console.log('Review submitted successfully:', response);
          this.router.navigate(['/profile']);
        },
        (error) => {
          // Log the complete error object
          console.error('Error submitting review:', error);

          // Log the specific error message if available
          console.error('Error message:', error.message);

          // Log the response body if available
          console.error('Response body:', error.error);

          // Handle error, e.g., show an error message
        }
      );
  }
}
