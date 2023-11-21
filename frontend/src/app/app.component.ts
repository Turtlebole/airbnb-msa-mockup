import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Emitters } from './emitters/emitters';
import { HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})

export class AppComponent {
  constructor(
    private http:HttpClient
  ){
  }
  title = 'home';
  message='';

  token=localStorage.getItem('token')

  ngOnInit():void{
    const httpOptions = {
      headers: new HttpHeaders({
        'Authorization': 'Bearer '+this.token
      }),
      withCredentials: true   // Assuming you need to send credentials along with the request
    };

    this.http.get('http://localhost:8000/users/get',httpOptions).subscribe(
      res=>{
        console.log(res);
        this.message='Welcome ${res.first_name}';

        Emitters.authEmitter.emit(true);
      },
      err=>{
        console.log(err);
        this.message='You are not logged in'
        Emitters.authEmitter.emit(true);
      }
    )
  }
}
