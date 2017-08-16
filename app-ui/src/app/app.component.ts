import { Component, OnInit } from '@angular/core';
import { Astilectron, Message } from './astilectron';
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

  public groupID: number;
  public taskID: number;

  constructor(
    private asti: Astilectron,
    private fb: FormBuilder,
    private store: Store<AppState>
  ) {
    this.newGroup = this.fb.group({
      name: ['', Validators.required]
    });

    this.newTask = this.fb.group({
      name: ['', Validators.required],
      costcode: ''
    });

    this.versions = this.store.select('version');
    this.groups = this.store.select('groups');
    this.tasks = this.store.select('groupTasks');
    this.activeTask = this.store.select('activeTask');
    this.selectedGroup = this.store.select('selectedGroup');

    this.selectedGroup.subscribe((v: Group) => {
      if (v) {
        this.groupID = v.id;
      }
    });

    this.activeTask.subscribe((v: Task) => {
      if (v) {
        console.log('active task changed to', v);
        this.taskID = v.id;
      }
    });
  }

  public ngOnInit() {
    this.asti.isReady.filter(v => v === true).take(1).subscribe(() => this.getVersions());
  }

  public addGroup() {
    if (this.newGroup.valid) {
      this.asti.send(message.REQUEST_ADD_GROUP, this.newGroup.value).subscribe(
        (m: Message) => {
          const v: Group = Object.assign(new Group, m.data);
          this.store.dispatch(new GroupsActions.Add(v));
          this.newGroup.reset();
        },
        (e: Error) => {
          console.log('whoops:', e.message);
        }
      );
    }
  }

  public addTask() {
    if (this.newTask.valid) {
      const payload = {
        groupid: this.groupID,
        name: this.newTask.get('name').value,
        costcode: this.newTask.get('costcode').value
      };

      this.asti.send(message.REQUEST_ADD_TASK, payload).subscribe(
        (m: Message) => {
          const v: Task = Object.assign(new Task, m.data);
          this.store.dispatch(new GroupTasksActions.Add(v));
          this.newTask.reset();
        },
        (e: Error) => {
          console.log('whoops:', e.message);
        }
      );
    }
  }

  public selectGroup(p: Group) {
    this.asti.send(message.REQUEST_GET_GROUP_TASKS, {groupid: p.id}).subscribe(
      (m: Message) => {
        const tasks: Task[] = [];
        m.data.forEach(v => {
          console.log(v);
          const t: Task = Object.assign(new Task, v);
          tasks.push(t);
        });

        this.store.dispatch(new GroupTasksActions.Set(tasks));
        this.store.dispatch(new SelectedGroupActions.Set(p));
      },
      (e: Error) => {
        console.log('whoops:', e.message);
      }
    );
  }

  public toggleTask(t: Task) {
    const payload = {
      groupid: t.groupid,
      id: t.id
    };

    if ((!!t.running) === false) {
      this.asti.send(message.REQUEST_START_TASK, payload).subscribe(
        (m: Message) => {
          const nt: Task = Object.assign(new Task, m.data);
          this.store.dispatch(new GroupTasksActions.Update(nt));
          this.store.dispatch(new ActiveTaskActions.Set(nt));
        },
        (e: Error) => {
          console.log('whoops:', e.message);
        }
      );
    } else {
      this.asti.send(message.REQUEST_STOP_TASK, payload).subscribe(
        (m: Message) => {
          const nt: Task = Object.assign(new Task, m.data);
          this.store.dispatch(new GroupTasksActions.Update(nt));
          this.store.dispatch(new ActiveTaskActions.Clear());
        },
        (e: Error) => {
          console.log('whoops:', e.message);
        }
      );
    }
  }

  private getVersions() {
    this.asti.send(message.REQUEST_APP_VERSIONS, null).subscribe(
      (m: Message) => {
        const v: AppVersion = Object.assign(new AppVersion, m.data);
        this.store.dispatch(new VersionActions.Set(v));

        this.openDatabase();
      },
      (e: Error) => {
        console.log('whoops:', e.message);
      }
    );
  }

  private openDatabase() {
    this.asti.send(message.REQUEST_OPEN_DATABASE, null).subscribe(
      () => this.getGroups(),
      (e: Error) => {
        console.log('whoops:', e.message);
      }
    );
  }

  private getGroups() {
    this.asti.send(message.REQUEST_GROUPS, null).subscribe(
      (m: Message) => {
        console.log(m.data);
        const all: Group[] = [];
        m.data.forEach(v => {
          all.push(Object.assign(new Group, v));
        });

        this.store.dispatch(new GroupsActions.Set(all));
      },
      (e: Error) => {
        console.log('whoops:', e.message);
      }
    );
  }
}
