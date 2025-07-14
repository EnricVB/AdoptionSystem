import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'MobileFooter',
  imports: [CommonModule],
  templateUrl: './mobile-footer.html',
})
export class MobileFooter {
  activeIndex = 0;
  buttons = [
    { icon: 'fa-paw' },
    { icon: 'fa-hands-helping' },
    { icon: 'fa-hand-holding-usd' },
    { icon: 'fa-info-circle' },
  ];

  setActive(index: number) {
    this.activeIndex = index;
  }
}
