import { Component, EventEmitter } from '@angular/core';
import { Emitters } from '../emitters/emitters';
import { HttpClient } from '@angular/common/http';

import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';



@Component({
  selector: 'app-nav-bar',
  templateUrl: './nav-bar.component.html',
  styleUrls: ['./nav-bar.component.css']
})
export class NavBarComponent {
authenticated=false;
token=localStorage.getItem('token')
constructor(private http:HttpClient){
}
ngOnInit():void{
Emitters.authEmitter.subscribe(
  (auth:boolean)=>{

  }
);

}
logout():void{
this.http.post('http://localhost:8000/users/logout', this.token,{withCredentials:true}).subscribe(()=>{
  this.authenticated=false
})
}
}
