import { Component, EventEmitter } from '@angular/core';
import { Emitters } from '../emitters/emitters';
import { HttpClient,HttpHeaders } from '@angular/common/http';

import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';
import { BehaviorSubject, Observable } from 'rxjs';
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
@Component({
  selector: 'app-nav-bar',
  templateUrl: './nav-bar.component.html',
  styleUrls: ['./nav-bar.component.css']
})
export class NavBarComponent {
authed=false
token=localStorage.getItem('token')
message=``
constructor(private http:HttpClient){
}
ngOnInit():void{
Emitters.authEmitter.subscribe(
  (auth:boolean)=>{

  }
);

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
    Authorization: 'Bearer ' + token,
  }),
  withCredentials: true,
};

this.http
  .get<any>('api/user/users/get', httpOptions)
  .subscribe(
    (res: any) => {
    //  console.log(res);
      this.message = `Welcome, ${res.user.first_name}`;
      Emitters.authEmitter.emit(true);
      this.authed=true
    },
    (err) => {
      console.log(err);
      this.message = 'You are not logged in';
      Emitters.authEmitter.emit(false);
    }
  );


}
logout():void{
  if (this.token) {
    this.getUserDataL(this.token);
    //console.log(this.token)
  } else {
    this.message = 'Please log in/Sign in.';
    Emitters.authEmitter.emit(false);
  }
  }
  getUserDataL(token: string): void {
    // const httpOptions = {
    //   headers: new HttpHeaders({
    //     Authorization: 'Bearer ' + token,
    //   }),
    //   withCredentials: true,
    // };
// this.http.post('http://localhost:8000/users/logout', httpOptions).subscribe(res=>{
//   this.authed=false
//   this.jaje=`aahahahah`
//   console.log(res)

// })

localStorage.removeItem('token')
location.reload();

}
}
