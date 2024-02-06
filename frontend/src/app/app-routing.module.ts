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
import { AccommodationViewComponent } from './accommodation-view/accommodation-view.component';
import { AccommodationReviewComponent } from './accommodation-review/accommodation-review.component';
import { HostReviewComponent } from './host-review/host-review.component';
import { UpdateReviewComponent } from './update-review/update-review.component';
import { HostReviewUpdateComponent } from './hostreview-update/hostreview-update.component';
const routes: Routes = [
  { path: 'login', component: LoginPageComponent },
  { path: 'sign-up', component: RegisterPageComponent },
  { path: 'accommodations', component: AccommodationsComponent },
  { path: 'accommodationsCreate', component: AccommodationCreateComponent },
  { path: 'edit/:id', component: AccommodationUpdateComponent },
  { path: 'profile', component: ProfileComponent },
  { path: 'reservation/:id', component: ReservationsComponent },
  { path: 'accommodation/:id', component: AccommodationViewComponent },
  { path: 'review/accommodation/:id', component: AccommodationReviewComponent },
  { path: 'review/host/:id', component: HostReviewComponent },
  { path: 'update-review/:reviewId', component: UpdateReviewComponent },
  { path: 'hostupdate-review/:reviewId', component: HostReviewUpdateComponent },
];
@NgModule({
  imports: [CommonModule, RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
