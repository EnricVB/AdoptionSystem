import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { Species } from '@app/models';
import { ApiService } from '@app/services/api.service';

@Component({
  selector: 'app-new',
  imports: [CommonModule],
  templateUrl: './new.html',
})
export class NewPet {

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================
  species: ReadonlyArray<Species> = [];
  vaccinationHistory: { vaccine: string; date: string }[] = [];

  // ======================================
  // CONSTRUCTOR
  // ======================================
  constructor(
    private apiService: ApiService,
    private router: Router
  ) {
  }

  // ======================================
  // LIFECYCLE
  // ======================================
  ngOnInit(): void {
    this.fetchSpecies();
  }

    // ======================================
  // DATA FETCHING & HANDLING
  // ======================================
  fetchSpecies(): void { 
    this.apiService.getSpecies().subscribe({
      next: (data) => this.handleSpeciesSuccess(data),
      error: (err) => this.handleSpeciesError(err)
    });
  }
  
  private handleSpeciesSuccess(data: { content: Species[] }): void {
    this.species = data.content;
  }

  private handleSpeciesError(error: any): void {
    console.error('Error fetching pets: ', error);
  }

  // ======================================
  // EVENT HANDLERS
  // ======================================
  onImageError(event: any): void {
    // Replace broken image with a placeholder
    event.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjQwMCIgdmlld0JveD0iMCAwIDQwMCA0MDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSI0MDAiIGhlaWdodD0iNDAwIiBmaWxsPSIjRjNGNEY2Ii8+CjxwYXRoIGQ9Ik0yMDAgMTAwQzE0NC43NzIgMTAwIDEwMCAxNDQuNzcyIDEwMCAyMDBTMTQ0Ljc3MiAzMDAgMjAwIDMwMFMyNTAgMjU1LjIyOCAyNTAgMjAwUzIwNS4yMjggMTAwIDIwMCAxMDBaTTIwMCAyNTBDMTcyLjM4NiAyNTAgMTUwIDIyNy42MTQgMTUwIDIwMEMxNTAgMTcyLjM4NiAxNzIuMzg2IDE1MCAyMDAgMTUwQzIyNy42MTQgMTUwIDI1MCAxNzIuMzg2IDI1MCAyMDBDMjUwIDIyNy42MTQgMjI3LjYxNCAyNTAgMjAwIDI1MFoiIGZpbGw9IiM5Q0EzQUYiLz4KPC9zdmc+';
  }

  onFileSelected(event: any): void { 
    
  }

  addVaccination(vaccineNameInput: HTMLInputElement, vaccinationDateInput: HTMLInputElement): void {
    const vaccineName = vaccineNameInput.value.trim();
    const vaccinationDate = vaccinationDateInput.value;

    if (vaccineName && vaccinationDate) {
      this.vaccinationHistory.push({ 
        vaccine: vaccineName, 
        date: vaccinationDate 
      });
      
      vaccineNameInput.value = '';
      vaccinationDateInput.value = '';
    } else {
      alert('Por favor, completa todos los campos');
    }
  }

  removeVaccination(vaccination: { vaccine: string; date: string }): void {
    this.vaccinationHistory = this.vaccinationHistory.filter(v => v !== vaccination);
  }

  // =======================================
  // ACTION BUTTONS
  // =======================================

  
  goBack(): void {
    this.router.navigate(['/adopt']);
  }
}
