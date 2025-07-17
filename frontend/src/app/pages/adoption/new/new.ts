import { CommonModule } from '@angular/common';
import { Component, ViewChild, ElementRef } from '@angular/core';
import { Router } from '@angular/router';
import { PopUp } from '@app/components/pop-up/pop-up';
import { Species, PetGender } from '@app/models';
import { ApiService } from '@app/services/api.service';

@Component({
  selector: 'app-new',
  imports: [CommonModule, PopUp],
  templateUrl: './new.html',
})
export class NewPet {

  // ======================================
  // VIEWCHILD REFERENCES
  // ======================================
  @ViewChild('petNameInput') petNameInput!: ElementRef<HTMLInputElement>;
  @ViewChild('speciesSelect') speciesSelect!: ElementRef<HTMLSelectElement>;
  @ViewChild('breedInput') breedInput!: ElementRef<HTMLInputElement>;
  @ViewChild('genderSelect') genderSelect!: ElementRef<HTMLSelectElement>;
  @ViewChild('weightInput') weightInput!: ElementRef<HTMLInputElement>;
  @ViewChild('birthdateInput') birthdateInput!: ElementRef<HTMLInputElement>;
  @ViewChild('vaccinationSelect') vaccinationSelect!: ElementRef<HTMLSelectElement>;
  @ViewChild('descriptionTextarea') descriptionTextarea!: ElementRef<HTMLTextAreaElement>;

  @ViewChild('successPopUp') successPopUp!: PopUp;
  @ViewChild('errorPopUp') errorPopUp!: PopUp;

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================

  species: ReadonlyArray<Species> = [];
  vaccinationHistory: { vaccine: string; date: string }[] = [];
  selectedFile: File | null = null;

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
    const file = event.target.files[0];
    if (file) {
      this.selectedFile = file;
      const reader = new FileReader();
      reader.onload = (e) => {
        const img = event.target.closest('.relative').querySelector('img');
        if (img) {
          img.src = e.target?.result as string;
        }
      };
      reader.readAsDataURL(file);
    }
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
      this.errorPopUp.start('Por favor, completa todos los campos');
    }
  }

  removeVaccination(vaccination: { vaccine: string; date: string }): void {
    this.vaccinationHistory = this.vaccinationHistory.filter(v => v !== vaccination);
  }

  // =======================================
  // ACTION BUTTONS
  // =======================================

  createAdoption(): void {
    // Collect form data
    const petData = {
      name: this.petNameInput.nativeElement.value.trim(),
      species: this.speciesSelect.nativeElement.value,
      breed: this.breedInput.nativeElement.value.trim() || 'Desconocido',
      gender: this.genderSelect.nativeElement.value as PetGender,
      weight: parseFloat(this.weightInput.nativeElement.value) || 0,
      birthdate: this.birthdateInput.nativeElement.value,
      is_vaccinated: this.vaccinationSelect.nativeElement.value === '1',
      description: this.descriptionTextarea.nativeElement.value.trim() || '',
      vaccination_history: this.vaccinationHistory
    };

    // Basic validation
    if (!petData.name || !petData.species) {
      this.errorPopUp.start('Por favor, completa al menos el nombre y la especie de la mascota');
      return;
    }

    // Create pet via API
    this.apiService.createPet(petData).subscribe({
      next: (response) => {
        this.successPopUp.start("Mascota creada correctamente.");
        setTimeout(() => {
          this.router.navigate(['/adopt']);
        }, 2000);
      },

      error: (error) => {
        this.errorPopUp.start("Ha ocurrido un problema al crear la mascota.");
      }
    });
  }

  clearForm(): void {
    // Clear all form fields
    this.petNameInput.nativeElement.value = '';
    this.speciesSelect.nativeElement.value = '';
    this.breedInput.nativeElement.value = '';
    this.genderSelect.nativeElement.value = '';
    this.weightInput.nativeElement.value = '';
    this.birthdateInput.nativeElement.value = '';
    this.vaccinationSelect.nativeElement.value = '';
    this.descriptionTextarea.nativeElement.value = '';
    this.vaccinationHistory = [];
    this.selectedFile = null;
    
    // Reset image
    const img = document.querySelector('.relative img') as HTMLImageElement;
    if (img) {
      img.src = 'error';
    }
  }

  deleteImage(): void {
    this.selectedFile = null;
    const img = document.querySelector('.relative img') as HTMLImageElement;
    if (img) {
      img.src = 'error';
    }
  }
  
  goBack(): void {
    this.router.navigate(['/adopt']);
  }
}
