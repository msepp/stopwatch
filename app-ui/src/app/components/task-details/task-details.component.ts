import { Component, OnInit, OnDestroy } from '@angular/core';
import { Location } from '@angular/common';
import { Router } from '@angular/router';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Observable } from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import { Store } from '@ngrx/store';
import { StopwatchService } from '../../services/stopwatch.service';
import { ErrorService } from '../../services/error.service';
import { AppState, Group, Task } from '../../model';

@Component({
  selector: 'app-task-details',
  templateUrl: './task-details.component.html',
  styleUrls: ['./task-details.component.less']
})
export class TaskDetailsComponent implements OnInit, OnDestroy {
  private task$: Subscription;
  public taskFG: FormGroup;

  constructor(
    private store: Store<AppState>,
    private router: Router,
    private location: Location,
    private stopwatch: StopwatchService,
    private fb: FormBuilder,
    private err: ErrorService
  ) { }

  ngOnInit() {
    this.taskFG = this.fb.group({
      id: [0],
      groupid: [0],
      name: ['', Validators.required],
      costcode: ['', Validators.required]
    });

    this.task$ = this.store.select('selectedTask').subscribe((t: Task) => {
      this.taskFG.reset({
        id: t.id,
        groupid: t.groupid,
        name: t.name,
        costcode: t.costcode
      });
    });
  }

  ngOnDestroy() {
    this.task$.unsubscribe();
  }

  public goBack() {
    this.location.back();
  }

  public saveTask() {
    if (this.taskFG.valid) {
      this.stopwatch.saveTask(this.taskFG.value).subscribe(
        () => {},
        (e: Error) => this.err.log(e)
      );
    }
  }
}
