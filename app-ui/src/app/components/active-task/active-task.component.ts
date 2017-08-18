import { Component, OnInit, OnDestroy, Input } from '@angular/core';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs/Observable';
import { StopwatchService } from '../../services/stopwatch.service';
import { ErrorService } from '../../services/error.service';
import { Task, AppState } from '../../model';

@Component({
  selector: 'app-active-task',
  templateUrl: './active-task.component.html',
  styleUrls: ['./active-task.component.less']
})
export class ActiveTaskComponent implements OnInit, OnDestroy {
  public active: Task = null;
  public elapsed = '-';
  private timer;
  private task$;

  constructor(
    private store: Store<AppState>,
    private stopwatch: StopwatchService,
    private err: ErrorService
  ) {
    this.task$ = this.store.select('activeTask').subscribe(
      (t: Task) => {
        this.active = t;
        console.log('active task now', t);
      }
    );

    this.timer = setInterval(() => {
      if (this.active) {
        const since = new Date(this.active.running);
        let seconds = Math.floor(((new Date()).getTime() - since.getTime()) / 1000);
        let hours = 0;
        let minutes = 0;

        const str = '';
        if (seconds > 3600) {
          hours = Math.floor(seconds / 3600);
          seconds = seconds - (hours * 3600);
        }

        if (seconds > 60) {
          minutes = Math.floor(seconds / 60);
          seconds = seconds - (minutes * 60);
        }

        this.elapsed = this.zeroPadded(hours) + ':' + this.zeroPadded(minutes) + ':' + this.zeroPadded(seconds);
      }
    }, 1000);
  }

  private zeroPadded(n: number): string {
    if (n < 10) {
      return '0' + n;
    }
    return '' + n;
  }

  ngOnInit() {
  }

  ngOnDestroy() {
    this.task$.unsubscribe();
    clearInterval(this.timer);
  }

  stop() {
    this.stopwatch.stopTask(this.active).subscribe(
      () => {},
      (e: Error) => this.err.log(e)
    );
  }
}
