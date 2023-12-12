import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginPageComponent } from './login-page/login-page.component';
import { RegisterPageComponent } from './register-page/register-page.component';
import { NavBarComponent } from './nav-bar/nav-bar.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { AccommodationsComponent } from './accommodations/accommodations.component';
import { AccommodationCreateComponent } from './accommodation-create/accommodation-create.component';
import { ProfileComponent } from './profile/profile.component';
import { UserService } from './user.service';

import { RECAPTCHA_SETTINGS, RecaptchaFormsModule, RecaptchaModule, RecaptchaSettings } from 'ng-recaptcha';

import { environment } from 'src/enviroments/enviroment';
import { AccommodationUpdateComponent } from './accommodation-update/accommodation-update.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginPageComponent,
    RegisterPageComponent,
    NavBarComponent,
    AccommodationsComponent,
    AccommodationCreateComponent,
    ProfileComponent,
    AccommodationUpdateComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    RecaptchaModule,
    RecaptchaFormsModule,
  ],
  providers: [UserService,  {
    provide: RECAPTCHA_SETTINGS,
    useValue: {
      siteKey: environment.recaptcha.siteKey,
    } as RecaptchaSettings,
  }],
  bootstrap: [AppComponent],
})
export class AppModule {}
