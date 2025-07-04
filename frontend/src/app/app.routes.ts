import { Routes } from '@angular/router';
import { Dashboard } from './pages/dashboard/dashboard';
import { Login } from './pages/login/login';
import { Twofa } from './pages/twofa/twofa';

export const routes: Routes = [
    {path: '', component: Login },
    {path:'twofa', component: Twofa},
    {path: 'dashboard', component: Dashboard},
];
