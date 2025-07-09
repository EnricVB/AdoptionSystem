import { Injectable } from '@angular/core';
import { CookieService as NgxCookieService } from 'ngx-cookie-service';

@Injectable({
    providedIn: 'root'
})
export class CookieService {

    constructor(private cookieService: NgxCookieService) {}

    setCookie(name: string, value: string, days: number = 7, path: string = '/'): void {
        this.cookieService.set(name, value, days, path);
    }

    getCookie(name: string): string {
        return this.cookieService.get(name);
    }

    deleteCookie(name: string, path: string = '/'): void {
        this.cookieService.delete(name, path);
    }

    checkCookie(name: string): boolean {
        return this.cookieService.check(name);
    }
}