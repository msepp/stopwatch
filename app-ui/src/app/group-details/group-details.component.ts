import { Component, OnInit, OnDestroy } from '@angular/core';
import { Location } from '@angular/common';
import { Router } from '@angular/router';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Observable } from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import { Store } from '@ngrx/store';
import { StopwatchService } from '../services/stopwatch.service';
import { AppState, Group, Task } from '../model';

@Component({
  selector: 'app-group-details',
  templateUrl: './group-details.component.html',
  styleUrls: ['./group-details.component.less']
})
export class GroupDetailsComponent implements OnInit, OnDestroy {
  public tasks$: Store<Task[]>;
  public groupForm: FormGroup;
  public group$: Subscription;

  constructor(
    private store: Store<AppState>,
    private router: Router,
    private location: Location,
    private stopwatch: StopwatchService,
    private fb: FormBuilder
  ) { }

  ngOnInit() {
    this.groupForm = this.fb.group({
      id: [0],
      name: ['', Validators.required]
    });

    this.group$ = this.store.select('selectedGroup').subscribe((g: Group) => {
      this.groupForm.reset({
        id: g.id,
        name: g.name
      });
    });

    this.tasks$ = this.store.select('groupTasks');
  }

  ngOnDestroy() {
    this.group$.unsubscribe();
  }

  public goBack() {
    this.location.back();
  }

  public saveGroup() {
    if (this.groupForm.valid) {
      this.stopwatch.saveGroup(this.groupForm.value).subscribe(
        () => {},
        e => console.log('error saving', e)
      );
    }
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
