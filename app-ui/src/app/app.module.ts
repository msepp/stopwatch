import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';

import { AppComponent } from './app.component';
import { AstilectronModule } from './astilectron';

// Store reducers
import { ActiveTaskReducer } from './store/reducers/active-task.reducers';
import { ProjectTasksReducer } from './store/reducers/project-tasks.reducers';
import { ProjectsReducer } from './store/reducers/projects.reducers';
import { SelectedProjectReducer } from './store/reducers/selected-project.reducers';
import { SelectedTaskReducer } from './store/reducers/selected-task.reducers';
import { VersionReducer } from './store/reducers/version.reducers';

// Store effects
import { ActiveTaskEffects } from './store/effects/active-task.effects';
import { ProjectTasksEffects } from './store/effects/project-tasks.effects';
import { ProjectsEffects } from './store/effects/projects.effects';
import { SelectedProjectEffects } from './store/effects/selected-project.effects';
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
      selectedProject: SelectedProjectReducer,
      selectedTask: SelectedTaskReducer,
      activeTask: ActiveTaskReducer,
      projects: ProjectsReducer,
      projectTasks: ProjectTasksReducer,
      version: VersionReducer
    }),
    EffectsModule.forRoot([
      ActiveTaskEffects,
      ProjectTasksEffects,
      ProjectsEffects,
      SelectedProjectEffects,
      SelectedTaskEffects
    ]),
    AstilectronModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
