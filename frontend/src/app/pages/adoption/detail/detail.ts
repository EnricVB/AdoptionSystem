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
  }

  private handlePetDetailsError(err: any): void {
    this.error = 'Error al cargar los detalles de la mascota';
  }

  getStatusClass(): string {
    if (!this.pet || !this.pet.status) return 'text-gray-600 bg-gray-100';
    
    switch (this.pet.status) {
      case PetStatus.Available:
        return 'text-green-600 bg-green-100';
      case PetStatus.Adopted:
        return 'text-red-600 bg-red-100';
      case PetStatus.FosterHome:
        return 'text-yellow-600 bg-yellow-100';
      default:
        return 'text-gray-600 bg-gray-100';
    }
  }

  getStatusText(): string {
    if (!this.pet || !this.pet.status) return 'Estado desconocido';
    
    switch (this.pet.status) {
      case PetStatus.Available:
        return 'Disponible';
      case PetStatus.Adopted:
        return 'Adoptado';
      case PetStatus.FosterHome:
        return 'En acogida';
      default:
        return 'Estado desconocido';
    }
  }

  onImageError(event: any): void {
    // Replace broken image with a placeholder
    event.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjQwMCIgdmlld0JveD0iMCAwIDQwMCA0MDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSI0MDAiIGhlaWdodD0iNDAwIiBmaWxsPSIjRjNGNEY2Ii8+CjxwYXRoIGQ9Ik0yMDAgMTAwQzE0NC43NzIgMTAwIDEwMCAxNDQuNzcyIDEwMCAyMDBTMTQ0Ljc3MiAzMDAgMjAwIDMwMFMyNTAgMjU1LjIyOCAyNTAgMjAwUzIwNS4yMjggMTAwIDIwMCAxMDBaTTIwMCAyNTBDMTcyLjM4NiAyNTAgMTUwIDIyNy42MTQgMTUwIDIwMEMxNTAgMTcyLjM4NiAxNzIuMzg2IDE1MCAyMDAgMTUwQzIyNy42MTQgMTUwIDI1MCAxNzIuMzg2IDI1MCAyMDBDMjUwIDIyNy42MTQgMjI3LjYxNCAyNTAgMjAwIDI1MFoiIGZpbGw9IiM5Q0EzQUYiLz4KPC9zdmc+';
  }

  goBack(): void {
    this.router.navigate(['/adopt']);
  }
}
