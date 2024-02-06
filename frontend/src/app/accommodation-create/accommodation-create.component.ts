import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { FormBuilder, FormGroup, FormArray } from '@angular/forms';
import { Router } from '@angular/router';
import { Emitters } from '../emitters/emitters';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { formatDate } from '@angular/common';

@Component({
  selector: 'app-accommodation-create',
  templateUrl: './accommodation-create.component.html',
  styleUrls: ['./accommodation-create.component.css'],
})
export class AccommodationCreateComponent implements OnInit {
  form: FormGroup;
  token: string | null;
  httpOptions: any;
  sanitizedContent!: SafeHtml;
  user_id: string | null;

  constructor(
    private formBuilder: FormBuilder,
    private http: HttpClient,
    private router: Router,
    private sanitizer: DomSanitizer
  ) {
    this.form = this.formBuilder.group({
      name: '',
      location: '',
      amenities: [],
      min_guests: '',
      max_guests: '',
      availability: this.formBuilder.array([]),
      availability_start:'',
      availability_end:'',
      price_per_night:'',
      price_on_weekends:'',
      price_type: ''
    });
    this.token = localStorage.getItem('token');
    this.user_id = localStorage.getItem('user_id');
  }

  sanitizeInput(input: any): any {
    if (typeof input === 'string') {
      const blockedCharactersPattern = /[<>"'`*/()\[\]?]/g;
      input = input.replace(blockedCharactersPattern, '');
    }
    return input;
  }

  ngOnInit(): void {
    if (this.user_id == null) {
      this.router.navigate(['/login']);
    }

    if (this.token) {
      this.getUserData(this.token);
    } else {
      console.log('No token found. Please log in.');
      Emitters.authEmitter.emit(false);
    }

    const currentDate = new Date();
    // const formattedDate = currentDate.toISOString().substring(0, 10);

    // this.form.patchValue({
    //   availability_start: formattedDate,
    //   availability_end: formattedDate
    // });
  }

  getUserData(token: string): void {
    this.httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + token,
      }),
    };
  }

  updateAmenities(event: any, amenity: string): void {
    const amenitiesArray = this.form.get('amenities')?.value || [];

    if (event.target.checked) {
      amenitiesArray.push(amenity);
    } else {
      const index = amenitiesArray.indexOf(amenity);
      if (index !== -1) {
        amenitiesArray.splice(index, 1);
      }
    }

    this.form.patchValue({ amenities: amenitiesArray });
  }

  submit(): void {
    const requestData = this.form.getRawValue();

    const availability = [{
      start: requestData.availability_start,
      end: requestData.availability_end,
      price_per_night: requestData.price_per_night,
      price_on_weekends: requestData.price_on_weekends
    }];

    requestData.availability = availability;

    delete requestData.availability_start;
    delete requestData.availability_end;
    delete requestData.price_per_night;
    delete requestData.price_on_weekends;

    requestData.name = this.sanitizeInput(requestData.name);
    requestData.location = this.sanitizeInput(requestData.location);
    requestData.amenities = this.sanitizeInput(requestData.amenities);
    requestData.max_guests = this.sanitizeInput(requestData.max_guests);
    requestData.min_guests = this.sanitizeInput(requestData.min_guests);
    requestData.price_type = this.sanitizeInput(requestData.price_type);
    console.log(requestData)
    console.log(availability)
    this.http
      .post(
        'https://localhost/api/accommodations/accommodations/create',
        requestData,
        this.httpOptions,
      )
      .subscribe(
        () => this.router.navigate(['/accommodations']),
        (error) => console.error('Error:', error)
      );
  }
}
