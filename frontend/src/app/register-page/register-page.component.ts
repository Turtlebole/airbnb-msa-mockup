import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';

@Component({
  selector: 'app-register-page',
  templateUrl: './register-page.component.html',
  styleUrls: ['./register-page.component.css'],
})
export class RegisterPageComponent implements OnInit {
  form: FormGroup = new FormGroup({});
  sanitizedContent!: SafeHtml;

  constructor(
    private formBuilder: FormBuilder,
    private http: HttpClient,
    private router: Router,
    private sanitizer: DomSanitizer
  ) {}

  sanitizeInput(input: any): any {
    if (typeof input === 'string') {
      const blockedCharactersPattern = /[<>"'`*/()\[\]?]/g;
      input = input.replace(blockedCharactersPattern, '');
    }
    return input;
  }

  ngOnInit(): void {
    this.form = this.formBuilder.group({
      first_name: '',
      last_name: '',
      password: ['',[Validators.pattern(/^(?=.*[A-Z])(?=.*[!@#$%^&*()_+.])[a-zA-Z0-9!@#$%^&*()_+.]{11,}$/)]],
      confirm_password: '',
      phone: ['',[Validators.pattern(/^(?:[0-9] ?){6,12}[0-9]$/)],],
      address: '',
      email: ['',[Validators.pattern(/^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/)],],
      user_type: 'User',
    });
  }

  submit(): void {
    if(this.form.controls['password'].value != this.form.controls['confirm_password'].value){
      alert('Passwords must match');
      return
    }
    if (this.form.invalid) {
      if (this.form.controls['password'].invalid) {
        alert('Password must be at least 11 characters, it needs to contain a capital letter and a special character.');
      }
      if (this.form.controls['phone'].invalid) {
        alert('Please enter a valid phone number.');
      }
      if (this.form.controls['email'].invalid) {
        alert('Please enter a valid email address.');
      }
      return;
    }
    const requestData = this.form.getRawValue();
    console.log('Request Data:', requestData);
    requestData.first_name = this.sanitizeInput(requestData.first_name);
    requestData.last_name = this.sanitizeInput(requestData.last_name);
    requestData.email = this.sanitizeInput(requestData.email);
    requestData.password = this.sanitizeInput(requestData.password);
    requestData.address = this.sanitizeInput(requestData.address);
    requestData.phone = this.sanitizeInput(requestData.phone);

    this.http
      .post('https://localhost/api/user/users/register', this.form.getRawValue())
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
