import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

interface Accommodation {
  id: string;
  name: string;
  location: string;
  amenities: string[];
  minGuests: number;
  maxGuests: number;
  images: string[];
  availability: string;
  pricePerNight: number;
}

@Component({
  selector: 'app-accommodations',
  templateUrl: './accommodations.component.html',
  styleUrls: ['./accommodations.component.css']
})
export class AccommodationsComponent implements OnInit {
  accommodations: Accommodation[] = [];

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    this.fetchAccommodations();
  }

  fetchAccommodations(): void {
    this.http.get<Accommodation[]>('http://localhost:8000/accommodations')
      .subscribe(
        (data: Accommodation[]) => {
          this.accommodations = data;
        },
        (error) => {
          console.error(error);
        }
      );
  }
}