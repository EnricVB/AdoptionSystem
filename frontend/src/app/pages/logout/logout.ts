import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { CookieService } from '@app/services/cookie.service';

@Component({
  selector: 'app-logout',
  imports: [],
  templateUrl: './logout.html',
})
export class Logout {

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================

  // This component does not have any specific properties or methods.
  // It serves as a placeholder for the logout functionality.

  // ======================================
  // CONSTRUCTOR
  // ======================================

  constructor(
    private router: Router,
    private cookieService: CookieService
  ) {
    this.cookieService.deleteCookie('sessionID');
    this.router.navigate(['/dashboard']);
  }
}
