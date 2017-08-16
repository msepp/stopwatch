import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';

import { AppComponent } from './app.component';
import { AstilectronModule } from './astilectron';
import { StopwatchService } from './services/stopwatch.service';

import './rxjs-operators';

// Store reducers
import { ActiveTaskReducer } from './store/reducers/active-task.reducers';
import { BackendConnReducer } from './store/reducers/backend-conn.reducers';
import { GroupTasksReducer } from './store/reducers/group-tasks.reducers';
import { GroupsReducer } from './store/reducers/groups.reducers';
import { SelectedGroupReducer } from './store/reducers/selected-group.reducers';
import { SelectedTaskReducer } from './store/reducers/selected-task.reducers';
import { VersionReducer } from './store/reducers/version.reducers';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    ReactiveFormsModule,
    HttpModule,
    BrowserAnimationsModule,
    StoreModule.forRoot({
      backendConn: BackendConnReducer,
      selectedGroup: SelectedGroupReducer,
      selectedTask: SelectedTaskReducer,
      activeTask: ActiveTaskReducer,
      groups: GroupsReducer,
      groupTasks: GroupTasksReducer,
      version: VersionReducer
    }),
    EffectsModule.forRoot([]),
    AstilectronModule
  ],
  providers: [
    StopwatchService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
