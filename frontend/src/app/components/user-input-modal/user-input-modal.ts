import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { User } from '@app/models/user.model';
import { ApiService } from '@app/services/api.service';
import {
  ButtonDirective,
  ModalBodyComponent,
  ModalComponent,
  ModalFooterComponent,
  ModalHeaderComponent,
  ModalTitleDirective
} from '@coreui/angular';

@Component({
  selector: 'app-user-input-modal',
  templateUrl: './user-input-modal.html',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ModalComponent,
    ModalHeaderComponent,
    ModalTitleDirective,
    ModalBodyComponent,
    ModalFooterComponent,
    ButtonDirective
  ]
})
export class UserInputModal {

  // ============================
  // COMPONENT PROPERTIES
  // ============================
  public visible = false;
  public selectedUserId: number | null = null;
  public users: ReadonlyArray<User> = [];
  public filteredUsers: Array<User> = [];
  nameFilter = "";

  @Output() onCancel = new EventEmitter<number>();
  @Output() onSaveChanges = new EventEmitter<number>();
  
  // ======================================
  // CONSTRUCTOR
  // ======================================
  constructor(
    private apiService: ApiService,
  ) {
  }

  // ======================================
  // LIFECYCLE
  // ======================================
  ngOnInit(): void {
    this.fetchUserData();
  }

  // ======================================
  // DATA FETCHING & HANDLING
  // ======================================
  fetchUserData(): void { 
    this.apiService.getUsers().subscribe({
      next: (data) => this.handleUserSuccess(data),
      error: (err) => this.handleUserError(err)
    });
  }

  private handleUserSuccess(data: { content: User[] }): void {
    this.users = data.content;
    this.filteredUsers = data.content;
  }

  private handleUserError(error: any): void {
    console.error('Error fetching users: ', error);
  }

  // ============================
  // SELECTION METHODS 
  // ============================
  
  selectUser(userId: number) {
    this.selectedUserId = userId;
  }

  filterUsers() {
    this.filteredUsers = this.users.filter(
      user =>
        user.name.includes(this.nameFilter) || user.email.includes(this.nameFilter) || this.nameFilter === ""
    );
  }

  // ============================
  // EVENT LISTENERS
  // ============================

  showModal() {
    this.visible = true;
  }

  onCloseListener() {
    this.visible = false;
    this.onCancel.emit(0);
  }

  onSaveChangesListener() {
    this.visible = false;
    this.onSaveChanges.emit(this.selectedUserId ?? 0);
  }
}