import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { App } from './app/app';
import { CardList } from '@app/components/card-list/card-list';

bootstrapApplication(CardList, appConfig)
  .catch((err) => console.error(err));
