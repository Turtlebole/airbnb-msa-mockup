import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormControl, FormArray } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-accommodation-update',
  templateUrl: './accommodation-update.component.html',
  styleUrls: ['./accommodation-update.component.css']
})
export class AccommodationUpdateComponent implements OnInit {
  form: FormGroup;
  accommodationId: string;
  accommodation: any;
  amenitiesList: any[] = [
    { value: 'wifi', label: 'Wi-Fi' },
    { value: 'pool', label: 'Pool' },
    { value: 'gym', label: 'Gym' },
    { value: 'wellness', label: 'Wellness' },
    { value: 'restaurant', label: 'Restaurant' },
    { value: 'allInclusive', label: 'All Inclusive' },
  ];

  constructor(
    private formBuilder: FormBuilder,
    private http: HttpClient,
    private route: ActivatedRoute,
    private router: Router
  ) {
    this.form = this.formBuilder.group({
      name: '',
      location: '',
      min_guests: '',
      max_guests: '',
      amenities: this.formBuilder.array([]),
      availability_start: '',
      availability_end: '',
      price_per_night: '',
      price_on_weekends: ''
    });
    this.accommodationId = this.route.snapshot.params['id'];
  }

  ngOnInit(): void {
    this.getAccommodationDetails();
  }

  updateAmenities(): void {
    const amenitiesArray = this.form.get('amenities') as FormArray;

    if (amenitiesArray) {
      this.amenitiesList.forEach((amenity) => {
        const isChecked = this.accommodation.amenities.includes(amenity.value);
        if (isChecked) {
          amenitiesArray.push(this.formBuilder.control(amenity.value));
        }
      });
    }
  }

  isAmenityChecked(amenityValue: string): boolean {
    const amenitiesArray = this.form.get('amenities')?.value || [];
    return amenitiesArray.includes(amenityValue);
  }

  getAccommodationDetails(): void {
    this.http.get(`https://localhost/api/accommodations/accommodations/${this.accommodationId}`).subscribe(
      (response: any) => {
        this.accommodation = response;
        this.populateFormFields();
      },
      (error) => {
        console.error('Error fetching accommodation details:', error);
      }
    );
  }

  populateFormFields(): void {
    this.form.patchValue({
      name: this.accommodation.name,
      location: this.accommodation.location,
      min_guests: this.accommodation.min_guests,
      max_guests: this.accommodation.max_guests,
      availability_start: this.accommodation.availability_start,
      availability_end: this.accommodation.availability_end,
      price_per_night: this.accommodation.price_per_night,
      price_on_weekends: this.accommodation.price_on_weekends
    });

    const amenitiesArray = this.form.get('amenities') as FormArray;
    this.amenitiesList.forEach((amenity) => {
      const isChecked = this.accommodation.amenities.includes(amenity.value);
      if (isChecked) {
        amenitiesArray.push(this.formBuilder.control(amenity.value));
      }
    });
  }

  submit(): void {
    const updatedData = this.form.value;
    this.http
      .put(`https://localhost/api/accommodations/accommodations/${this.accommodationId}`, updatedData)
      .subscribe(
        () => {
          console.log('Accommodation updated');
          this.router.navigate(['/accommodations']);
        },
        (error) => {
          console.error('Error updating accommodation:', error);
        }
      );
  }
}