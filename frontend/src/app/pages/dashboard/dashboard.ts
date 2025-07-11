import { Component } from '@angular/core';
import { InfoBox } from '../../components/info-box/info-box';
import { CallToAction } from "@app/components/call-to-action/call-to-action";
import { Wave } from "@app/components/wave/wave";
import { Footer } from "@app/components/footer/footer";
import { Router } from '@angular/router';
import { PetHeader } from '@app/components/header/header';
import { DeviceDetectorService } from 'ngx-device-detector';
import { MobileFooter } from '@app/components/mobile-footer/mobile-footer';

@Component({
  selector: 'app-dashboard',
  imports: [InfoBox, CallToAction, Wave, Footer, PetHeader, MobileFooter],
  templateUrl: './dashboard.html',
  host: {
    'style': 'view-transition-name: dashboard-page'
  }
})
export class Dashboard {
  
  // ======================================
  // CONSTRUCTOR
  // ======================================

  constructor(
    private router: Router,
    private deviceService: DeviceDetectorService
  ) {
    this.ngOnInit();
  }

  // ======================================
  // PROPERTIES
  // ======================================

  public isMobile!: boolean;

  // ======================================
  // LIFECYCLE HOOKS
  // ======================================

  ngOnInit() {
    this.isMobile = this.deviceService.isMobile();
  }
}