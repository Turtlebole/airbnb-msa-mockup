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
  filteredAccommodations: any[] = [];
  locationFilter: string = '';
  guestsFilter: number | null = null;

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
        this.applyFilters();
      },
      (error) => {
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  applyFilters(): void {
    this.filteredAccommodations = this.accommodations.filter((accommodation) => {
      let passLocationFilter = true;
      let passGuestsFilter = true;

      if (this.locationFilter) {
        passLocationFilter = accommodation.location.toLowerCase().includes(this.locationFilter.toLowerCase());
      }

      if (this.guestsFilter !== null) {
        passGuestsFilter =
          accommodation.min_guests <= this.guestsFilter && accommodation.max_guests >= this.guestsFilter;
      }

      return passLocationFilter && passGuestsFilter;
    });
  }

  viewAccommodation(accommodationId: string): void {
    this.accommodationService.deleteAccommodation(accommodationId).subscribe(
      () => {
        this.loadAccommodations();
      },
      (error) => {
        console.error('Error deleting accommodation:', error);
      }
    );
  }

  editAccommodation(accommodationId: string): void {
    this.router.navigate(['/edit', accommodationId]);
  }


  makeReservation(accommodationId:string):void{
    this.router.navigate(['/reservation',accommodationId])
  }
}