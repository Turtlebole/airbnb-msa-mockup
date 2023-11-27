import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'app-register-page',
  templateUrl: './register-page.component.html',
  styleUrls: ['./register-page.component.css'],
})
export class RegisterPageComponent implements OnInit {
  form: FormGroup = new FormGroup({});

  constructor(
    private formBuilder: FormBuilder,
    private http: HttpClient,
    private router: Router
  ) {}
  ngOnInit(): void {
    this.form = this.formBuilder.group({
      first_name: '',
      last_name: '',
      password: '',
      phone: '',
      address: '',
      email: '',
      user_type: 'UUser',
    });
  }

  submit(): void {
    console.log(this.form.getRawValue());

    this.http
      .post('http://localhost:8000/users/register', this.form.getRawValue())
      .subscribe(
        (response: any) => {
          if (
            response &&
            response['message'] === 'The chosen password is blacklisted'
          ) {
            // Show a pop-up with the blacklisted password message
            alert('The chosen password is blacklisted');
          } else {
            // Continue with your existing logic for other responses
            this.router.navigate(['/login']);
          }
        },
        (error) => {
          // Handle other errors if needed
          alert('An error occurred during registration');
        }
      );
  }
}
