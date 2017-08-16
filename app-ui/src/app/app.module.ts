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

// Store reducers
import { ActiveTaskReducer } from './store/reducers/active-task.reducers';
import { GroupTasksReducer } from './store/reducers/group-tasks.reducers';
import { GroupsReducer } from './store/reducers/groups.reducers';
import { SelectedGroupReducer } from './store/reducers/selected-group.reducers';
import { SelectedTaskReducer } from './store/reducers/selected-task.reducers';
import { VersionReducer } from './store/reducers/version.reducers';

// Store effects
import { ActiveTaskEffects } from './store/effects/active-task.effects';
import { GroupTasksEffects } from './store/effects/group-tasks.effects';
import { GroupsEffects } from './store/effects/groups.effects';
import { SelectedGroupEffects } from './store/effects/selected-group.effects';
import { SelectedTaskEffects } from './store/effects/selected-task.effects';
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
      selectedGroup: SelectedGroupReducer,
      selectedTask: SelectedTaskReducer,
      activeTask: ActiveTaskReducer,
      groups: GroupsReducer,
      groupTasks: GroupTasksReducer,
      version: VersionReducer
    }),
    EffectsModule.forRoot([
      ActiveTaskEffects,
      GroupTasksEffects,
      GroupsEffects,
      SelectedGroupEffects,
      SelectedTaskEffects
    ]),
    AstilectronModule
  ],
  providers: [
    StopwatchService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
