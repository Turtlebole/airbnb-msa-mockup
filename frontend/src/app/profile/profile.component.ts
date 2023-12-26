import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  profiles: any = {}; // Define the type of your profile data here
  reservations: any[] = []; // Array to hold reservations
  token = localStorage.getItem('token');
  userID: string | undefined;

  constructor(private http: HttpClient ,private router:Router) {}

  ngOnInit(): void {
    this.loadUserProfileAndReservations();
  }

  loadUserProfileAndReservations(): void {
    const httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + this.token,
      }),
      withCredentials: true,
    };

    this.http.get<any>(`http://localhost:8088/profiles/get-id`, httpOptions)
      .subscribe(
        (res: any) => {
          this.userID = res.user_id;
          if (this.userID) {
            this.http.get<any>(`http://localhost:8088/profiles/${this.userID}`, httpOptions)
              .subscribe(
                (res1: any) => {
                  this.profiles = res1; // Assign retrieved profile data
                },
                (error) => {
                  console.error(error);
                }
              );

            // Fetch reservations for this user
            this.http.get<any>(`http://localhost:8002/reservations/by_guest/${this.userID}`, httpOptions)
              .subscribe(
                (reservations: any) => {
                  this.reservations = reservations; // Assign retrieved reservations
                },
                (error) => {
                  console.error(error);
                }
              );
          }
        },
        (error) => {
          console.error(error);
        }
      );
  }

  cancelReservation(reservationId: string, roomId: string): void {
    const httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + this.token,
        'Content-Type': 'application/json', // Specify JSON content type
      }),
      withCredentials: true,
      body: JSON.stringify({ reservation_id: reservationId, room_id: roomId }), // Prepare the JSON body
    };

    this.http
      .request<any>('delete', `http://localhost:8002/reservations/cancel`, httpOptions)
      .subscribe(
        (response: any) => {
          console.log('Reservation canceled:', response);
          // Update the UI or handle the canceled reservation response here
          // For example, remove the canceled reservation from the 'reservations' array
          this.reservations = this.reservations.filter((res) => res.reservation_id !== reservationId);
        },
        (error) => {
          console.error('Cancellation failed:', error);
          // Handle error response if needed
        }
      );
  }
  deleteProfile():void{
    const httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + this.token,
        'Content-Type': 'application/json', // Specify JSON content type
      }),
      withCredentials: true,
    };
if(this.reservations.length==0){
    this.http
    .request<any>('delete', `http://localhost:8080/users/delete/${this.userID}`, httpOptions)
    .subscribe(
      (response: any) => {
        console.log(response)
        localStorage.removeItem('token')
        this.router.navigate(['/']).then(() => {
          window.location.reload();
        });

      },
      (error) => {
        console.error('Account delete failed:', error);
      }
    );
    }
    else{
      alert("You must cancel your reservations first")
    }
  }
}
