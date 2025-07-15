import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

export interface BreadcrumbItem {
  label: string;
  icon?: string;
  routerLink?: string;
  url?: string;
  isActive?: boolean;
}

@Component({
  selector: 'app-breadcrumb',
  imports: [CommonModule, RouterModule],
  templateUrl: './breadcrumb.html',
  standalone: true
})
export class Breadcrumb {
  @Input() items: BreadcrumbItem[] = [];
  
  /**
   * Método para determinar si un item es navegable
   */
  isNavigable(item: BreadcrumbItem, isLast: boolean): boolean {
    return !isLast && (!!item.routerLink || !!item.url);
  }
}
