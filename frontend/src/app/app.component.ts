import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Emitters } from './emitters/emitters';
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

  ngOnInit():void{
    this.http.get('http//localhost:8080/users',{withCredentials:true}).subscribe(
      res=>{
        console.log(res);
        this.message='Welcome ${res.name}';
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
