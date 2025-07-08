import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

/**
 * API Service
 * 
 * Centralizes all API endpoints and HTTP operations for the adoption system.
 * Provides a clean interface for components to interact with the backend.
 * 
 * Features:
 * - User authentication and management
 * - Pet management operations
 * - Species management operations
 * - Consistent error handling
 * - Centralized configuration
 */
@Injectable({
  providedIn: 'root'
})
export class ApiService {
  // ========================================
  // CONFIGURATION
  // ========================================
  
  /** Base URL for all API calls */
  private readonly baseUrl = '/api';
  
  /** Default HTTP headers for all requests */
  private readonly defaultHeaders = new HttpHeaders({
    'Content-Type': 'application/json'
  });

  // ========================================
  // CONSTRUCTOR
  // ========================================
  
  constructor(private http: HttpClient) {}

  // ========================================
  // AUTHENTICATION ENDPOINTS
  // ========================================
  
  /**
   * Authenticate user with email and password
   * 
   * @param credentials - User login credentials
   * @returns Observable with authentication response containing session_id
   */
  login(credentials: { email: string; password: string }): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/auth/login`, 
      credentials, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Register a new user account
   * 
   * @param userData - User registration data
   * @returns Observable with registration response
   */
  register(userData: { 
    Name: string; 
    Surname: string; 
    Email: string; 
    Password: string 
  }): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/register`, 
      userData, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Authenticate user with Google OAuth
   * 
   * @param googleData - Google authentication data
   * @returns Observable with authentication response
   */
  loginWithGoogle(googleData: { email: string; id_token: string }): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/auth/login/google`, 
      googleData, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Verify 2FA code for authenticated user
   * 
   * @param twoFactorData - Session ID and 2FA code
   * @returns Observable with verification response
   */
  verify2FA(twoFactorData: { session_id: string; code: string }): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/auth/verify-2fa`, 
      twoFactorData, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Refresh 2FA token - sends new verification code to user's email
   * 
   * @param email - User's email address
   * @returns Observable with refresh response
   */
  refresh2FAToken(email: { email: string }): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/auth/refresh-token`, 
      email, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Reset user password with a random one
   * 
   * @param emailData - User's email address
   * @returns Observable with recovery response
   */
  resetPassword(emailData: { email: string }): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/auth/reset-password`, 
      emailData, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Send email with password recovery instructions
   * 
   * @param emailData - User's email address
   * @returns Observable with email sending response
   */
  sendRecoverPasswordMail(emailData: { email: string }): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/auth/forgot-password`, 
      emailData, 
      { headers: this.defaultHeaders }
    );
  }

  // ========================================
  // USER MANAGEMENT ENDPOINTS
  // ========================================
  
  /**
   * Get all users in the system
   * 
   * @returns Observable with array of all users
   */
  getUsers(): Observable<any> {
    return this.http.get<any>(
      `${this.baseUrl}/users`, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Get a specific user by ID
   * 
   * @param userId - Unique identifier of the user
   * @returns Observable with user data
   */
  getUserById(userId: number): Observable<any> {
    return this.http.get<any>(
      `${this.baseUrl}/users/${userId}`, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Create a new user account
   * 
   * @param userData - User information for account creation
   * @returns Observable with creation response
   */
  createUser(userData: {
    name: string;
    surname: string;
    email: string;
    address?: string;
    password?: string;
    provider?: string;
    provider_id?: string;
  }): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/users`, 
      userData, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Update an existing user's information
   * 
   * @param userId - Unique identifier of the user
   * @param userData - Updated user information
   * @returns Observable with update response
   */
  updateUser(userId: number, userData: any): Observable<any> {
    return this.http.put<any>(
      `${this.baseUrl}/users/${userId}`, 
      userData, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Delete a user account
   * 
   * @param userId - Unique identifier of the user to delete
   * @returns Observable with deletion response
   */
  deleteUser(userId: number): Observable<any> {
    return this.http.delete<any>(
      `${this.baseUrl}/users/${userId}`, 
      { headers: this.defaultHeaders }
    );
  }

  // ========================================
  // PET MANAGEMENT ENDPOINTS
  // ========================================
  
  /**
   * Get all pets in the system
   * 
   * @returns Observable with array of all pets
   */
  getPets(): Observable<any> {
    return this.http.get<any>(
      `${this.baseUrl}/pets`, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Get a specific pet by ID
   * 
   * @param petId - Unique identifier of the pet
   * @returns Observable with pet data
   */
  getPetById(petId: number): Observable<any> {
    return this.http.get<any>(
      `${this.baseUrl}/pets/${petId}`, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Create a new pet record
   * 
   * @param petData - Pet information for registration
   * @returns Observable with creation response
   */
  createPet(petData: {
    name: string;
    species: string;
    breed?: string;
    birth_date?: string;
    description?: string;
  }): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/pets`, 
      petData, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Update an existing pet's information
   * 
   * @param petId - Unique identifier of the pet
   * @param petData - Updated pet information
   * @returns Observable with update response
   */
  updatePet(petId: number, petData: any): Observable<any> {
    return this.http.put<any>(
      `${this.baseUrl}/pets/${petId}`, 
      petData, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Delete a pet record
   * 
   * @param petId - Unique identifier of the pet to delete
   * @returns Observable with deletion response
   */
  deletePet(petId: number): Observable<any> {
    return this.http.delete<any>(
      `${this.baseUrl}/pets/${petId}`, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Adopt a pet (assign to a user)
   * 
   * @param petId - Unique identifier of the pet to adopt
   * @param userId - Unique identifier of the adopting user
   * @returns Observable with adoption response
   */
  adoptPet(petId: number, userId: number): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/pets/${petId}/adopt`, 
      { user_id: userId }, 
      { headers: this.defaultHeaders }
    );
  }

  // ========================================
  // SPECIES MANAGEMENT ENDPOINTS
  // ========================================
  
  /**
   * Get all species in the system
   * 
   * @returns Observable with array of all species
   */
  getSpecies(): Observable<any> {
    return this.http.get<any>(
      `${this.baseUrl}/species`, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Get a specific species by ID
   * 
   * @param speciesId - Unique identifier of the species
   * @returns Observable with species data
   */
  getSpeciesById(speciesId: number): Observable<any> {
    return this.http.get<any>(
      `${this.baseUrl}/species/${speciesId}`, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Create a new species
   * 
   * @param speciesData - Species information
   * @returns Observable with creation response
   */
  createSpecies(speciesData: { name: string }): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/species`, 
      speciesData, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Update an existing species
   * 
   * @param speciesId - Unique identifier of the species
   * @param speciesData - Updated species information
   * @returns Observable with update response
   */
  updateSpecies(speciesId: number, speciesData: { name: string }): Observable<any> {
    return this.http.put<any>(
      `${this.baseUrl}/species/${speciesId}`, 
      speciesData, 
      { headers: this.defaultHeaders }
    );
  }

  /**
   * Delete a species
   * 
   * @param speciesId - Unique identifier of the species to delete
   * @returns Observable with deletion response
   */
  deleteSpecies(speciesId: number): Observable<any> {
    return this.http.delete<any>(
      `${this.baseUrl}/species/${speciesId}`, 
      { headers: this.defaultHeaders }
    );
  }

  // ========================================
  // UTILITY METHODS
  // ========================================
  
  /**
   * Get custom headers for specific requests
   * 
   * @param additionalHeaders - Additional headers to merge with defaults
   * @returns HttpHeaders object with merged headers
   */
  private getHeaders(additionalHeaders?: { [key: string]: string }): HttpHeaders {
    let headers = this.defaultHeaders;
    
    if (additionalHeaders) {
      Object.keys(additionalHeaders).forEach(key => {
        headers = headers.set(key, additionalHeaders[key]);
      });
    }
    
    return headers;
  }

  /**
   * Build full URL for endpoint
   * 
   * @param endpoint - API endpoint path
   * @returns Full URL string
   */
  private buildUrl(endpoint: string): string {
    return `${this.baseUrl}${endpoint.startsWith('/') ? endpoint : '/' + endpoint}`;
  }
}
