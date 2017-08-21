import { Injectable } from '@angular/core';
import { Astilectron } from '../astilectron';
import * as messaging from '../astilectron/message';
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { ReplaySubject } from 'rxjs/ReplaySubject';

import { Store } from '@ngrx/store';
import { ErrorService } from './error.service';
import { AppState, Group, Task } from '../model';
import * as ActiveTaskActions from '../store/actions/active-task.actions';
import * as BackendConnActions from '../store/actions/backend-conn.actions';
import * as GroupsActions from '../store/actions/groups.actions';
import * as GroupTasksActions from '../store/actions/group-tasks.actions';
import * as SelectedGroupActions from '../store/actions/selected-group.actions';
import * as SelectedTaskActions from '../store/actions/selected-task.actions';
import * as TaskHistoryActions from '../store/actions/task-history.actions';

@Injectable()
export class StopwatchService {
  private _ready = false;
  private _dbopen$: BehaviorSubject<boolean>;
  private _selectedGroup: Group = null;
  private _activeTask: Task = null;

  constructor(
    private backend: Astilectron,
    private store: Store<AppState>,
    private err: ErrorService
  ) {
    this._dbopen$ = new BehaviorSubject<boolean>(false);

    this.store.select('selectedGroup').subscribe(g => this._selectedGroup = g);
    this.store.select('activeTask').subscribe(t => this._activeTask = t);

    // Save history when it has content. Skip first (initial value) and the
    // first loaded value (don't want to re-save it.)
    this.store.select('taskHistory')
      .skip(2)
      .filter((list: Task[]) => list.length > 0)
      .debounceTime(1000)
      .subscribe(
        (history: Task[]) => this.saveTaskHistory(history).subscribe(
          () => {console.log('saved history:', history); },
          (e: Error) => this.err.log(e)
        )
      );
  }

  public get ready(): Observable<boolean> {
    return this._dbopen$.asObservable();
  }

  public init(): Observable<boolean> {
    if (this._ready) {
      return this.ready;
    }

    const s = new Subject<boolean>();

    console.log('not initialized yet, doing now');
    this.backend.isReady.filter(v => v === true).take(1).concatMap(() =>
      this.backend.send(messaging.REQUEST_OPEN_DATABASE, null).map(() => true)
    ).subscribe(
      () => {
        this._ready = true;
        this._dbopen$.next(this._ready);
        this.store.dispatch(new BackendConnActions.Set(true));
        s.next(this._ready);
      },
      (e: Error) => s.error(e),
      () => s.complete()
    );

    return s.asObservable();
  }

  public loadGroups(): Observable<Group[]> {
    const s = new Subject<Group[]>();

    this.backend.send(messaging.REQUEST_GROUPS, null).subscribe(
      (m: messaging.Message) => {
        const groups: Group[] = [];
        m.data.forEach(v => groups.push(Object.assign(new Group, v)));
        console.log('loaded groups:', groups);
        this.store.dispatch(new GroupsActions.Set(groups));
        s.next(groups);
      },
      (e: Error) => s.error(e),
      () => s.complete()
    );

    return s.asObservable();
  }

  public addGroup(group: Group): Observable<Group> {
    const s = new Subject<Group>();

    this.backend.send(messaging.REQUEST_ADD_GROUP, group).subscribe(
      (m: messaging.Message) => {
        const newGroup: Group = Object.assign(new Group, m.data);
        this.store.dispatch(new GroupsActions.Add(newGroup));
        s.next(newGroup);
      },
      (e: Error) => s.error(e),
      () => s.complete()
    );

    return s.asObservable();
  }

  public saveGroup(group: Group): Observable<Group> {
    const s = new Subject<Group>();

    this.backend.send(messaging.REQUEST_UPDATE_GROUP, group).subscribe(
      (m: messaging.Message) => {
        const newGroup: Group = Object.assign(new Group, m.data);
        this.store.dispatch(new GroupsActions.Update(newGroup));
        s.next(newGroup);
      },
      (e: Error) => s.error(e),
      () => s.complete()
    );

    return s.asObservable();
  }

  public selectGroup(tgt: Group): Observable<boolean> {
    const s = new Subject<boolean>();

    this.store.select('groups').take(1).subscribe((groups: Group[]) => {
      const grp = groups.find((v) => v.id === tgt.id);
      if (!grp) {
        s.error(new Error('group not found'));
        return;
      }

      console.log('loading tasks for group', grp);
      // Load tasks for the group
      this.backend.send(messaging.REQUEST_GROUP_TASKS, grp).subscribe(
        (m: messaging.Message) => {
          console.log('got tasks', m.data);
          const tasks: Task[] = [];
          m.data.forEach(t => tasks.push(Object.assign(new Task, t)));
          this.store.dispatch(new SelectedGroupActions.Set(grp));
          this.store.dispatch(new GroupTasksActions.Set(tasks));
          s.next(true);
        },
        (e: Error) => s.error(e),
        () => s.complete()
      );
    });

    return s.asObservable();
  }

  public loadTaskHistory(): Observable<Task[]> {
    const s = new Subject<Group>();

    this.backend.send(messaging.REQUEST_GET_TASK_HISTORY, null).subscribe(
      (m: messaging.Message) => {
        console.log('loaded task history:', m.data);
        const tasks: Task[] = [];
        m.data.forEach(t => tasks.push(Object.assign(new Task, t)));

        this.store.dispatch(new TaskHistoryActions.Set(tasks));
        s.next(tasks);
      },
      (e: Error) => s.error(e),
      () => s.complete()
    );

    return s.asObservable();
  }

  public saveTaskHistory(history: Task[]): Observable<boolean> {
    const s = new Subject<Group>();

    console.log('saving history:', history);
    this.backend.send(messaging.REQUEST_SET_TASK_HISTORY, {history}).subscribe(
      (m: messaging.Message) => {
        console.log('updated task history');
        s.next(true);
      },
      (e: Error) => s.error(e),
      () => s.complete()
    );

    return s.asObservable();
  }

  public loadActiveTask(): Observable<Task> {
    const s = new Subject<Group>();

    this.backend.send(messaging.REQUEST_ACTIVE_TASK, null).subscribe(
      (m: messaging.Message) => {
        console.log('loaded active task:', m.data);
        let task: Task = null;

        if (m.data !== null) {
          task = Object.assign(new Task, m.data);
        }

        this.store.dispatch(new ActiveTaskActions.Set(task));
        s.next(task);
      },
      (e: Error) => s.error(e),
      () => s.complete()
    );

    return s.asObservable();
  }

  public addTask(task: Task): Observable<Task> {
    const s = new Subject<Task>();

    this.backend.send(messaging.REQUEST_ADD_TASK, task).subscribe(
      (m: messaging.Message) => {
        const newTask: Task = Object.assign(new Task, m.data);

        // If selected group is the tasks group, then add.
        if (this._selectedGroup && this._selectedGroup.id === newTask.groupid) {
          this.store.dispatch(new GroupTasksActions.Add(newTask));
        }

        s.next(newTask);
      },
      (e: Error) => s.error(e),
      () => s.complete()
    );

    return s.asObservable();
  }

  public saveTask(task: Task): Observable<Task> {
    const s = new Subject<Task>();

    this.backend.send(messaging.REQUEST_UPDATE_TASK, task).subscribe(
      (m: messaging.Message) => {
        const newTask: Task = Object.assign(new Task, m.data);
        this.store.dispatch(new GroupTasksActions.Update(newTask));
        s.next(newTask);
      },
      (e: Error) => s.error(e),
      () => s.complete()
    );

    return s.asObservable();
  }

  public startTask(task: Task): Observable<Task> {
    const s = new Subject<Task>();

    console.log('starting task', task);
    // If we have an active task right now, we stop it first.
    if (this._activeTask) {
      this.stopTask(this._activeTask).subscribe(
        () => {
          // Wait until activeTask isn't set.
          this.store.select('activeTask').filter(v => v === null).take(1).subscribe(() => {
            this.startTask(task).subscribe(
              (t: Task) => s.next(t),
              (e: Error) => s.error(e),
              () => s.complete()
            );
          });
        },
        (e: Error) => s.error(e)
      );
    } else {
      this.backend.send(messaging.REQUEST_START_TASK, task).subscribe(
        (m: messaging.Message) => {
          const newTask: Task = Object.assign(new Task, m.data);
          this.store.dispatch(new GroupTasksActions.Update(newTask));
          this.store.dispatch(new ActiveTaskActions.Set(newTask));
          this.store.dispatch(new TaskHistoryActions.Remove(newTask));
          s.next(newTask);
        },
        (e: Error) => s.error(e),
        () => s.complete()
      );
    }

    return s.asObservable();
  }

  public stopTask(task: Task): Observable<Task> {
    const s = new Subject<Task>();

    console.log('stopping task', task);
    this.backend.send(messaging.REQUEST_STOP_TASK, task).subscribe(
      (m: messaging.Message) => {
        const newTask: Task = Object.assign(new Task, m.data);
        if (
          this._activeTask &&
          newTask.id === this._activeTask.id &&
          newTask.groupid === this._activeTask.groupid
        ) {
          this.store.dispatch(new ActiveTaskActions.Clear());
        }

        this.store.dispatch(new GroupTasksActions.Update(newTask));
        this.store.dispatch(new TaskHistoryActions.Add(newTask));
        s.next(newTask);
      },
      (e: Error) => s.error(e),
      () => s.complete()
    );

    return s.asObservable();
  }

  public selectTask(tgt: Task): Observable<Task> {
    const s = new ReplaySubject<Task>();

    this.backend.send(messaging.REQUEST_GET_TASK, tgt).subscribe(
      (m: messaging.Message) => {
        const task: Task = Object.assign(new Task, m.data);
        this.store.dispatch(new SelectedTaskActions.Set(task));
        s.next(task);
      },
      e => s.error(e),
      () => s.complete()
    );

    return s.asObservable();
  }
}
