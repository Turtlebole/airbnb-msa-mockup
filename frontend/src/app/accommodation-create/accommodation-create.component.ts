import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { FormBuilder, FormGroup } from '@angular/forms';
import { Router } from '@angular/router';
import { Emitters } from '../emitters/emitters';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';

@Component({
  selector: 'app-accommodation-create',
  templateUrl: './accommodation-create.component.html',
  styleUrls: ['./accommodation-create.component.css'],
})
export class AccommodationCreateComponent implements OnInit {
  form: FormGroup = new FormGroup({});
  token: string | null = localStorage.getItem('token');
  httpOptions: any; // Declare httpOptions at the class level

  constructor(
    private formBuilder: FormBuilder,
    private http: HttpClient,
    private router: Router
  ) {}
  
  ngOnInit(): void {
    this.form = this.formBuilder.group({
      name: '',
      location: '',
      amenities: [], // Initialize amenities as an empty array
      min_guests: '',
      max_guests: '',
      price_per_night: '',
    });

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
      withCredentials: true,
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
    console.log('Request Data:', requestData);

    this.http
      .post(
        'http://localhost:8000/api/accommodations/accommodations/create',
        requestData,
        this.httpOptions
      )
      .subscribe(
        () => this.router.navigate(['/accommodationsCreate']),
        (error) => console.error('Error:', error)
      );
  }
}