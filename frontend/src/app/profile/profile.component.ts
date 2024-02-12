import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';
import { AccommodationService } from 'src/app/accommodation.service';
import { ReservationService } from 'src/app/reservation.service';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css'],
})
export class ProfileComponent implements OnInit {
  profiles: any = {}; // Define the type of your profile data here
  reservations: any[] = []; // Array to hold reservations
  token = localStorage.getItem('token');
  userID: string | undefined;
  reviews: any[] = []; // Array to hold reviews
  updatedReview: any = {};
  accommodationReviews: any[] = [];
  hostReviews: any[] = [];

  constructor(
    private http: HttpClient,
    private router: Router,
    private accommodationService: AccommodationService,
    private reservationService: ReservationService
  ) {}

  ngOnInit(): void {
    this.loadUserProfileAndReservations();
    this.loadUserReviews();
  }

  loadUserProfileAndReservations(): void {
    const httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + this.token,
      }),
      withCredentials: true,
    };

    this.http
      .get<any>(`http://localhost:8088/profiles/get-id`, httpOptions)
      .subscribe(
        (res: any) => {
          this.userID = res.user_id;
          console.log('User ID:', this.userID); // Added console.log
          if (this.userID) {
            this.http
              .get<any>(
                `http://localhost:8088/profiles/${this.userID}`,
                httpOptions
              )
              .subscribe(
                (res1: any) => {
                  this.profiles = res1;
                  console.log('User Profiles:', this.profiles); // Added console.log
                },
                (error) => {
                  console.error('Profile Fetch Error:', error);
                }
              );

            // Fetch reservations for this user
            this.http
              .get<any>(
                `http://localhost:8002/reservations/by_guest/${this.userID}`
              )
              .subscribe(
                (reservations: any) => {
                  this.reservations = reservations;
                  console.log('User Reservations:', this.reservations); // Added console.log
                  // Call the new method to load user reviews
                  this.loadUserReviews();
                },
                (error) => {
                  console.error('Reservations Fetch Error:', error);
                }
              );
          }
        },
        (error) => {
          console.error('User ID Fetch Error:', error);
        }
      );
  }

  cancelReservation(reservationId: string, roomId: string): void {
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json', // Specify JSON content type
      }),
      body: JSON.stringify({ reservation_id: reservationId, room_id: roomId }), // Prepare the JSON body
    };

    this.http
      .request<any>(
        'delete',
        `http://localhost:8002/reservations/cancel`,
        httpOptions
      )
      .subscribe(
        (response: any) => {
          console.log('Reservation canceled:', response);
          // Update the UI or handle the canceled reservation response here
          // For example, remove the canceled reservation from the 'reservations' array
          this.reservations = this.reservations.filter(
            (res) => res.reservation_id !== reservationId
          );
        },
        (error) => {
          console.error('Cancellation failed:', error);
          // Handle error response if needed
        }
      );
  }
  deleteProfile(): void {
    const httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + this.token,
        'Content-Type': 'application/json', // Specify JSON content type
      }),
    };
    if (this.reservations.length == 0) {
      this.http
        .request<any>(
          'delete',
          `http://localhost:8080/users/delete/${this.userID}`,
          httpOptions
        )
        .subscribe(
          (response: any) => {
            console.log(response);
            localStorage.removeItem('token');
            this.router.navigate(['/']).then(() => {
              window.location.reload();
            });
          },
          (error) => {
            console.error('Account delete failed:', error);
            alert(error.error.error);
          }
        );
    } else {
      alert('You must cancel your reservations first');
    }
  }
  editProfile(): void {
    this.router.navigate(['/edit-profile',this.userID]);
  }

  //   updateReview(reviewId: string, updatedReview: any): void {
  //     this.accommodationService.updateReview(reviewId, updatedReview).subscribe(
  //       (response) => {
  //         console.log('Review updated successfully:', response);
  //         // You can update the UI or perform any other actions after a successful update
  //       },
  //       (error) => {
  //         console.error('Error updating review:', error);
  //       }
  //     );
  //   }
  reviewAccommodation(reservationId: string): void {
    console.log('Selected Reservation:', reservationId);
    console.log('All Reservations:', this.reservations);

    const selectedReservation = this.reservations.find(
      (reservation) => reservation.reservation_id === reservationId
    );

    if (selectedReservation) {
      const accommodationId = selectedReservation.room_id;

      // Navigate to the accommodation review page
      this.router.navigate(['review/accommodation/', accommodationId]);

      // Additional logging for debugging
      console.log('Navigating to accommodation review page:');
      console.log('Accommodation ID:', accommodationId);
    } else {
      console.error('Selected reservation not found in the list.');
    }
  }

  updateReview(reviewId: string, updatedReview: any): void {
    this.accommodationService.updateReview(reviewId, updatedReview).subscribe(
      (response) => {
        console.log('Review updated successfully:', response);
        // You can update the UI or perform any other actions after a successful update
      },
      (error) => {
        console.error('Error updating review:', error);
      }
    );
  }

  deleteReview(reviewId: string): void {
    this.accommodationService.deleteReview(reviewId).subscribe(
      () => {
        // Successfully deleted, now update the reviews
        this.loadUserReviews();
        console.log('Review deleted successfully.');
      },
      (error) => {
        console.error('Error deleting review:', error);
      }
    );
  }
  async reviewHost(reservationId: string): Promise<void> {
    console.log('All Reservations:', this.reservations);

    // Get the reservation details
    const reservation = this.reservations.find(
      (res) => res.reservation_id === reservationId
    );

    console.log('Selected Reservation:', reservation);

    if (!reservation) {
      console.error('Selected reservation not found in the list.');
      return;
    }

    // Fetch the accommodation details using the room_id from the reservation
    const accommodationId = reservation.room_id;
    console.log('Accommodation ID:', accommodationId);

    // Make a request to get accommodation details
    this.accommodationService.getAccommodationById(accommodationId).subscribe(
      (accommodation: any) => {
        console.log('Accommodation:', accommodation);

        if (!accommodation) {
          console.error(
            'Accommodation not found for reservation:',
            reservationId
          );
          return;
        }

        // Now you have the host_id, navigate to the host review page
        const hostId = accommodation.host_id;
        console.log('Host ID:', hostId);

        this.router.navigate(['/review/host', hostId]);
      },
      (error) => {
        console.error('Error fetching accommodation details:', error);
      }
    );
  }
  loadUserReviews(): void {
    if (this.userID) {
      // Fetch accommodation reviews for this user
      this.accommodationService.getReviewsByUser(this.userID).subscribe(
        (reviews: any) => {
          this.accommodationReviews = reviews;
          console.log('Accommodation Reviews:', this.accommodationReviews); // Added console.log
        },
        (error) => {
          console.error('Error fetching accommodation reviews:', error);
        }
      );

      // Fetch host reviews for this user
      this.accommodationService.getReviewsByUserAndHosts(this.userID).subscribe(
        (reviews: any) => {
          this.hostReviews = reviews;
          console.log('Host Reviews:', this.hostReviews); // Added console.log
        },
        (error) => {
          console.error('Error fetching host reviews:', error);
        }
      );
    }
  }
}
function editProfile() {
  throw new Error('Function not implemented.');
}

