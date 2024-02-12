import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-profile-edit',
  templateUrl: './profile-edit.component.html',
  styleUrls: ['./profile-edit.component.css']
})
export class ProfileEditComponent implements OnInit {
  profile: any = {};
  token = localStorage.getItem('token');
  userID: string | undefined;

  constructor(
    private http: HttpClient,
    private route: ActivatedRoute,
    private router: Router
  ) { }

  ngOnInit(): void {
    // Extract user ID from the route parameters
    this.route.params.subscribe(params => {
      const userId = params['id'];
      if (userId) {
        this.getProfileById(userId);
        this.userID = userId; // Set userID when receiving profile data
      } else {
        // Handle error or redirection if ID is not available
      }
    });
  }

  getProfileById(id: string): void {
    this.http.get<any>(`http://localhost:8088/profiles/${id}`).subscribe(
      (res: any) => {
        this.profile = res;
        console.log(this.profile.id);
        // Set userID property when receiving profile data
        this.userID = this.profile.id;
      },
      (error) => {
        console.error('Error fetching profile:', error);
      }
    );
  }

  updateProfile(): void {
    const httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + localStorage.getItem('token'), // Retrieve token from local storage
      })
    };
    // Use this.userID instead of id
    const id = this.userID;
    console.log(id);
    this.http.put<any>(`http://localhost:8088/profiles/${id}`, this.profile, httpOptions).subscribe(
      (res: any) => {
        console.log('Profile updated successfully:', res);
        // Redirect to profile page or handle success
        this.router.navigate(['/profile']);
      },
      (error) => {
        console.error('Error updating profile:', error);
        // Handle error response
      }
    );
  }
}