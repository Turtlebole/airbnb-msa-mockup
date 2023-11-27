import { Component, Input, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { UserService } from '../user.service';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css'],
})
export class ProfileComponent implements OnInit {
  @Input() user: any;

  constructor(private http: HttpClient, private userService: UserService) {}

  ngOnInit(): void {
    this.userService.user$.subscribe((user) => {
      // Ensure user is defined before updating the property
      if (user) {
        this.user = user;
        console.log(user.user.first_name);
      }
    });
  }

  becomeHost() {
    const userId = this.user.user_id;

    // Make a POST request to the backend to become a host
    this.http
      .post(`http://localhost:8000/users/${userId}/become-host`, {})
      .subscribe(
        (response: any) => {
          // Handle success
          this.user.user_type = 'Host';
          this.userService.updateUser(this.user); // Notify other components about the user update
        },
        (error) => {
          // Handle error
          console.error('Error becoming host:', error);
        }
      );
  }
}
