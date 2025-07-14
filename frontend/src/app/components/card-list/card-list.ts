import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Card } from '../card/card';
import { ApiService } from '@app/services/api.service';

interface Pet {
  id: number;
  name: string;
  species: any;
  breed?: string;
  isAdopted: boolean;
  description?: string;
  imageURL?: string;
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
  loading = false;
  error: string | null = null;
  currentPage = 1;
  itemsPerPage = 4;

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
    this.loading = true;
    this.error = null;
    this.apiService.getPets().subscribe({
      next: (data) => this.handlePetsSuccess(data),
      error: (err) => this.handlePetsError(err)
    });
  }

  private handlePetsSuccess(data: any): void {
    this.animals = data.content;
    this.loading = false;
  }

  private handlePetsError(error: any): void {
    this.loading = false;
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
