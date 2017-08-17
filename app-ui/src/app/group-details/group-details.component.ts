import { Component, OnInit } from '@angular/core';
import { Location } from '@angular/common';
import { Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Store } from '@ngrx/store';
import { StopwatchService } from '../services/stopwatch.service';
import { AppState, Group, Task } from '../model';

@Component({
  selector: 'app-group-details',
  templateUrl: './group-details.component.html',
  styleUrls: ['./group-details.component.less']
})
export class GroupDetailsComponent implements OnInit {
  public tasks$: Store<Task[]>;

  constructor(
    private store: Store<AppState>,
    private router: Router,
    private location: Location,
    private stopwatch: StopwatchService
  ) { }

  ngOnInit() {
    this.tasks$ = this.store.select('groupTasks');
  }

  public goBack() {
    this.location.back();
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
      e => console.log('error starting:', e)
    );
  }
}
