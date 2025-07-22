import { CommonModule } from '@angular/common';
import { Component, ViewChild, ElementRef, OnInit, ChangeDetectorRef } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { PopUp } from '@app/components/pop-up/pop-up';
import { Species, PetGender, Base64Image, Pet, PetStatus } from '@app/models';
import { ApiService } from '@app/services/api.service';
import { Carousel } from "@app/components/carousel/carousel";

@Component({
  selector: 'app-new',
  standalone: true,
  imports: [CommonModule, PopUp, Carousel],
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
  @ViewChild('statusSelect') statusSelect!: ElementRef<HTMLSelectElement>;
  @ViewChild('vaccinationSelect') vaccinationSelect!: ElementRef<HTMLSelectElement>;
  @ViewChild('descriptionTextarea') descriptionTextarea!: ElementRef<HTMLTextAreaElement>;

  @ViewChild('successPopUp') successPopUp!: PopUp;
  @ViewChild('errorPopUp') errorPopUp!: PopUp;

  @ViewChild('carousel') carousel!: Carousel;

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================

  species: ReadonlyArray<Species> = [];
  vaccinationHistory: { vaccine: string; date: string }[] = [];
  showVaccinationHistory: boolean = false;
  
  // Edit mode properties
  isEditMode: boolean = false;
  petId: number | null = null;
  currentPet: Pet | null = null;
  
  // Available pet statuses
  petStatuses = [
    { value: PetStatus.Available, label: 'Disponible' },
    { value: PetStatus.Adopted, label: 'Adoptado' },
    { value: PetStatus.FosterHome, label: 'Casa de Acogida' }
  ];

  // ======================================
  // CONSTRUCTOR
  // ======================================
  constructor(
    private apiService: ApiService,
    private router: Router,
    private route: ActivatedRoute,
    private cdr: ChangeDetectorRef,
  ) {
  }

  // ======================================
  // LIFECYCLE
  // ======================================
  ngOnInit(): void {
    this.checkRouteParams();
    this.fetchSpecies();
  }

  // ======================================
  // ROUTE HANDLING
  // ======================================
  private checkRouteParams(): void {
    // Check if we have a pet ID in the route (edit mode)
    const petIdParam = this.route.snapshot.paramMap.get('id');

    if (petIdParam) {
      this.petId = parseInt(petIdParam, 10);
      this.isEditMode = true;
      this.loadPetData();
    }
    
    // Check for query parameters (pre-filled data)
    this.route.queryParams.subscribe(params => {
      if (Object.keys(params).length > 0) {
        this.prefillFormData(params);
      }
    });
  }

  private loadPetData(): void {
    if (this.petId) {
      this.apiService.getPetById(this.petId).subscribe({
        next: (data) => this.handlePetDataSuccess(data),
        error: (err) => this.handlePetDataError(err)
      });
    }
  }

  private handlePetDataSuccess(data: { content: Pet }): void {
    this.currentPet = data.content;
    this.populateFormWithPetData();
  }

  private handlePetDataError(error: any): void {
    console.error('Error fetching pet data: ', error);
    this.errorPopUp.start('Error al cargar los datos de la mascota');

    setTimeout(() => {
      this.goBack();
    }, 2000);
  }

  private populateFormWithPetData(): void {
    if (!this.currentPet) return;

    // Wait for ViewChild elements to be available
    setTimeout(() => {
      if (this.petNameInput) this.petNameInput.nativeElement.value = this.currentPet!.name || '';
      if (this.speciesSelect) this.speciesSelect.nativeElement.value = this.currentPet!.species_id?.toString() || '';
      if (this.breedInput) this.breedInput.nativeElement.value = this.currentPet!.breed || '';
      if (this.genderSelect) this.genderSelect.nativeElement.value = this.currentPet!.gender || '';
      if (this.weightInput) this.weightInput.nativeElement.value = this.currentPet!.weight?.toString() || '';
      if (this.birthdateInput) this.birthdateInput.nativeElement.value = this.formatDateForInput(this.currentPet!.birthdate) || '';
      if (this.statusSelect) this.statusSelect.nativeElement.value = this.currentPet!.status || 'Available';
      if (this.descriptionTextarea) this.descriptionTextarea.nativeElement.value = this.currentPet!.description || '';
      if (this.vaccinationSelect) this.vaccinationSelect.nativeElement.value = this.currentPet!.is_vaccinated ? '1' : '2';

      // Load vaccination history
      if (this.currentPet!.vaccination_history) {
        this.vaccinationHistory = this.currentPet!.vaccination_history.map(vh => ({
          vaccine: vh.vaccine_name,
          date: this.formatDateForInput(vh.vaccination_date)
        }));
        this.showVaccinationHistory = this.vaccinationSelect.nativeElement.value === '1';
      }

      // Load images
      if (this.currentPet!.image_url) {
        this.apiService.getPetImageUrlsByFolder(this.currentPet!.image_url).subscribe({
          next: (urls) => {
            if (urls.length > 0) {
              this.carousel.images = urls;
            } else {
              this.carousel.images = [this.apiService.getPetImageUrl(this.currentPet!.image_url)];
            }
          },
          error: (err) => {
            this.errorPopUp.start('Error al cargar las imágenes de la mascota');
            console.error('Error fetching pet images: ', err);
          }
        });
      }

      this.cdr.detectChanges();
    }, 100);
  }

  private prefillFormData(params: any): void {
    // Wait for ViewChild elements to be available
    setTimeout(() => {
      if (params.name && this.petNameInput) this.petNameInput.nativeElement.value = params.name;
      if (params.species_id && this.speciesSelect) this.speciesSelect.nativeElement.value = params.species_id;
      if (params.breed && this.breedInput) this.breedInput.nativeElement.value = params.breed;
      if (params.gender && this.genderSelect) this.genderSelect.nativeElement.value = params.gender;
      if (params.weight && this.weightInput) this.weightInput.nativeElement.value = params.weight;
      if (params.birthdate && this.birthdateInput) this.birthdateInput.nativeElement.value = this.formatDateForInput(params.birthdate);
      if (params.description && this.descriptionTextarea) this.descriptionTextarea.nativeElement.value = params.description;
      if (params.is_vaccinated && this.vaccinationSelect) this.vaccinationSelect.nativeElement.value = params.is_vaccinated;

      this.cdr.detectChanges();
    }, 100);
  }

  private formatDateForInput(dateString: string): string {
    if (!dateString) return '';
    try {
      const date = new Date(dateString);
      return date.toISOString().split('T')[0];
    } catch {
      return '';
    }
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

  // =======================================
  // ACTION BUTTONS
  // =======================================

  onVaccinationChange(event: any): void {
    this.showVaccinationHistory = event.target.value === '1';
  }

  async savePet(): Promise<void> {
    if (this.isEditMode && this.petId) {
      await this.updatePet();
    } else {
      await this.createPet();
    }
  }

  async createPet(): Promise<void> {
    // Upload images if selected
    let base64Image: Base64Image | null = null;

    for (const file of this.carousel.selectedFiles) {
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
    const statusValue = this.statusSelect.nativeElement.value || 'Available';
    
    // Validate status value and provide fallback
    let validStatus: PetStatus;
    switch (statusValue) {
      case PetStatus.Available:
      case PetStatus.Adopted:
      case PetStatus.FosterHome:
        validStatus = statusValue as PetStatus;
        break;
      default:
        console.warn(`Invalid status value: ${statusValue}, defaulting to Available`);
        validStatus = PetStatus.Available;
    }

    const petData = {
      name: this.petNameInput.nativeElement.value.trim(),
      species_id: parseInt(this.speciesSelect.nativeElement.value),
      breed: this.breedInput.nativeElement.value.trim() || 'Desconocido',
      gender: this.genderSelect.nativeElement.value as PetGender,
      weight: parseFloat(this.weightInput.nativeElement.value) || 0,
      birthdate: this.toRFC3339(this.birthdateInput.nativeElement.value),
      description: this.descriptionTextarea.nativeElement.value.trim() || '',
      status: validStatus,
      image_url: base64Image ? base64Image.name.substring(0, base64Image.name.indexOf('/')) : '',
      is_vaccinated: this.vaccinationSelect.nativeElement.value === '1',
      vaccination_history: this.vaccinationHistory.map(vaccination => ({
        vaccine_name: vaccination.vaccine,
        vaccination_date: this.toRFC3339(vaccination.date)
      })),
    };

    // Validate required fields
    if (!petData.name) {
      this.errorPopUp.start('El nombre de la mascota es obligatorio.');
      return;
    } else if (petData.species_id <= 0) {
      this.errorPopUp.start('Debe seleccionar una especie válida.');
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

  async updatePet(): Promise<void> {
    if (!this.petId) return;

    // Upload new images if selected
    let base64Image: Base64Image | null = null;
    let imageFolder = this.currentPet?.image_url || '';

    // Process image uploads first if there are new files
    if (this.carousel.selectedFiles.length > 0) {
      for (const file of this.carousel.selectedFiles) {
        try {
          base64Image = await this.toBase64(file, this.currentPet?.image_url);

          // Upload the image to the server and wait for completion
          await new Promise<void>((resolve, reject) => {
            this.apiService.uploadPetImage(base64Image!).subscribe({
              next: (response) => {
                this.handleImageUploadSuccess(response);
                // Extract folder name from the uploaded image
                if (base64Image && base64Image.name.includes('/')) {
                  imageFolder = base64Image.name.substring(0, base64Image.name.indexOf('/'));
                }
                resolve();
              },
              error: (error) => {
                this.handleImageUploadError(error);
                reject(error);
              }
            });
          });
        } catch (error) {
          this.handleImageUploadError(error);
          return; // Stop execution if image upload fails
        }
      }
    }

    // Collect form data
    const statusValue = this.statusSelect.nativeElement.value || 'Available';
    
    // Validate status value and provide fallback
    let validStatus: PetStatus;
    switch (statusValue) {
      case PetStatus.Available:
      case PetStatus.Adopted:
      case PetStatus.FosterHome:
        validStatus = statusValue as PetStatus;
        break;
      default:
        console.warn(`Invalid status value: ${statusValue}, defaulting to Available`);
        validStatus = PetStatus.Available;
    }

    const petData = {
      id: this.petId,
      name: this.petNameInput.nativeElement.value.trim(),
      species_id: parseInt(this.speciesSelect.nativeElement.value),
      breed: this.breedInput.nativeElement.value.trim() || 'Desconocido',
      gender: this.genderSelect.nativeElement.value as PetGender,
      weight: parseFloat(this.weightInput.nativeElement.value) || 0,
      birthdate: this.toRFC3339(this.birthdateInput.nativeElement.value),
      description: this.descriptionTextarea.nativeElement.value.trim() || '',
      status: validStatus,
      image_url: imageFolder, // Include image folder even in update
      is_vaccinated: this.vaccinationSelect.nativeElement.value === '1',
      vaccination_history: this.vaccinationHistory.map(vaccination => ({
        vaccine_name: vaccination.vaccine,
        vaccination_date: this.toRFC3339(vaccination.date)
      }))
    };

    console.log('Updating pet with data:', petData);

    // Validate required fields
    if (!petData.name) {
      this.errorPopUp.start('El nombre de la mascota es obligatorio.');
      return;
    } else if (petData.species_id <= 0) {
      this.errorPopUp.start('Debe seleccionar una especie válida.');
      return;
    } else if (isNaN(petData.weight) || petData.weight <= 0) {
      this.errorPopUp.start('El peso de la mascota debe ser un número positivo.');
      return;
    } else if (!petData.birthdate) {
      this.errorPopUp.start('La fecha de nacimiento es obligatoria.');
      return;
    }

    // Update pet via API
    this.apiService.updatePet(this.petId, petData).subscribe({
      next: (response) => this.handleUpdatePetSuccess(response),
      error: (error) => this.handleUpdatePetError(error)
    });
  }

  // ======================================
  // SUCCESS & ERROR HANDLERS
  // ======================================
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
    
    // Reset carousel index
    if (this.carousel) {
      this.carousel.clearImages();
    }
  }

  // ======================================
  // HANDLERS FOR API RESPONSES
  // ======================================
  private handleCreatePetSuccess(response: any): void {
    this.successPopUp.start("Mascota creada correctamente.");
    setTimeout(() => {
      this.goBack();
    }, 2000);
  }

  private handleCreatePetError(error: any): void {
    this.errorPopUp.start("Ha ocurrido un problema al crear la mascota.");
    console.log('Error creating pet: ', error);
  }

  private handleUpdatePetSuccess(response: any): void {
    this.successPopUp.start("Mascota actualizada correctamente.");
    setTimeout(() => {
      this.goBack();
    }, 2000);
  }

  private handleUpdatePetError(error: any): void {
    this.errorPopUp.start("Ha ocurrido un problema al actualizar la mascota.");
    console.log('Error updating pet: ', error);
  }

  private handleImageUploadSuccess(response: any): void {
  }

  private handleImageUploadError(error: any): void {
    console.log('Error uploading image: ', error);
  }
  
  goBack(): void {
    if (this.isEditMode) {
      this.router.navigate(['/adopt/', this.petId]);
    } else {
      this.router.navigate(['/adopt']);
    }
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

  private toBase64(file: File, foldersName?: string): Promise<Base64Image> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      
      reader.onload = (e) => {
        try {
          const result = reader.result as string;
          const base64String = result.split(',')[1]; // Remove data:image/...;base64, prefix
          const folderName = foldersName ? `${foldersName}/` : `$${Math.random().toString(36).substring(2, 8)}`;
          
          const base64Image: Base64Image = {
            base64: base64String,
            name: `${folderName}/${Date.now()}.png`
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
