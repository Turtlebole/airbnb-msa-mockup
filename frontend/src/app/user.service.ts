import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class UserService {
  private userSource = new BehaviorSubject<any>(null);
  user$ = this.userSource.asObservable();

  updateUser(user: any): void {
    this.userSource.next(user);
  }
}
