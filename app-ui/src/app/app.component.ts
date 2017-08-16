import { Component, OnInit } from '@angular/core';
import { Store } from '@ngrx/store';
import { AppState } from './model/app-state';
import { Observable } from 'rxjs/Observable';

import * as VersionActions from './store/actions/version.actions';
import * as GroupsActions from './store/actions/groups.actions';
import * as ActiveTaskActions from './store/actions/active-task.actions';
import * as GroupTasksActions from './store/actions/group-tasks.actions';
import * as SelectedGroupActions from './store/actions/selected-group.actions';

import * as message from './astilectron/message';
import { FormControl, FormGroup, FormBuilder, Validators } from '@angular/forms';

import { AppVersion, Group, Task } from './model';
import { StopwatchService } from './services/stopwatch.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.less']
})
export class AppComponent implements OnInit {
  public ready = false;

  constructor(
    private stopwatch: StopwatchService,
    private store: Store<AppState>
  ) {
    this.store.select('backendConn')
      .filter(ok => (ok === true))
      .take(1)
      .concatMap(() => this.stopwatch.loadGroups())
      .concatMap(() => this.stopwatch.loadActiveTask())
      .subscribe(
        () => { console.log('ready'); this.ready = true; },
        e => console.log('error:', e)
      );
  }

  public ngOnInit() {
    this.stopwatch.init();
  }
}
