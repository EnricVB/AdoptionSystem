import { CommonModule } from '@angular/common';
import { Component, ViewChild, ElementRef, OnInit, ChangeDetectorRef } from '@angular/core';
import { Router } from '@angular/router';
import { PopUp } from '@app/components/pop-up/pop-up';
import { Species, PetGender, Base64Image } from '@app/models';
import { ApiService } from '@app/services/api.service';
import {
  CarouselComponent,
  CarouselControlComponent,
  CarouselInnerComponent,
  CarouselItemComponent
} from '@coreui/angular';

@Component({
  selector: 'app-new',
  standalone: true,
  imports: [CommonModule, PopUp, CarouselComponent, CarouselControlComponent, CarouselInnerComponent, CarouselItemComponent,],
  templateUrl: './new.html',
})
export class NewPet implements OnInit {

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

  @ViewChild('carousel') carousel!: CarouselComponent;
  @ViewChild('modal') modal!: ElementRef<HTMLElement>;
  @ViewChild('modalImage') modalImage!: ElementRef<HTMLImageElement>;
  @ViewChild('file') fileInput!: ElementRef<HTMLInputElement>;

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================

  species: ReadonlyArray<Species> = [];
  vaccinationHistory: { vaccine: string; date: string }[] = [];
  selectedFiles: File[] = [];
  carouselImages: string[] = [];
  showVaccinationHistory: boolean = false;

  // ======================================
  // CONSTRUCTOR
  // ======================================
  constructor(
    private apiService: ApiService,
    private router: Router,
    private cdr: ChangeDetectorRef,
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
    console.warn('Error loading image:', event.target.src);

    
    if(event.target.src.toString().includes('error')) {
      // Replace broken image with a placeholder
      event.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjQwMCIgdmlld0JveD0iMCAwIDQwMCA0MDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSI0MDAiIGhlaWdodD0iNDAwIiBmaWxsPSIjRjNGNEY2Ii8+CjxwYXRoIGQ9Ik0yMDAgMTAwQzE0NC43NzIgMTAwIDEwMCAxNDQuNzcyIDEwMCAyMDBTMTQ0Ljc3MiAzMDAgMjAwIDMwMFMyNTAgMjU1LjIyOCAyNTAgMjAwUzIwNS4yMjggMTAwIDIwMCAxMDBaTTIwMCAyNTBDMTcyLjM4NiAyNTAgMTUwIDIyNy42MTQgMTUwIDIwMEMxNTAgMTcyLjM4NiAxNzIuMzg2IDE1MCAyMDAgMTUwQzIyNy42MTQgMTUwIDI1MCAxNzIuMzg2IDI1MCAyMDBDMjUwIDIyNy42MTQgMjI3LjYxNCAyNTAgMjAwIDI1MFoiIGZpbGw9IiM5Q0EzQUYiLz4KPC9zdmc+';
    } else {
      this.errorPopUp.start("Error al cargar la imagen, intente nuevamente.")
    }
  }

  onImageClick(index: number): void {
    this.modalImage.nativeElement.src = this.carouselImages[index];
    this.modal.nativeElement.hidden = false;
  }

  closeImage(): void {
    this.modal.nativeElement.hidden = true;
  }

  onFileSelected(event: any): void { 
    const file = event.target.files[0];
    const reader = new FileReader();

    if (!file) {
      return;
    }

    // Add the file to selected files
    this.selectedFiles.push(file);

    reader.onload = (e) => {
      const imageUrl = e.target?.result as string;
      
      // Add the image URL to carousel images
      this.carouselImages.push(imageUrl);
      
      // Move carousel to the last image (newly added)
      const lastIndex = this.carouselImages.length - 1;
      
      // Use setTimeout to ensure the DOM is updated before changing the index
      setTimeout(() => {
        if (this.carousel && lastIndex >= 0) {
          this.carousel.activeIndex.set(lastIndex);
        }
        // Force change detection to update the view
        this.cdr.detectChanges();
      }, 50);
    };

    reader.onerror = () => {
      // Remove the file from selected files if reading fails
      this.selectedFiles.pop();
      this.errorPopUp.start('Error al cargar la imagen');
    };

    // Read the file as a data URL
    reader.readAsDataURL(file);
    
    // Clear the input to allow selecting the same file again
    event.target.value = '';
  }

  addVaccination(vaccineNameInput: HTMLInputElement, vaccinationDateInput: HTMLInputElement): void {
    const vaccineName = vaccineNameInput.value.trim();
    const vaccinationDate = vaccinationDateInput.value;

    if (vaccineName && vaccinationDate) {
      // Validate date format (YYYY-MM-DD)
      if (isNaN(Date.parse(vaccinationDate))) {
        this.errorPopUp.start('La fecha de vacunación no es válida');
        return;
      }

      this.vaccinationHistory.push({ 
        vaccine: vaccineName, 
        date: vaccinationDate 
      });
      
      vaccineNameInput.value = '';
      vaccinationDateInput.value = '';
    } else {
      this.errorPopUp.start('Por favor, completa todos los campos de vacunación');
    }
  }

  removeVaccination(vaccination: { vaccine: string; date: string }): void {
    this.vaccinationHistory = this.vaccinationHistory.filter(v => v !== vaccination);
  }

  onVaccinationChange(event: any): void {
    this.showVaccinationHistory = event.target.value === '1';
  }

  // =======================================
  // ACTION BUTTONS
  // =======================================

  async createAdoption(): Promise<void> {
    // Upload image if selected
    let base64Image: Base64Image | null = null;

    for (const file of this.selectedFiles) {
      try {
        base64Image = await this.toBase64(file);

        // Upload the image to the server
        this.apiService.uploadPetImage(base64Image).subscribe({
          next: (response) => this.handleImageUploadSuccess(response),
          error: (error) => this.handleImageUploadError(error)
        });
      } catch (error) {
        this.handleImageUploadError(error);
      }
    }

    // Collect form data
    const petData = {
      name: this.petNameInput.nativeElement.value.trim(),
      species_id: parseInt(this.speciesSelect.nativeElement.value),
      breed: this.breedInput.nativeElement.value.trim() || 'Desconocido',
      gender: this.genderSelect.nativeElement.value as PetGender,
      weight: parseFloat(this.weightInput.nativeElement.value) || 0,
      birthdate: this.toRFC3339(this.birthdateInput.nativeElement.value),
      description: this.descriptionTextarea.nativeElement.value.trim() || '',
      status: 'Available',
      image_url: base64Image ? base64Image.name.substring(0, base64Image.name.indexOf('/')) : '',
      is_vaccinated: this.vaccinationSelect.nativeElement.value === '1',
      vaccination_history: this.vaccinationHistory.map(vaccination => ({
      vaccine_name: vaccination.vaccine,
      vaccination_date: this.toRFC3339(vaccination.date)
      })),
    };

    // Basic validation
    if (!petData.name) {
      this.errorPopUp.start('El nombre de la mascota es obligatorio.');
      return;
    } else if (!petData.species_id || isNaN(petData.species_id)) {
      this.errorPopUp.start('La especie de la mascota es obligatoria.');
      return;
    } else if (!petData.gender) {
      this.errorPopUp.start('El género de la mascota es obligatorio.');
      return;
    } else if (isNaN(petData.weight) || petData.weight <= 0) {
      this.errorPopUp.start('El peso de la mascota debe ser un número positivo.');
      return;
    } else if (!petData.birthdate) {
      this.errorPopUp.start('La fecha de nacimiento es obligatoria.');
      return;
    }

    // Create pet via API
    this.apiService.createPet(petData).subscribe({
      next: (response) => this.handleCreatePetSuccess(response),
      error: (error) => this.handleCreatePetError(error)
    });
  }

  // Clear all form fields
  clearForm(): void {
    this.petNameInput.nativeElement.value = '';
    this.speciesSelect.nativeElement.value = '';
    this.breedInput.nativeElement.value = '';
    this.genderSelect.nativeElement.value = '';
    this.weightInput.nativeElement.value = '';
    this.birthdateInput.nativeElement.value = '';
    this.vaccinationSelect.nativeElement.value = '';
    this.descriptionTextarea.nativeElement.value = '';
    this.vaccinationHistory = [];
    this.showVaccinationHistory = false;
    
    // Reset image arrays
    this.selectedFiles = [];
    this.carouselImages = [];
    
    // Reset carousel index
    if (this.carousel) {
      this.carousel.activeIndex.set(0);
    }
    
    // Reset file input
    if (this.fileInput) {
      this.fileInput.nativeElement.value = '';
    }
    
    // Force change detection
    this.cdr.detectChanges();
  }

  deleteImage(): void {
    if (this.carouselImages.length === 0) {
      return;
    }
    
    const activeIndex = this.carousel.activeIndex() || 0;
    
    // Remove from arrays
    this.selectedFiles.splice(activeIndex, 1);
    this.carouselImages.splice(activeIndex, 1);

    // Adjust carousel index if necessary
    if (this.carouselImages.length > 0) {
      const newIndex = activeIndex >= this.carouselImages.length ? 0 : activeIndex;
      this.carousel.activeIndex.set(newIndex);
    } else {
      this.carousel.activeIndex.set(0);
    }
    
    // Force change detection
    this.cdr.detectChanges();
  }

  // ======================================
  // HANDLERS FOR API RESPONSES
  // ======================================
  private handleCreatePetSuccess(response: any): void {
    this.successPopUp.start("Mascota creada correctamente.");
    setTimeout(() => {
      this.router.navigate(['/adopt']);
    }, 2000);
  }

  private handleCreatePetError(error: any): void {
    this.errorPopUp.start("Ha ocurrido un problema al crear la mascota.");
    console.log('Error creating pet: ', error);
  }

  private handleImageUploadSuccess(response: any): void {
    console.log('Image uploaded successfully:', response);
  }

  private handleImageUploadError(error: any): void {
    console.log('Error uploading image: ', error);
  }
  
  goBack(): void {
    this.router.navigate(['/adopt']);
  }

  // ======================================
  // UTILITY METHODS
  // ======================================
 
  private toRFC3339(dateStr: string): string {
    // El input date ya viene en formato YYYY-MM-DD
    // Solo necesitamos validar que sea una fecha válida
    if (!dateStr) {
      throw new Error('Date string is empty');
    }

    const date = new Date(dateStr);
    if (isNaN(date.getTime())) {
      throw new Error('Invalid date format');
    }

    // Devolver la fecha en formato RFC3339 completo que Go puede parsear
    return `${dateStr}T00:00:00Z`; // "2025-12-31T00:00:00Z"
  }

  private toBase64(file: File): Promise<Base64Image> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      
      reader.onload = (e) => {
        try {
          const result = reader.result as string;
          const base64String = result.split(',')[1]; // Remove data:image/...;base64, prefix
          
          const base64Image: Base64Image = {
            base64: base64String,
            name: `$${Math.random().toString(36).substring(2, 8)}/${Date.now()}.png`
          };
          
          resolve(base64Image);
        } catch (error) {
          reject(error);
        }
      };
      
      reader.onerror = () => {
        reject(new Error('Error reading file'));
      };
      
      reader.readAsDataURL(file);
    });
  }
}
