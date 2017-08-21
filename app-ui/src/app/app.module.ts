import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import { MaterialModule, MdDatepickerModule, MdNativeDateModule } from '@angular/material';
import { RouterModule } from '@angular/router';
import './rxjs-operators';

import { AppComponent } from './app.component';
import { AstilectronModule } from './astilectron';
import { TaskDurationPipe } from './task-duration.pipe';
import { StopwatchService } from './services/stopwatch.service';
import { ErrorService } from './services/error.service';

// Components
import { ActiveTaskComponent } from './components/active-task/active-task.component';
import { GroupListComponent } from './components/group-list/group-list.component';
import { AddGroupFormComponent } from './components/add-group-form/add-group-form.component';
import { GroupsHomeComponent } from './components/groups-home/groups-home.component';
import { GroupDetailsComponent } from './components/group-details/group-details.component';
import { AddTaskFormComponent } from './components/add-task-form/add-task-form.component';
import { TaskDetailsComponent } from './components/task-details/task-details.component';
import { TaskHistoryComponent } from './components/task-history/task-history.component';
import { TaskComponent } from './components/task/task.component';
import { GroupUsageComponent } from './components/group-usage/group-usage.component';

// Store reducers
import { ActiveTaskReducer } from './store/reducers/active-task.reducers';
import { BackendConnReducer } from './store/reducers/backend-conn.reducers';
import { GroupTasksReducer } from './store/reducers/group-tasks.reducers';
import { GroupsReducer } from './store/reducers/groups.reducers';
import { SelectedGroupReducer } from './store/reducers/selected-group.reducers';
import { SelectedTaskReducer } from './store/reducers/selected-task.reducers';
import { VersionReducer } from './store/reducers/version.reducers';
import { TaskHistoryReducer } from './store/reducers/task-history.reducers';

// Effects
import { GroupTasksEffects } from './store/effects/group-tasks.effects';
import { TaskHistoryEffects } from './store/effects/task-history.effects';

@NgModule({
  declarations: [
    AppComponent,
    ActiveTaskComponent,
    GroupListComponent,
    AddGroupFormComponent,
    GroupsHomeComponent,
    GroupDetailsComponent,
    TaskDurationPipe,
    AddTaskFormComponent,
    TaskDetailsComponent,
    TaskHistoryComponent,
    TaskComponent,
    GroupUsageComponent
  ],
  imports: [
    BrowserModule,
    ReactiveFormsModule,
    HttpModule,
    BrowserAnimationsModule,
    MaterialModule,
    MdDatepickerModule,
    MdNativeDateModule,
    RouterModule.forRoot([
      {path: 'groups', component: GroupsHomeComponent},
      {path: 'group/:id', component: GroupDetailsComponent},
      {path: 'task/:id', component: TaskDetailsComponent},
      {path: 'usage/:id', component: GroupUsageComponent},
      {path: '**', redirectTo: 'groups'}
    ]),
    StoreModule.forRoot({
      backendConn: BackendConnReducer,
      selectedGroup: SelectedGroupReducer,
      selectedTask: SelectedTaskReducer,
      activeTask: ActiveTaskReducer,
      groups: GroupsReducer,
      groupTasks: GroupTasksReducer,
      version: VersionReducer,
      taskHistory: TaskHistoryReducer
    }),
    EffectsModule.forRoot([
      GroupTasksEffects,
      TaskHistoryEffects
    ]),
    AstilectronModule
  ],
  providers: [
    StopwatchService,
    ErrorService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
