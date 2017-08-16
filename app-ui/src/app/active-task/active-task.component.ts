import { Component, OnInit, OnDestroy, Input } from '@angular/core';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs/Observable';
import { Task, AppState } from '../model';
@Component({
  selector: 'app-active-task',
  templateUrl: './active-task.component.html',
  styleUrls: ['./active-task.component.less']
})
export class ActiveTaskComponent implements OnInit, OnDestroy {
  public active: Task = null;
  private task$;

  constructor(private store: Store<AppState>) {
    this.task$ = this.store.select('activeTask').subscribe(
      (t: Task) => {
        this.active = t;
        console.log('active task now', t);
      }
    );
  }

  ngOnInit() {
  }

  ngOnDestroy() {
    this.task$.unsubscribe();
  }

}
