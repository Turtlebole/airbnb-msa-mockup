import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-profile-edit',
  templateUrl: './profile-edit.component.html',
  styleUrls: ['./profile-edit.component.css']
})
export class ProfileEditComponent implements OnInit {
  profile: any = {};

  constructor(
    private http: HttpClient,
    private route: ActivatedRoute
  ) { }

  ngOnInit(): void {
    const profileId = this.route.snapshot.paramMap.get('id');
    this.getProfileById(profileId);
  }

  getProfileById(id: string | null): void {
    if (id) {
      this.http.get<any>(`http://localhost:8088/profiles/${id}`).subscribe(
        (res: any) => {
          this.profile = res;
        },
        (error) => {
          console.error(error);
        }
      );
    }
  }

  updateProfile(): void {
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json'
      })
    };
    this.http.put<any>(`http://localhost:8088/profiles/${this.profile._id}`, this.profile, httpOptions).subscribe(
      (response: any) => {
        console.log('Profile updated successfully:', response);
      },
      (error) => {
        console.error('Error updating profile:', error);
      }
    );
  }
}
