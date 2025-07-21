import { ChangeDetectorRef, Component, ElementRef, Input, ViewChild } from '@angular/core';
import {
  CarouselComponent,
  CarouselInnerComponent,
  CarouselItemComponent
} from '@coreui/angular';

@Component({
  selector: 'app-carousel',
  imports: [CarouselComponent, CarouselInnerComponent, CarouselItemComponent],
  templateUrl: './carousel.html',
})
export class Carousel {

  // ======================================
  // CHILD COMPONENTS
  // ======================================

  @ViewChild('carousel') carousel!: CarouselComponent;
  @ViewChild('modal') modal!: ElementRef<HTMLElement>;
  @ViewChild('modalImage') modalImage!: ElementRef<HTMLImageElement>;
  @ViewChild('file') fileInput!: ElementRef<HTMLInputElement>;

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================
  
  @Input() images: string[] = [];
  @Input() isEditable: boolean = false;
  @Input() selectedFiles: File[] = [];

  // ======================================
  // CONSTRUCTOR
  // ======================================
  constructor(
    private cdr: ChangeDetectorRef,
  ) {
    // Initialization logic can go here if needed
  }

  // ======================================
  // PUBLIC METHODS
  // ======================================

  public clearImages(): void {
    this.images = [];
    this.selectedFiles = [];
    this.fileInput.nativeElement.value = '';
    this.cdr.detectChanges();
  }

  public deleteImage(indexToDelete: number): void {
    if (this.images.length === 0 || indexToDelete < 0 || indexToDelete >= this.images.length) {
      return;
    }
    
    // Remove from arrays
    this.selectedFiles.splice(indexToDelete, 1);
    this.images.splice(indexToDelete, 1);

    // Only adjust carousel index if there are still images
    if (this.images.length > 0) {
      // Calculate safe new index
      const currentIndex = this.carousel?.activeIndex() || 0;
      let newIndex: number;
      
      // If we deleted the currently active image
      if (indexToDelete === currentIndex) {
        // If we deleted the last image, go to the previous one
        newIndex = indexToDelete >= this.images.length ? this.images.length - 1 : indexToDelete;
      } else if (indexToDelete < currentIndex) {
        // If we deleted an image before the current one, adjust index
        newIndex = currentIndex - 1;
      } else {
        // If we deleted an image after the current one, keep the same index
        newIndex = currentIndex;
      }
      
      // Ensure index is within bounds
      newIndex = Math.max(0, Math.min(newIndex, this.images.length - 1));
      
      setTimeout(() => {
        if (this.carousel) {
          this.carousel.activeIndex.set(newIndex);
        }
      }, 50);
    } else {
      // No images left, reset to 0
      setTimeout(() => {
        if (this.carousel) {
          this.carousel.activeIndex.set(0);
        }
      }, 50);
    }
    
    // Force change detection
    this.cdr.detectChanges();
  }

  onImageClick(index: number): void {
    this.modalImage.nativeElement.src = this.images[index];
    this.modal.nativeElement.hidden = false;
  }

  closeImage(): void {
    this.modal.nativeElement.hidden = true;
  }

  onImageError(event: any): void {
    // Replace broken image with a placeholder
    event.target.style.backgroundImage = 'url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjQwMCIgdmlld0JveD0iMCAwIDQwMCA0MDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSI0MDAiIGhlaWdodD0iNDAwIiBmaWxsPSIjRjNGNEY2Ii8+CjxwYXRoIGQ9Ik0yMDAgMTAwQzE0NC43NzIgMTAwIDEwMCAxNDQuNzcyIDEwMCAyMDBTMTQ0Ljc3MiAzMDAgMjAwIDMwMFMyNTAgMjU1LjIyOCAyNTAgMjAwUzIwNS4yMjggMTAwIDIwMCAxMDBaTTIwMCAyNTBDMTcyLjM4NiAyNTAgMTUwIDIyNy42MTQgMTUwIDIwMEMxNTAgMTcyLjM4NiAxNzIuMzg2IDE1MCAyMDAgMTUwQzIyNy42MTQgMTUwIDI1MCAxNzIuMzg2IDI1MCAyMDBDMjUwIDIyNy42MTQgMjI3LjYxNCAyNTAgMjAwIDI1MFoiIGZpbGw9IiM5Q0EzQUYiLz4KPC9zdmc+)';
  }

  // ======================================
  // CAROUSEL NAVIGATION METHODS
  // ======================================
  getSafeActiveIndex(): number {
    if (this.images.length === 0) {
      return 0;
    }
    const activeIndex = this.carousel?.activeIndex() || 0;
    return Math.max(0, Math.min(activeIndex, this.images.length - 1));
  }

  goToPreviousSlide(): void {
    if (this.images.length === 0) {
      return;
    }
    
    const currentIndex = this.carousel.activeIndex() || 0;
    const newIndex = currentIndex > 0 ? currentIndex - 1 : this.images.length - 1;
    
    // Validate index is within bounds
    if (newIndex >= 0 && newIndex < this.images.length) {
      this.carousel.activeIndex.set(newIndex);
      this.cdr.detectChanges();
    }
  }

  goToNextSlide(): void {
    if (this.images.length === 0) {
      return;
    }
    
    const currentIndex = this.carousel.activeIndex() || 0;
    const newIndex = currentIndex < this.images.length - 1 ? currentIndex + 1 : 0;
    
    // Validate index is within bounds
    if (newIndex >= 0 && newIndex < this.images.length) {
      this.carousel.activeIndex.set(newIndex);
      this.cdr.detectChanges();
    }
  }

  
  onFileSelected(event: any): void { 
    const file = event.target.files[0];
    const reader = new FileReader();

    if (!file) {
      return;
    }

    // Add the file to selected files
    this.selectedFiles.push(file);

    reader.onload = (e) => {
      const imageUrl = e.target?.result as string;
      
      // Add the image URL to carousel images
      this.images.push(imageUrl);
      
      // Move carousel to the last image (newly added)
      const lastIndex = this.images.length - 1;
      
      // Use setTimeout to ensure the DOM is updated before changing the index
      setTimeout(() => {
        if (this.carousel && lastIndex >= 0) {
          this.carousel.activeIndex.set(lastIndex);
        }
        // Force change detection to update the view
        this.cdr.detectChanges();
      }, 50);
    };

    reader.onerror = () => {
      // Remove the file from selected files if reading fails
      this.selectedFiles.pop();
    };

    // Read the file as a data URL
    reader.readAsDataURL(file);
    
    // Clear the input to allow selecting the same file again
    event.target.value = '';
  }
}
