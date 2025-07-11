import { Component } from '@angular/core';
import { InfoBox } from '../../components/info-box/info-box';
import { CallToAction } from "@app/components/call-to-action/call-to-action";
import { Wave } from "@app/components/wave/wave";
import { Footer } from "@app/components/footer/footer";

@Component({
  selector: 'app-dashboard',
  imports: [InfoBox, CallToAction, Wave, Footer],
  templateUrl: './dashboard.html',
  host: {
    'style': 'view-transition-name: dashboard-page'
  }
})
export class Dashboard {

}
