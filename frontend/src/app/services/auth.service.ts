import { Injectable } from "@angular/core";
import { User } from "@app/models";
import { BehaviorSubject } from "rxjs";
import { ApiService } from "./api.service";
import { CookieService } from "./cookie.service";

@Injectable({
  providedIn: 'root'
})
export class AuthService {
    private loggedInUser$ = new BehaviorSubject<User | null>(null);

    constructor(
        private apiService: ApiService, 
        private cookieService: CookieService) 
    {
        this.initializeLoggedInUser();
    }

    isLoggedIn(): boolean {
    return this.loggedInUser$ !== null;
    }

    setLoggedInUserId(user: User): void {
        this.loggedInUser$.next(user);
    }

    public get getLoggedInUser(): User | null {
        if (!this.loggedInUser$.getValue()) {
            this.initializeLoggedInUser();
        }

        return this.loggedInUser$.getValue();
    }

    private initializeLoggedInUser() {
        const sessionID = this.cookieService.getCookie('sessionID');

        if (sessionID) {
            this.apiService.getUserBySessionId(sessionID).subscribe({
                next: (data) => this.setLoggedInUserId(data.content),
                error: (err) => console.error('Failed to initialize user:', err),
            });
        }
    }
}