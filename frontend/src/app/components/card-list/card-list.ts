import { Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { Card } from '../card/card';
import { Breadcrumb, BreadcrumbItem } from '../breadcrumb/breadcrumb';
import { ApiService } from '@app/services/api.service';
import { SimplifiedPet, Species } from '@app/models';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-card-list',
  imports: [CommonModule, Card, FormsModule, RouterModule, Breadcrumb],
  templateUrl: './card-list.html',
  host: {
    'style': 'view-transition-name: card-list'
  }
})
export class CardList implements OnInit {

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================
  animals: ReadonlyArray<SimplifiedPet> = [];
  species: ReadonlyArray<Species> = [];
  error: string | null = null;
  currentPage = 1;

  breadcrumbItems: BreadcrumbItem[] = [
    { 
      label: 'Dashboard', 
      icon: 'fa-solid fa-house', 
      routerLink: '/dashboard' 
    },
    { 
      label: 'Adopción', 
      icon: 'fa-solid fa-paw' 
    }
  ];

  @Input() showPagination = true;
  @Input() showFilter = true;
  @Input() itemsPerPage: number = 8; 

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
    this.fetchPets();
    this.fetchSpecies();
  }

  // ======================================
  // DATA FETCHING & HANDLING
  // ======================================
  fetchPets(): void {
    this.error = null;
    this.apiService.getPets().subscribe({
      next: (data) => this.handlePetsSuccess(data),
      error: (err) => this.handlePetsError(err)
    });
  }

  fetchSpecies(): void { 
    this.apiService.getSpecies().subscribe({
      next: (data) => this.handleSpeciesSuccess(data),
      error: (err) => this.handleSpeciesError(err)
    });
  }

  filters = {
    name: '',
    status: '',
    species_id: '',
    gender: '',
    vaccinated: ''
  };

  fetchFilteredPets(): void {
    this.error = null;
    const filters = {
      name: this.filters.name || '',
      status: this.filters.status || '',
      species_id: this.filters.species_id || '',
      gender: this.filters.gender || '',
      vaccinated: this.filters.vaccinated || ''
    };
    this.apiService.getFilteredPets(filters).subscribe({
      
      next: (data) => this.handleFilteredPetsSuccess(data),
      error: (err) => this.handleFilteredPetsError(err)
    });
    console.log('Filters applied:', filters);
  }

  private handlePetsSuccess(data: any): void {
    this.animals = data.content.map((pet: any) => ({
      ...pet,
      birthdate: pet.birthdate && pet.birthdate !== "0001-01-01T00:00:00Z"
        ? this.parseDateToDDMMYYYY(pet.birthdate)
        : ''
    }));
  }

  private parseDateToDDMMYYYY(dateStr: string): string {
    const date = new Date(dateStr);
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
  }

  private handlePetsError(error: any): void {
    this.error = 'Error fetching pets.';
  }

  private handleSpeciesSuccess(data: { content: Species[] }): void {
    this.species = data.content;
  }

  private handleSpeciesError(error: any): void {
    this.error = 'Error fetching pets.';
    console.error('Error fetching pets: ', error);
  }

  private handleFilteredPetsSuccess(data: any): void {
    if (!data.content || !Array.isArray(data.content)) {
    this.animals = [];
    console.warn('No pets found or content is null.');
    return;
  }

    this.animals = data.content.map((pet: any) => ({
      ...pet,
      birthdate: pet.birthdate && pet.birthdate !== "0001-01-01T00:00:00Z"
        ? this.parseDateToDDMMYYYY(pet.birthdate)
        : ''
    }));
  }

  private handleFilteredPetsError(error: any): void {
    this.error = 'Error fetching filtered pets.';
    console.error('Error fetching filtered pets:', error);
  }

  public onImageError(event: any): void {
    // Replace broken image with a placeholder
    event.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjQwMCIgdmlld0JveD0iMCAwIDQwMCA0MDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSI0MDAiIGhlaWdodD0iNDAwIiBmaWxsPSIjRjNGNEY2Ii8+CjxwYXRoIGQ9Ik0yMDAgMTAwQzE0NC43NzIgMTAwIDEwMCAxNDQuNzcyIDEwMCAyMDBTMTQ0Ljc3MiAzMDAgMjAwIDMwMFMyNTAgMjU1LjIyOCAyNTAgMjAwUzIwNS4yMjggMTAwIDIwMCAxMDBaTTIwMCAyNTBDMTcyLjM4NiAyNTAgMTUwIDIyNy42MTQgMTUwIDIwMEMxNTAgMTcyLjM4NiAxNzIuMzg2IDE1MCAyMDAgMTUwQzIyNy42MTQgMTUwIDI1MCAxNzIuMzg2IDI1MCAyMDBDMjUwIDIyNy42MTQgMjI3LjYxNCAyNTAgMjAwIDI1MFoiIGZpbGw9IiM5Q0EzQUYiLz4KPC9zdmc+';
  }

  // ======================================
  // PAGINATION
  // ======================================
  get totalPages(): number {
    return Math.ceil(this.animals.length / this.itemsPerPage);
  }

  get paginatedAnimals(): SimplifiedPet[] {
    const start = (this.currentPage - 1) * this.itemsPerPage;
    const end = start + this.itemsPerPage;
    return this.animals.slice(start, end);
  }

  goToPage(page: number): void {
    if (page >= 1 && page <= this.totalPages) {
      this.currentPage = page;
    }
  }

  nextPage(): void {
    if (this.currentPage < this.totalPages) {
      this.currentPage++;
    }
  }

  previousPage(): void {
    if (this.currentPage > 1) {
      this.currentPage--;
    }
  }

  // ======================================
  // NAVIGATION
  // ======================================
  goToDashboard(): void {
    this.router.navigate(['/dashboard']);
  }

  // ======================================
  // FILTER MANAGEMENT
  // ======================================
  clearFilters(): void {
    this.filters = {
      name: '',
      status: '',
      species_id: '',
      gender: '',
      vaccinated: ''
    };
    this.currentPage = 1; // Reset to first page
    this.fetchPets(); // Reload all pets without filters
  }

  // ======================================
}
