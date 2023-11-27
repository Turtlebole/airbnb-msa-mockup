import { Component, Input } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css'],
})
export class ProfileComponent {
  @Input() user: any;

  constructor(private http: HttpClient) {}

  becomeHost() {
    const userId = this.user.user_id;

    // Make a POST request to the backend to become a host
    this.http
      .post(`http://localhost:8000/users/${userId}/become-host`, {})
      .subscribe(
        (response) => {
          // Handle success
          this.user.user_type = 'Host';
        },
        (error) => {
          // Handle error
          console.error('Error becoming host:', error);
        }
      );
  }
}
