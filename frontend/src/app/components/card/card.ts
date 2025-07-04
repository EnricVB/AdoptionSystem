import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-card',
  imports: [CommonModule],
  templateUrl: './card.html',
  styleUrl: './card.css'
})
export class Card {
  @Input() title: string = 'Nombre';
  @Input() species: string = 'especie';

  constructor() {
    // Initialization logic can go here if needed
  } 
}
