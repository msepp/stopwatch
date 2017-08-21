import { Component, OnInit, OnDestroy } from '@angular/core';
import { Location } from '@angular/common';
import { Router, RouterStateSnapshot } from '@angular/router';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Observable } from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import { Store } from '@ngrx/store';
import { StopwatchService } from '../../services/stopwatch.service';
import { ErrorService } from '../../services/error.service';
import { AppState, Group } from '../../model';

@Component({
  selector: 'app-group-usage',
  templateUrl: './group-usage.component.html',
  styleUrls: ['./group-usage.component.less']
})
export class GroupUsageComponent implements OnInit, OnDestroy {
  public usageForm: FormGroup;
  public group$: Subscription;
  public minDate = new Date(2017, 7, 1);
  public group: Group;

  constructor(
    private store: Store<AppState>,
    private location: Location,
    private router: Router,
    private stopwatch: StopwatchService,
    private fb: FormBuilder,
    private err: ErrorService
  ) { }

  ngOnInit() {
    const snapshot: RouterStateSnapshot = this.router.routerState.snapshot;
    console.log(snapshot);

    this.usageForm = this.fb.group({
      id: [snapshot.root.params['id']],
      start: [snapshot.root.queryParams['start'], Validators.required],
      end: [snapshot.root.queryParams['end'], Validators.required]
    });

    this.group$ = this.store.select('selectedGroup').subscribe((g: Group) => {
      this.group = g;
      this.usageForm.get('id').setValue(g.id);
    });
  }

  ngOnDestroy() {
    this.group$.unsubscribe();
  }

  public getUsage() {
    if (this.usageForm.valid) {
    }
  }

  public goBack() {
    this.location.back();
  }
}
