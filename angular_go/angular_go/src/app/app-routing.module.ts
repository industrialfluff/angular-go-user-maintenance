import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AboutComponent } from './about/about.component';
import { NewsComponent } from './news/news.component';
import { PageNotFoundComponent } from './pagenotfound/pagenotfound.component';
import { UserlistComponent } from './userlist/userlist.component';
import { UserComponent } from './user/user.component';
const routes: Routes = [
    { path: 'news', component: NewsComponent },
    { path: 'about', component: AboutComponent },
    { path: 'userlist', component: UserlistComponent },
    { path: 'user/new', component: UserComponent },
    { path: 'user/:user_id/edit', component: UserComponent },
    { path: 'user/:user_id/delete', component: UserlistComponent },
    { path: '', redirectTo: '/news', pathMatch: 'full' },
    { path: '**', component: PageNotFoundComponent },  // ALWAYS MAKE SURE THIS IS LAST
  ];


@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
