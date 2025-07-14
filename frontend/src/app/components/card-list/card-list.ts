import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Card } from '../card/card';import { HttpClient } from '@angular/common/http';
import { ApiService } from '@app/services/api.service';
;

interface Pet {
  id: number;
  name: string;
  species: string;
  breed?: string; // Optional field for breed
  isAdopted: boolean;
  description: string; // Optional field for description
  image: string; // URL to the pet's image
}

@Component({
  selector: 'app-card-list',
  imports: [CommonModule, Card],
  templateUrl: './card-list.html',
  styleUrl: './card-list.css',
})
export class CardList implements OnInit{
  constructor(
    private apiService: ApiService,
  ){}
  
  animals: Pet[] = [];
  currentPage = 1;
  itemsPerPage = 4;

  ngOnInit(): void {
    this.apiService.getPets().subscribe({
      next: (data) => {
        console.log('Pets fetched successfully:', data);
        this.animals = data;

      },
      error: (error) => {
        console.error('Error fetching pets: ', error)
      }
    })
  }

  get totalPages() {
    return Math.ceil(this.animals.length / this.itemsPerPage);
  }

  get paginatedAnimals() {
    const start = (this.currentPage - 1) * this.itemsPerPage;
    const end = start + this.itemsPerPage;
    return this.animals.slice(start, end);
  }
  nextPage() {
    if (this.currentPage < this.totalPages) {
      this.currentPage++;
    }
  }
  previousPage() {
    if (this.currentPage > 1) {
      this.currentPage--;
    }
  }
}
