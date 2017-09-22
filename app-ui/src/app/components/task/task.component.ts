import { Component, OnInit, Input } from '@angular/core';
import { Router } from '@angular/router';
import { StopwatchService } from '../../services/stopwatch.service';
import { ErrorService } from '../../services/error.service';
import { Task } from '../../model';

@Component({
  selector: 'app-task, [app-task]',
  templateUrl: './task.component.html',
  styleUrls: ['./task.component.less']
})
export class TaskComponent implements OnInit {
  @Input() task: Task;

  constructor(
    private stopwatch: StopwatchService,
    private router: Router,
    private err: ErrorService
  ) { }

  ngOnInit() {
  }

  public startTask(t: Task) {
    this.stopwatch.startTask(t).subscribe(
      () => {},
      (e: Error) => this.err.log(e)
    );
  }

  public stopTask(t: Task) {
    this.stopwatch.stopTask(t).subscribe(
      () => {},
      (e: Error) => this.err.log(e)
    );
  }

  public editTask(t: Task) {
    console.log('selecting task', t);
    this.stopwatch.selectTask(t).subscribe(
      () => this.router.navigate(['/task', t.groupid + '-' + t.id]),
      (e: Error) => this.err.log(e)
    );
  }
}
