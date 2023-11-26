import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Emitters } from './emitters/emitters';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  title = 'home';
  message = '';
  token: string | null = localStorage.getItem('token');

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    if (this.token) {
      this.getUserData(this.token);
    } else {
      this.message = 'No token found. Please log in.';
      Emitters.authEmitter.emit(false);
    }
  }

  getUserData(token: string): void {
    const httpOptions = {
      headers: new HttpHeaders({
        'Authorization': `Bearer ${token}`
      }),
      withCredentials: true
    };

    this.http.get<any>('http://localhost:8000/users/get', httpOptions).subscribe(
      (res: any) => {
        console.log(res);
        this.message = `Welcome ${res.first_name}`;
        Emitters.authEmitter.emit(true);
      },
      (err) => {
        console.log(err);
        this.message = 'You are not logged in';
        Emitters.authEmitter.emit(false);
      }
    );
  }
}
