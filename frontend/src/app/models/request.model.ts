/**
 * Request Models for Pet Adoption System
 * 
 * Defines TypeScript interfaces for API request payloads
 * related to pet adoption, foster home, and contact requests.
 */

/**
 * Pet adoption request payload
 * Used for submitting pet adoption requests to the organization
 */
export interface PetAdoptionRequest {
  pet_id: number;
  user_id: number;
}

/**
 * Foster home request payload
 * Used for submitting foster home requests to the organization
 */
export interface PetFosterHomeRequest {
  pet_id: number;
  user_id: number;
}

/**
 * Foster home contact request payload
 * Used for sending contact messages to foster families
 */
export interface PetFosterHomeContactRequest {
  pet_id: number;
  user_id: number;
  contact_user_id: number;
  reason: string;
  message: string;
}

/**
 * Generic API response for request operations
 */
export interface RequestResponse {
  status: string;
  message?: string;
}
