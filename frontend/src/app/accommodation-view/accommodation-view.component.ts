import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-accommodation-view',
  templateUrl: './accommodation-view.component.html',
  styleUrls: ['./accommodation-view.component.css']
})
export class AccommodationViewComponent implements OnInit {
  accommodationId: string = '';
  accommodation: any;

  constructor(
    private http: HttpClient,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit(): void {
    // Retrieve the accommodation ID from the route parameters
    this.accommodationId = this.route.snapshot.params['id'];
    // Fetch accommodation details based on the obtained ID
    this.getAccommodationDetails();
  }
  
  getAccommodationDetails(): void {
    this.http.get(`https://localhost/api/accommodations/accommodations/${this.accommodationId}`).subscribe(
      (response: any) => {
        this.accommodation = response;
      },
      (error) => {
        console.error('Error fetching accommodation details:', error);
      }
    );
  }
}