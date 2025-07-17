import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-pop-up',
  imports: [CommonModule],
  templateUrl: './pop-up.html',
})
export class PopUp {

  @Input() message: string = 'La mascota se ha registrado correctamente.';
  @Input() type: 'success' | 'error' = 'error';
  @Input() show: boolean = false;

  public start(message : string, duration : number = 3000) : void {
    this.show = true;
    this.message = message;

    setTimeout(() => {
      this.show = false;
    }, duration);
  }
}
