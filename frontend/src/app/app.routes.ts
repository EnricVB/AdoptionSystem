import { Routes } from '@angular/router';
import { Dashboard } from './pages/dashboard/dashboard';
import { Login } from './pages/login/login';
import { Logout } from './pages/logout/logout';
import { RecoverPassword } from './pages/login/recover-password/recover-password';
import { ChangePass } from './pages/login/change-pass/change-pass';
import { Register } from './pages/register/register';
import { Twofa } from './pages/twofa/twofa';

export const routes: Routes = [
    { 
        path: '', 
        component: Dashboard,
        data: { animation: 'dashboard' }
    },{ 
        path: 'dashboard', 
        component: Dashboard,
        data: { animation: 'dashboard' }
    },
    {
        path: 'logout',
        component: Logout
    },
    { 
        path: 'login', 
        component: Login,
        data: { animation: 'login' }
    },
    {
        path: 'recover-password',
        component: RecoverPassword,
        data: { animation: 'recover-password' }
    },
    {
        path: 'change-pass',
        component: ChangePass,
        data: { animation: 'change-pass' }
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
];
