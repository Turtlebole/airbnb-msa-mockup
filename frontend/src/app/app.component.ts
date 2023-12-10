import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Emitters } from './emitters/emitters';
import { UserService } from './user.service';
import { NgForm } from '@angular/forms';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent implements OnInit {
  title = 'home';
  message = '';

  constructor(private http: HttpClient, private userService: UserService) {
    
  }

  ngOnInit(): void {
   
  }

  getUserData(token: string): void {
    const httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + token,
      }),
      withCredentials: true,
    };

    this.http
      .get<any>('https://localhost/api/user/users/get', httpOptions)
      .subscribe(
        (res: any) => {
          console.log(res);
          this.message = `Welcome ${res.first_name}`;
          Emitters.authEmitter.emit(true);
          this.userService.updateUser(res); // Notify other components about the user update
        },
        (err) => {
          console.log(err);
          this.message = 'You are not logged in';
          Emitters.authEmitter.emit(false);
        }
      );
  }

  
}
