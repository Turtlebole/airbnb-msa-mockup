import { Component, OnInit } from '@angular/core';
import { AccommodationService } from 'src/app/accommodation.service';
import { Router } from '@angular/router';

interface AccommodationWithRating {
  id: string;
  name: string;
  location: string;
  min_guests: number;
  max_guests: number;
  amenities?: string[];
  averageRating?: number;
  availability: AvailabilityInterval[];
}

interface AvailabilityInterval {
  start: string;
  end: string;
  price_per_night: number;
  priceOnWeekends: number;
}

@Component({
  selector: 'app-accommodations',
  templateUrl: './accommodations.component.html',
  styleUrls: ['./accommodations.component.css'],
})
export class AccommodationsComponent implements OnInit {
  accommodations: AccommodationWithRating[] = [];
  filteredAccommodations: AccommodationWithRating[] = [];
  locationFilter: string = '';
  guestsFilter: number | null = null;
  averageRatings: { [key: string]: number | undefined } = {};
  ratingFilter: boolean | undefined = false;
  priceFilter: { from: number | null; to: number | null } = {
    from: null,
    to: null,
  };
  amenitySearch: string = '';

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
        const requests = data.map((accommodation) =>
          this.accommodationService
            .getAverageRating(accommodation.id)
            .toPromise()
            .then((averageRating) => {
              console.log(
                `Accommodation ID: ${accommodation.id}, Average Rating: ${averageRating}`
              );

              const rating = averageRating || 0;

              this.averageRatings[accommodation.id] = rating;
              return { ...accommodation, averageRating: rating };
            })
        );

        Promise.all(requests).then((accommodationsWithAvgRating) => {
          this.accommodations = accommodationsWithAvgRating;
          this.applyFilters();
        });
      },
      (error) => {
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  applyFilters(): void {
    this.filteredAccommodations = this.accommodations.filter(
      (accommodation) => {
        let passLocationFilter = true;
        let passGuestsFilter = true;
        let passRatingFilter = true;
        let passPriceFilter = true;
        let passAmenitySearch = true;

        // Existing location and guests filters
        if (this.locationFilter) {
          passLocationFilter = accommodation.location
            .toLowerCase()
            .includes(this.locationFilter.toLowerCase());
        }

        if (this.guestsFilter !== null) {
          passGuestsFilter =
            accommodation.min_guests <= this.guestsFilter &&
            accommodation.max_guests >= this.guestsFilter;
        }

        // Toggle filters
        if (this.ratingFilter !== undefined && this.ratingFilter) {
          const accommodationRating = this.averageRatings[accommodation.id];
          passRatingFilter =
            accommodationRating !== undefined && accommodationRating > 4.7;
        }

        if (
          this.priceFilter.from !== null &&
          this.priceFilter.to !== null &&
          this.priceFilter.from > 0
        ) {
          passPriceFilter =
            accommodation.availability[0].price_per_night >= this.priceFilter.from &&
            accommodation.availability[0].price_per_night <= this.priceFilter.to;
        }

        if (
          this.amenitySearch !== undefined &&
          this.amenitySearch.trim() !== ''
        ) {
          const amenities = accommodation.amenities || [];
          passAmenitySearch = amenities
            .join(' ')
            .toLowerCase()
            .includes(this.amenitySearch.toLowerCase());
        }

        return (
          passLocationFilter &&
          passGuestsFilter &&
          passRatingFilter &&
          passPriceFilter &&
          passAmenitySearch
        );
      }
    );
  }

  viewAccommodation(accommodationId: string): void {
    this.router.navigate(['/accommodation', accommodationId]);
  }

  editAccommodation(accommodationId: string): void {
    this.router.navigate(['/edit', accommodationId]);
  }

  makeReservation(accommodationId: string): void {
    this.router.navigate(['/reservation', accommodationId]);
  }
}
