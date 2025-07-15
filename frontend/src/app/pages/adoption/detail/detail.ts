import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ApiService } from '@app/services/api.service';
import { CommonModule } from '@angular/common';
import { Pet, PetStatus } from '@app/models';

@Component({
  selector: 'app-detail',
  imports: [CommonModule],
  templateUrl: './detail.html',
})
export class Detail implements OnInit {

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================
  
  petId!: number;
  pet: Pet | null = null;
  error: string | null = null;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private apiService: ApiService
  ) {
    // Initialization logic can go here if needed
  }

  ngOnInit(): void {
    this.petId = Number(this.route.snapshot.paramMap.get('id'));
    
    if (this.petId) {
      this.fetchPetDetails();
    } else {
      this.error = 'ID de mascota no válido';
    }
  }

  private fetchPetDetails(): void {
    this.error = null;
    
    this.apiService.getPetById(this.petId).subscribe({
      next: (data) => this.handlePetDetailsSuccess(data.content),
      error: (err) => this.handlePetDetailsError(err)
    });
  }

  private handlePetDetailsSuccess(data: Pet): void {
    this.pet = data;
    console.log('Fetched pet details:', this.pet);
    console.log(data.image_url);
  }

  private handlePetDetailsError(err: any): void {
    this.error = 'Error al cargar los detalles de la mascota';
  }
}
