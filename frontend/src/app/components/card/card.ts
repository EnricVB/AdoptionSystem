import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { Wave } from "../wave/wave";
import { SimplifiedPet } from '@app/models';

@Component({
  selector: 'app-card',
  imports: [CommonModule, Wave],
  templateUrl: './card.html',
})
export class Card {
  @Input() pet!: SimplifiedPet;

  constructor(private router: Router) {
  } 

  getAgeFromBirthdate(birthdate: string): string {
    if (!birthdate) return 'Edad no especificada';
    
    const birth = new Date(birthdate);
    const today = new Date();
    const ageInMonths = (today.getFullYear() - birth.getFullYear()) * 12 + today.getMonth() - birth.getMonth();
    
    if (ageInMonths < 12) {
      return `${ageInMonths} ${ageInMonths === 1 ? 'mes' : 'meses'}`;
    } else {
      const years = Math.floor(ageInMonths / 12);
      return `${years} ${years === 1 ? 'año' : 'años'}`;
    }
  }

  onCardClick(): void {
    this.router.navigate(['/adopt/', this.pet.id]);
  }
}
