import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { AccommodationService } from 'src/app/accommodation.service';

@Component({
  selector: 'app-accommodation-details',
  templateUrl: './accommodation-details.component.html',
  styleUrls: ['./accommodation-details.component.css'],
})
export class AccommodationDetailsComponent implements OnInit {
  accommodationId: string | null = null;
  accommodationDetails: any | null = null;

  constructor(
    private route: ActivatedRoute,
    private accommodationService: AccommodationService
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe((params) => {
      this.accommodationId = params['id'];

      if (this.accommodationId) {
        this.loadAccommodationDetails();
      }
    });
  }

  loadAccommodationDetails(): void {
    this.accommodationService
      .getAccommodationById(this.accommodationId!)
      .subscribe(
        (data) => {
          this.accommodationDetails = data;
          console.log('Accommodation Details:', this.accommodationDetails);
        },
        (error) => {
          console.error('Error fetching accommodation details:', error);
        }
      );
  }
}
