import { Component } from '@angular/core';
import { CookieService } from '@app/services/cookie.service';

@Component({
  selector: 'PetHeader',
  imports: [],
  templateUrl: './header.html',
})
export class PetHeader {

  constructor(
      private cookieService: CookieService
  ) {
  }

  
  // Check if the user is already logged in by checking for a session ID cookie
  public isAlreadyLoggedIn(): boolean {
    return !!this.cookieService.getCookie('sessionID');
  }
}
