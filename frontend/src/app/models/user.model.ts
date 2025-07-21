// User-related models for the frontend application
// These interfaces correspond to the Go models in the backend

export interface FullUser {
  id: number;
  name: string;
  surname: string;
  email: string;
  session_id: string;
  address: string;
  failed_logins: number;
  is_blocked: boolean;
  two_factor_auth: string;
  password?: string; // Optional since it's omitted from JSON
  provider: string;
  provider_id: string;
  change_password: boolean;
  is_admin: boolean;
  crt_date: string; // ISO date string
  upt_date: string; // ISO date string
}

export interface User {
  id: number;
  name: string;
  surname: string;
  email: string;
  session_id: string;
  address: string;
  provider: string;
  provider_id: string;
  password: string;
  change_pass: boolean;
  failed_logins: number;
  is_blocked: boolean;
  is_admin: boolean;
  crt_date: string; // ISO date string
  upt_date: string; // ISO date string
}

export interface NonValidatedUser {
  id: number;
  name: string;
  surname: string;
  email: string;
  address: string;
  failed_logins: number;
  provider: string;
  is_blocked: boolean;
  is_admin: boolean;
  crt_date: string; // ISO date string
  upt_date: string; // ISO date string
}

export interface SimplifiedUser {
  id: number;
  name: string;
  surname: string;
  email: string;
  address: string;
  is_admin: boolean;
}
