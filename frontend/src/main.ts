import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { App } from './app/app';
import { Login } from './app/pages/login/login';

bootstrapApplication(Login, appConfig)
  .catch((err) => console.error(err));
