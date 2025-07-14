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
  @Input() isAdopted!: boolean;
  @Input() description!: string; 

  constructor() {
  } 

  ngOnInit(): void {
    console.log('title:', this.title);
    console.log('species:', this.species);
    console.log('image:', this.image);
    console.log('isAdopted:', this.isAdopted);
    console.log('description:', this.description);
  }
}
