import { Component, OnInit } from '@angular/core';
import { AccommodationService } from 'src/app/accommodation.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-accommodations',
  templateUrl: './accommodations.component.html',
  styleUrls: ['./accommodations.component.css'],
})
export class AccommodationsComponent implements OnInit {
  accommodations: any[] = [];

  constructor(
    private accommodationService: AccommodationService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.loadAccommodations();
  }

  loadAccommodations(): void {
    this.accommodationService.getAllAccommodations().subscribe(
      (data) => {
        this.accommodations = data;
      },
      (error) => {
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  viewAccommodation(accommodation: any): void {
    console.log('Viewing accommodation:', accommodation);
    const accommodationId = accommodation?.id;
    if (accommodationId) {
      this.router.navigate(['/accommodations', accommodationId]);
    } else {
      console.error('Invalid accommodationId:', accommodationId);
    }
  }
}
