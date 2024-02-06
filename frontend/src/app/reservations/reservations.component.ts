import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-reservations',
  templateUrl: './reservations.component.html',
  styleUrls: ['./reservations.component.css']
})
export class ReservationsComponent implements OnInit {

  form: FormGroup;
  httpOptions: any;
  token = localStorage.getItem('token');
  constructor(
    private formBuilder: FormBuilder,
    private http: HttpClient,
    private router: Router,
    private route: ActivatedRoute

  ) {
    this.form = this.formBuilder.group({
      room_id: [''],
      guest_username: [''],
      guest_id: [''],
      checkin_date: [''],
      checkout_date: [''],
      number_of_guests: ['', [Validators.required, Validators.min(1)]]
    });
  }
  getUserData(token: string): void {
    this.httpOptions = {
      headers: new HttpHeaders({
        Authorization: 'Bearer ' + token,
      }),
    };
  }

  ngOnInit(): void {
    const user_id = localStorage.getItem('user_id');
    const user_first_name = localStorage.getItem('user_first_name');

    this.route.params.subscribe(params => {
      this.form.patchValue({
        room_id: params['id'] || '',
        guest_username: user_first_name,
        guest_id: user_id
      });
    });
    if (this.token) {
      this.getUserData(this.token);
    } else {
      console.log('No token found. Please log in.');
    }
  }


  onSubmit(): void {
    if (this.form.invalid) {
      const requestData = this.form.getRawValue();
      const totalPrice = this.calculatePrice(requestData);
      console.log('Total Price:', totalPrice);
  }

    const requestData = this.form.getRawValue();
    requestData.room_id = this.sanitizeInput(requestData.room_id);
    requestData.guest_username = this.sanitizeInput(requestData.guest_username);
    requestData.guest_id = this.sanitizeInput(requestData.guest_id);

    this.http
    .post<any>(
      'http://localhost:8002/reservations/by_guest/insert',
      requestData,
      this.httpOptions
      )
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

  sanitizeInput(input: any): any {
    if (typeof input === 'string') {
      const blockedCharactersPattern = /[<>"'`*/()\[\]?]/g;
      input = input.replace(blockedCharactersPattern, '');
    }
    return input;
  }

  calculatePrice(formValue: any): number {
    const checkinDate = new Date(formValue.checkin_date);
    const checkoutDate = new Date(formValue.checkout_date);
    const days = Math.floor((checkoutDate.getTime() - checkinDate.getTime()) / (1000 * 60 * 60 * 24));

    const weekendDays = this.countWeekendDays(checkinDate, checkoutDate);

    const guestNum = formValue.number_of_guests || 1;

    if (formValue.price_type === 'Whole') {
      return formValue.price_per_night * days + (formValue.price_on_weekends * weekendDays);
    } else {
      return guestNum * formValue.price_per_night * days + (formValue.price_on_weekends * guestNum * weekendDays);
    }
  }

  countWeekendDays(startDate: Date, endDate: Date): number {
    let weekendDays = 0;
    const currentDate = new Date(startDate);

    while (currentDate <= endDate) {
      const dayOfWeek = currentDate.getDay();
      if (dayOfWeek === 0 || dayOfWeek === 6) {
        weekendDays++;
      }
      currentDate.setDate(currentDate.getDate() + 1);
    }

    return weekendDays;
  }
}

