import { Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Card } from '../card/card';
import { ApiService } from '@app/services/api.service';
import { MobileFooter } from "../mobile-footer/mobile-footer";

interface Pet {
  id: number;
  name: string;
  species: any;
  breed?: string;
  description?: string;
  image_url?: string;
  birthdate?: string;
  status?: string;
  is_vaccinated?: boolean;
}

@Component({
  selector: 'app-card-list',
  imports: [CommonModule, Card, MobileFooter],
  templateUrl: './card-list.html',
  host: {
    'style': 'view-transition-name: card-list'
  }
})
export class CardList implements OnInit {

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================
  animals: ReadonlyArray<Pet> = [];
  species: ReadonlyArray<any> = [];
  error: string | null = null;
  currentPage = 1;

  @Input() showPagination = true;
  @Input() showFilter = true;
  @Input() itemsPerPage: number = 8; 

  // ======================================
  // CONSTRUCTOR
  // ======================================
  constructor(private apiService: ApiService) {
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

  private handlePetsSuccess(data: any): void {
    this.animals = data.content.map((pet: any) => ({
      ...pet,
      birthdate: pet.birthdate && pet.birthdate !== "0001-01-01T00:00:00Z"
        ? this.parseDateToDDMMYYYY(pet.birthdate)
        : ''
    }));
    console.log('Fetched pets:', data.content);
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

  private handleSpeciesSuccess(data: any): void {
    console.log('Fetched species:', data);
    this.species = data.content;
  }

  private handleSpeciesError(error: any): void {
    this.error = 'Error fetching pets.';
    console.error('Error fetching pets: ', error);
  }

  // ======================================
  // PAGINATION
  // ======================================
  get totalPages(): number {
    return Math.ceil(this.animals.length / this.itemsPerPage);
  }

  get paginatedAnimals(): Pet[] {
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
}
