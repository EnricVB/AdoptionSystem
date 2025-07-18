import { Component, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ApiService } from '@app/services/api.service';
import { CommonModule } from '@angular/common';
import { Pet, PetStatus, User } from '@app/models';
import { CookieService } from '@app/services/cookie.service';
import { AuthService } from '@app/services/auth.service';
import { PopUp } from '@app/components/pop-up/pop-up';

@Component({
  selector: 'app-detail',
  imports: [CommonModule, PopUp],
  templateUrl: './detail.html',
})
export class Detail implements OnInit {

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================

  @ViewChild('successPopUp') successPopUp!: PopUp;
  @ViewChild('errorPopUp') errorPopUp!: PopUp;
  
  petId!: number;
  pet: Pet | null = null;
  error: string | null = null;
  petImageUrl: string = '';

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private apiService: ApiService,
    public authService: AuthService,
    private cookieService: CookieService
  ) {
    // Initialization logic can go here if needed
  }

  async ngOnInit(): Promise<void> {
    this.petId = Number(this.route.snapshot.paramMap.get('id'));

    if (this.petId) {
      await this.fetchPetDetails();
    } else {
      this.error = 'ID de mascota no válido';
    }

    this.loadPetImage();
  }

  private fetchPetDetails(): Promise<void> {
    this.error = null;
    
    return new Promise((resolve, reject) => {
      this.apiService.getPetById(this.petId).subscribe({
        next: (data) => {
          this.handlePetDetailsSuccess(data.content);
          resolve();
        },
        error: (err) => {
          this.handlePetDetailsError(err);
          reject(err);
        }
      });
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
    event.target.style.backgroundImage = 'url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjQwMCIgdmlld0JveD0iMCAwIDQwMCA0MDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSI0MDAiIGhlaWdodD0iNDAwIiBmaWxsPSIjRjNGNEY2Ii8+CjxwYXRoIGQ9Ik0yMDAgMTAwQzE0NC43NzIgMTAwIDEwMCAxNDQuNzcyIDEwMCAyMDBTMTQ0Ljc3MiAzMDAgMjAwIDMwMFMyNTAgMjU1LjIyOCAyNTAgMjAwUzIwNS4yMjggMTAwIDIwMCAxMDBaTTIwMCAyNTBDMTcyLjM4NiAyNTAgMTUwIDIyNy42MTQgMTUwIDIwMEMxNTAgMTcyLjM4NiAxNzIuMzg2IDE1MCAyMDAgMTUwQzIyNy42MTQgMTUwIDI1MCAxNzIuMzg2IDI1MCAyMDBDMjUwIDIyNy42MTQgMjI3LjYxNCAyNTAgMjAwIDI1MFoiIGZpbGw9IiM5Q0EzQUYiLz4KPC9zdmc+)';
  }

  goBack(): void {
    this.router.navigate(['/adopt']);
  }

  async fosterHomeContact(): Promise<void> {
    if (!this.pet) return;

    const reason = "Solicitud de adopción para " + this.pet.name;
    const message = "Hola, estoy interesado en acoger a " + this.pet.name + ", que está en su hogar de acogida. Por favor, contáctenme para más detalles.";

    this.apiService.sendPetFosterHomeContact(
      this.pet.id,
      this.authService.getLoggedInUser?.id || 0,
      this.pet.adopt_user_id || 0,
      reason,
      message
    ).subscribe({
      next: (response) => this.nextSendPetFosterHomeContact(response),
      error: (err) => this.submitError(err)
    });
  }

  async submitAdoption(): Promise<void> {
    if (!this.pet) return;
    
    this.apiService.submitPetAdoptionRequest(
      this.pet.id,
      this.authService.getLoggedInUser?.id || 0,
    ).subscribe({
      next: (response) => this.nextSubmitPetAdoptionRequest(response),
      error: (err) => this.submitError(err)
    });
  }

  async submitFosterHome(): Promise<void> {
    if (!this.pet) return;

    this.apiService.submitPetFosterHomeRequest(
      this.pet.id,
      this.authService.getLoggedInUser?.id || 0,
    ).subscribe({
      next: (response) => this.nextSubmitPetFosterHomeRequest(response),
      error: (err) => this.submitError(err)
    });
  }

  // ======================================
  // SUBSCRIBE METHODS
  // ======================================
  private nextSendPetFosterHomeContact(response: any): void {
    this.successPopUp.start('Solicitud de contacto enviada correctamente.');
  }

  private nextSubmitPetAdoptionRequest(response: any): void {
    this.successPopUp.start('Solicitud de adopción enviada correctamente.');
  }

  private nextSubmitPetFosterHomeRequest(response: any): void {
    this.successPopUp.start('Solicitud de acogida enviada correctamente.');
  }

  private submitError(err: any): void {
    this.errorPopUp.start('Error al enviar la solicitud.');
  }

  // ======================================
  // UTILITY METHODS
  // ======================================
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
}
