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
    { icon: 'fa-paw', label: 'Adopt' },
    { icon: 'fa-hands-helping', label: 'Help' },
    { icon: 'fa-hand-holding-usd', label: 'Donate' },
    { icon: 'fa-info-circle', label: 'Info' },
  ];

  setActive(index: number) {
    this.activeIndex = index;
  }
}
