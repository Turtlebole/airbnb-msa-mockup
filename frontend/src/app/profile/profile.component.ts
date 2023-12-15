import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  profiles: any[] = []; // Define the type of your profile data here
  token=localStorage.getItem('token')
  userID: string|undefined;

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    const httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + this.token,
      }),
      withCredentials: true,
    };
    
    this.http.get<any>(`http://localhost:8088/profiles/get-id`,httpOptions)
    .subscribe(
      (res: any) => {
        console.log("get id res: "+ res.user_id)
        this.userID=res.user_id
        this.profiles.push(res); // Add the retrieved profile to the profiles array
        this.http.get<any>(`http://localhost:8088/profiles/${this.userID}`)
        .subscribe(
          (res1: any) => {
            this.profiles.push(res1); // Add the retrieved profile to the profiles array
          },
          (error) => {
            console.error(error);
          }
        );
      },
      (error) => {
        console.error(error);
      }
    );
    // Fetch a specific user profile by ID
   
  }
}