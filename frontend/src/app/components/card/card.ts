import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Wave } from "../wave/wave";

@Component({
  selector: 'app-card',
  imports: [CommonModule, Wave],
  templateUrl: './card.html',
})
export class Card {
  @Input() title!: string;
  @Input() species!: string;
  @Input() image!: string;
  @Input() status!: string;
  @Input() description!: string; 
  @Input() age!: string;
  @Input() vaccinated!: string;

  constructor() {
  } 
}
