import { Component, Input } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'call-to-action',
  imports: [],
  templateUrl: './call-to-action.html',
})
export class CallToAction {
  @Input() icon!: string;
  @Input() boxTitle!: string;
  @Input() url!: string;

  constructor(private router: Router) {}

  navigateTo(url: string): void {
    if (url && url.startsWith('/')) {
      this.router.navigate([url]);
    } else if (url) {
      window.open(url, '_blank');
    }
  }
}