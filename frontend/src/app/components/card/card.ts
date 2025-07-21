import { Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { SimplifiedPet } from '@app/models';
import { ApiService } from '@app/services/api.service';

@Component({
  selector: 'app-card',
  imports: [CommonModule],
  templateUrl: './card.html',
})
export class Card implements OnInit {
  @Input() pet!: SimplifiedPet;
  
  petImageUrl: string = '';

  constructor(
    private router: Router,
    public apiService: ApiService
  ) {
  } 

  ngOnInit(): void {
    this.loadPetImage();
  }

  loadPetImage(): void {
    if (!this.pet?.image_url) {
      return;
    }

    this.apiService.getPetImageUrlsByFolder(this.pet.image_url).subscribe({
      next: (urls) => {
        if (urls.length > 0) {
          this.petImageUrl = urls[0];
        }
      },
      error: (error) => {
        console.error(`Error loading images for folder ${this.pet?.image_url}:`, error);
        if (this.pet?.image_url) {
          this.petImageUrl = this.apiService.getPetImageUrl(this.pet.image_url);
        }
      }
    });
  } 
  
  getAgeFromBirthdate(birthdate: string): string {
    const [day, month, year] = birthdate.split('/');
    const birth = new Date(parseInt(year), parseInt(month) - 1, parseInt(day));
    const today = new Date();
    const ageInMonths = (today.getFullYear() - birth.getFullYear()) * 12 + today.getMonth() - birth.getMonth();

    if (ageInMonths < 12) {
      return `${ageInMonths} ${ageInMonths === 1 ? 'mes' : 'meses'}`;
    } else if(ageInMonths >= 12) {
      const years = Math.floor(ageInMonths / 12);
      return `${years} ${years === 1 ? 'año' : 'años'}`;
    }

    return 'Edad no especificada';
  }

  getGenderName() {
    switch (this.pet.gender) {
      case 'Male':
        return 'Macho';
      case 'Female':
        return 'Hembra';
      default:
        return 'Desconocido';
    }
  }

  onCardClick(): void {
    this.router.navigate(['/adopt/', this.pet.id]);
  }
}
