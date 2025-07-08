import { Routes } from '@angular/router';
import { Dashboard } from './pages/dashboard/dashboard';
import { Login } from './pages/login/login';
import { Register } from './pages/register/register';
import { Twofa } from './pages/twofa/twofa';

export const routes: Routes = [
    { 
        path: '', 
        component: Login,
        data: { animation: 'login' }
    },
    { 
        path: 'login', 
        component: Login,
        data: { animation: 'login' }
    },
    { 
        path: 'register', 
        component: Register,
        data: { animation: 'register' }
    },
    { 
        path: 'twofa', 
        component: Twofa,
        data: { animation: 'twofa' }
    },
    { 
        path: 'dashboard', 
        component: Dashboard,
        data: { animation: 'dashboard' }
    },
];
