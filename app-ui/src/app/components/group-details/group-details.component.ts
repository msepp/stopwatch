import { Component, OnInit, OnDestroy } from '@angular/core';
import { Location } from '@angular/common';
import { Router } from '@angular/router';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Observable } from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import { Store } from '@ngrx/store';
import { StopwatchService } from '../../services/stopwatch.service';
import { ErrorService } from '../../services/error.service';
import { AppState, Group, Task } from '../../model';

@Component({
  selector: 'app-group-details',
  templateUrl: './group-details.component.html',
  styleUrls: ['./group-details.component.less']
})
export class GroupDetailsComponent implements OnInit, OnDestroy {
  public tasks$: Store<Task[]>;
  public groupForm: FormGroup;
  public usageForm: FormGroup;
  public group$: Subscription;
  public minDate = new Date(2017, 7, 1);

  constructor(
    private store: Store<AppState>,
    private location: Location,
    private router: Router,
    private stopwatch: StopwatchService,
    private fb: FormBuilder,
    private err: ErrorService
  ) { }

  ngOnInit() {
    this.groupForm = this.fb.group({
      id: [0],
      name: ['', Validators.required]
    });

    const end: Date = new Date();
    const start = new Date(end.getTime());
    start.setDate(start.getDate() - 5);

    this.usageForm = this.fb.group({
      id: [0],
      start: [start, Validators.required],
      end: [end, Validators.required]
    });

    this.group$ = this.store.select('selectedGroup').subscribe((g: Group) => {
      this.groupForm.reset({
        id: g.id,
        name: g.name
      });

      this.usageForm.get('id').setValue(g.id);
    });

    this.tasks$ = this.store.select('groupTasks');
  }

  ngOnDestroy() {
    this.group$.unsubscribe();
  }

  public goBack() {
    this.location.back();
  }

  public getUsage() {
    if (this.usageForm.valid) {
      this.router.navigate(['/usage', this.usageForm.get('id').value], {
        queryParams: {
          start: this.usageForm.get('start').value,
          end: this.usageForm.get('end').value
        }
      });
    }
  }

  public saveGroup() {
    if (this.groupForm.valid) {
      this.stopwatch.saveGroup(this.groupForm.value).subscribe(
        () => {},
        (e: Error) => this.err.log(e)
      );
    }
  }
}
