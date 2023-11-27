import { NgModule, createComponent } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginPageComponent } from './login-page/login-page.component';
import { CommonModule } from '@angular/common';
import { RegisterPageComponent } from './register-page/register-page.component';
import { AccommodationsComponent } from './accommodations/accommodations.component';
import { AccommodationCreateComponent } from './accommodation-create/accommodation-create.component';
import { ProfileComponent } from './profile/profile.component';
const routes: Routes = [
  { path: 'login', component: LoginPageComponent },
  { path: 'sign-up', component: RegisterPageComponent },
  { path: 'accommodations', component: AccommodationsComponent },
  { path: 'accommodationsCreate', component: AccommodationCreateComponent },
  { path: 'profile', component: ProfileComponent },
];
@NgModule({
  imports: [CommonModule, RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
