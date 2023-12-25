import { NgModule, createComponent } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginPageComponent } from './login-page/login-page.component';
import { CommonModule } from '@angular/common';
import { RegisterPageComponent } from './register-page/register-page.component';
import { AccommodationsComponent } from './accommodations/accommodations.component';
import { AccommodationCreateComponent } from './accommodation-create/accommodation-create.component';
import { ProfileComponent } from './profile/profile.component';
import { AccommodationUpdateComponent } from './accommodation-update/accommodation-update.component';
import { ReservationsComponent } from './reservations/reservations.component';
const routes: Routes = [
  { path: 'login', component: LoginPageComponent },
  { path: 'sign-up', component: RegisterPageComponent },
  { path: 'accommodations', component: AccommodationsComponent },
  { path: 'accommodationsCreate', component: AccommodationCreateComponent },
  { path: 'edit/:id', component: AccommodationUpdateComponent},
  { path: 'profile', component: ProfileComponent },
  {path: 'reservation/:id',component:ReservationsComponent}
];
@NgModule({
  imports: [CommonModule, RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
