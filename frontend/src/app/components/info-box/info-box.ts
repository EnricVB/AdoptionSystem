import { Component, Input } from '@angular/core';

@Component({
  selector: 'info-box',
  imports: [],
  templateUrl: './info-box.html',
})
export class InfoBox {
  @Input() icon!: string;
  @Input() boxTitle!: string;
  @Input() description!: string;
  @Input() quantity!: number;
}