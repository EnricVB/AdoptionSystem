import { Component } from '@angular/core';
import { InfoBox } from '../../components/info-box/info-box';
import { CallToAction } from "@app/components/call-to-action/call-to-action";
import { Wave } from "@app/components/wave/wave";
import { Footer } from "@app/components/footer/footer";
import { Router } from '@angular/router';
import { PetHeader } from '@app/components/header/header';

@Component({
  selector: 'app-dashboard',
  imports: [InfoBox, CallToAction, Wave, Footer, PetHeader],
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
      private router: Router
  ) {
  }
}