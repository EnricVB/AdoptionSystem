import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { App } from './app/app';
import { Card } from './app/components/card/card';

bootstrapApplication(Card, appConfig)
  .catch((err) => console.error(err));
