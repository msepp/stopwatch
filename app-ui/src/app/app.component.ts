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
  public title = 'Stopwatch';
  public groups: Observable<Group[]>;
  public versions: Observable<AppVersion>;
  public activeTask: Observable<Task>;
  public selectedGroup: Observable<Group>;
  public tasks: Observable<Task[]>;
  public newGroup: FormGroup;
  public newTask: FormGroup;
  public backendReady: Observable<boolean>;

  constructor(
    private stopwatch: StopwatchService,
    private fb: FormBuilder,
    private store: Store<AppState>
  ) {
    this.newGroup = this.fb.group({
      name: ['', Validators.required]
    });

    this.newTask = this.fb.group({
      groupid: 0,
      name: ['', Validators.required],
      costcode: ''
    });

    this.backendReady = this.store.select('backendConn');
    this.versions = this.store.select('version');
    this.groups = this.store.select('groups');
    this.tasks = this.store.select('groupTasks');
    this.activeTask = this.store.select('activeTask');
    this.selectedGroup = this.store.select('selectedGroup');

    this.selectedGroup.subscribe((v: Group) => {
      if (v) {
        this.newTask.get('groupid').setValue(v.id);
      }
    });

    this.backendReady.filter(ok => (ok === true)).take(1).subscribe(
      (v) => this.stopwatch.loadGroups()
    );
  }

  public ngOnInit() {
    this.stopwatch.init();
  }

  public addGroup() {
    if (this.newGroup.valid) {
      this.stopwatch.addGroup(this.newGroup.value).subscribe(
        () => this.newGroup.reset()
      );
    }
  }

  public addTask() {
    if (this.newTask.valid) {
      this.stopwatch.addTask(this.newTask.value).subscribe(
        () => this.newTask.reset()
      );
    }
  }

  public selectGroup(g: Group) {
    this.stopwatch.selectGroup(g).subscribe((grp: Group) => {
      console.log('selected', grp);
    }, (e: Error) => console.log('error', e));
  }

  public toggleTask(t: Task) {
    if (!!t.running) {
      this.stopwatch.stopTask(t).subscribe(() => {}, e => console.log(e));
    } else {
      this.stopwatch.startTask(t).subscribe(() => {}, e => console.log(e));
    }
  }
}
