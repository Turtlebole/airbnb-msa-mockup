import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormBuilder,FormGroup,Validators } from '@angular/forms';
import { Router ,ActivatedRoute } from '@angular/router';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';

@Component({
  selector: 'app-reservations',
  templateUrl: './reservations.component.html',
  styleUrls: ['./reservations.component.css']
})
export class ReservationsComponent {

  form: FormGroup = new FormGroup({});
  token: string|undefined;
  id: string|undefined;
  checkin_date:string = ''
  checkout_date:string = ''
  
  constructor(
    private formBuilder: FormBuilder,
    private http:HttpClient,
    private router:Router,
    private route: ActivatedRoute
    
    ){
      this.token = undefined;
      this.id=undefined

  }
  
  convertDateFormat(inputDate: string) {
    const parts = inputDate.split('/');
    if (parts.length === 3) {
      const day = parseInt(parts[0], 10);
      const month = parseInt(parts[1], 10) - 1; // Month starts from 0 in JavaScript Date
      const year = parseInt(parts[2], 10);
      
      const date = new Date(year, month, day); // Create a Date object in the original format
      inputDate = date.toLocaleDateString('en-CA'); // Convert to yyyy-MM-dd format
    } else {
      console.error('Invalid date format!');
    }
    return inputDate
  }

  sanitizeInput(input: any): any {
    if (typeof input === 'string') {
      const blockedCharactersPattern = /[<>"'`*/()\[\]?]/g;
      input = input.replace(blockedCharactersPattern, '');
    }
    return input;
  }


  ngOnInit():void{
    const user_id=localStorage.getItem('user_id')
    const user_first_name=localStorage.getItem('user_first_name')
    this.form = this.formBuilder.group({
      room_id: [''],
      guest_username: [''],
      guest_id: [''], 
      checkin_date: [''],
      checkout_date: ['']
    });
    this.route.params.subscribe(params => {
      this.form.patchValue({
        room_id: params['id'] || '',
        guest_username: user_first_name,
        guest_id: user_id
      });
    });
  }

  onSubmit(): void {
    const requestData = this.form.getRawValue();
    requestData.room_id = this.sanitizeInput(requestData.room_id);
    requestData.guest_username = this.sanitizeInput(requestData.guest_username);
    requestData.guest_id = this.sanitizeInput(requestData.guest_id);
    requestData.checkin_date=this.convertDateFormat(this.checkin_date)
    requestData.checkout_date=this.convertDateFormat(this.checkout_date)

    this.http.post<any>('http://localhost:8002/reservations/by_guest/insert', this.form.getRawValue(), { withCredentials: true })
    .subscribe(
      (res: any) => {
        this.router.navigate(['/']).then(() => {
          window.location.reload();
        });
      },
      (error) => {
        console.error(error);
      }
    );
  }
  
}
