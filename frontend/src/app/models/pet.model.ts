// Pet-related models for the frontend application
// These interfaces correspond to the Go models in the backend

import { SimplifiedUser } from './user.model';
import { Species } from './species.model';

export enum PetStatus {
  Available = 'Available',
  FosterHome = 'FosterHome',
  Adopted = 'Adopted'
}

export enum PetGender {
  Male = 'Male',
  Female = 'Female'
}

export interface VaccinationHistory {
  pet_id: number;
  vaccination_date: string; // ISO date string
  vaccine_name: string;
}

export interface Pet {
  id: number;
  name: string;
  species_id: number;
  species: Species;
  breed: string;
  gender: PetGender;
  weight: number;
  status: PetStatus;
  birthdate: string; // ISO date string
  adopt_date: string; // ISO date string
  description: string;
  adopt_user_id: number | null;
  adopt_user: SimplifiedUser | null;
  image_url: string;
  is_vaccinated: boolean;
  vaccination_history: VaccinationHistory[];
  is_urgent: boolean;
  crt_date: string; // ISO date string
  upt_date: string; // ISO date string
}

export interface SimplifiedPet {
  id: number;
  name: string;
  description: string;
  species_id: number;
  species: Species;
  breed: string;
  gender: PetGender;
  weight: number;
  status: PetStatus;
  adopt_user_id: number | null;
  adopt_user: SimplifiedUser | null;
  birthdate: string; // ISO date string
  image_url: string;
  is_vaccinated: boolean;
  vaccination_history: VaccinationHistory[];
  is_urgent: boolean;
}

export interface Base64Image {
	base64: string;
	name: string;
}