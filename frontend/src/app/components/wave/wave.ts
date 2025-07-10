import { Component, Input } from '@angular/core';

@Component({
  selector: 'wave',
  imports: [],
  templateUrl: './wave.html',
  host: {
    'class': 'block w-full',
  }
})
export class Wave {
  @Input() color: string = '';
}
