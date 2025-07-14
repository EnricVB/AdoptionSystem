import { Component, input, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Wave } from "../wave/wave";
import { NgIf } from '@angular/common';

@Component({
  selector: 'app-card',
  imports: [CommonModule, Wave],
  templateUrl: './card.html',
})
export class Card {
  @Input() title: string = 'Nombre';
  @Input() species: string = 'Especie';
  @Input() image: string = 'https://picsum.photos/200';
  @Input() isAdopted: boolean = false;
  @Input() description: string = 'Descripción del animal.'; 

  constructor() {
    // Initialization logic can go here if needed
  } 
}
