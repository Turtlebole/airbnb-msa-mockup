import { Component, Input, OnInit } from '@angular/core';
import { AuthService } from '../auth-service.service';
import { UserService } from '../user.service';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css'],
})
export class ProfileComponent implements OnInit {
  @Input() user: any;

  constructor(
    private http: HttpClient,
    private userService: UserService,
    private authService: AuthService // Inject the AuthService
  ) {}

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

    // Make a request to the AuthService to become a host
    this.authService.becomeHost(userId).subscribe(
      (response: any) => {
        // Handle success
        this.user.user_type = 'Host';
        this.userService.updateUser(this.user); // Notify other components about the user update
      },
      (error) => {
        // Handle error
        console.error('Error becoming host:', error);

        if (error.status === 500) {
          // Internal Server Error
          console.error('Internal Server Error. Please try again later.');
          // Display a user-friendly message to the user, e.g., using a toast or alert
        } else if (error.status === 404) {
          // Not Found
          console.error('User not found. Please check the user ID.');
          // Display a user-friendly message to the user
        } else {
          // Other error
          console.error('An unexpected error occurred.');
          // Display a user-friendly message to the user
        }
      }
    );
  }
}
