import { Component, OnInit } from '@angular/core';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs/Observable';
import { AppState, Task } from '../../model';

@Component({
  selector: 'app-task-history',
  templateUrl: './task-history.component.html',
  styleUrls: ['./task-history.component.less']
})
export class TaskHistoryComponent implements OnInit {
  public tasks: Observable<Task[]>;

  constructor(
    private store: Store<AppState>
  ) {
    this.tasks = this.store.select('taskHistory').map((tasks: Task[]) => {
      return tasks.slice(0, 5);
    });
  }

  ngOnInit() {
  }
}
