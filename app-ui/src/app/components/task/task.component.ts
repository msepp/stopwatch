import { Component, OnInit, Input } from '@angular/core';
import { Router } from '@angular/router';
import { StopwatchService } from '../../services/stopwatch.service';
import { Task } from '../../model';

@Component({
  selector: 'app-task',
  templateUrl: './task.component.html',
  styleUrls: ['./task.component.less']
})
export class TaskComponent implements OnInit {
  @Input() task: Task;

  constructor(
    private stopwatch: StopwatchService,
    private router: Router
  ) { }

  ngOnInit() {
  }

  public startTask(t: Task) {
    this.stopwatch.startTask(t).subscribe(
      () => {},
      e => console.log('error starting:', e)
    );
  }

  public stopTask(t: Task) {
    this.stopwatch.stopTask(t).subscribe(
      () => {},
      e => console.log('error stopping:', e)
    );
  }

  public editTask(t: Task) {
    console.log('selecting task', t);
    this.stopwatch.selectTask(t).subscribe(
      () => {
        console.log('navigating...');
        this.router.navigate(['/task', t.groupid + '-' + t.id]);
      },
      e => console.log('error selecting task:', e),
      () => console.log('select done')
    );
  }
}
