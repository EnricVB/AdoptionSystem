import { Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Card } from '../card/card';
import { ApiService } from '@app/services/api.service';

interface Pet {
  id: number;
  name: string;
  species: any;
  breed?: string;
  is_adopted: boolean;
  description?: string;
  image_url?: string;
}

@Component({
  selector: 'app-card-list',
  imports: [CommonModule, Card],
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

  private handlePetsSuccess(data: any): void {
    console.log('Fetched pets:', data);
    this.animals = data.content;
  }

  private handlePetsError(error: any): void {
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
