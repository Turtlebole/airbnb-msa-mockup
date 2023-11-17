import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormBuilder,FormGroup } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.css']
})
export class LoginPageComponent {
  form: FormGroup = new FormGroup({});

  constructor(
    private formBuilder: FormBuilder,
    private http:HttpClient,
    private router:Router
    ){
  }
  ngOnInit():void{
    this.form = this.formBuilder.group({
      email:'',
      password:''

  });
  }
  submit():void{
    console.log(this.form.getRawValue());

    this.http.post('http://localhost:8000/users/login', this.form.getRawValue(),
    {withCredentials:true}).subscribe(res =>{
      this.router.navigate(['/'])
    console.log(res)});
}
}
