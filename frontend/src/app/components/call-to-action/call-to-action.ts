import { Component, Input } from '@angular/core';

@Component({
  selector: 'call-to-action',
  imports: [],
  templateUrl: './call-to-action.html',
})
export class CallToAction {
  @Input() icon!: string;
  @Input() boxTitle!: string;
  @Input() url!: string;
}