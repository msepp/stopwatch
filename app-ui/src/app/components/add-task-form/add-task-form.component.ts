import { Component, ViewChild, OnInit, OnDestroy} from '@angular/core';
import { FormBuilder, FormGroup, Validators, FormGroupDirective } from '@angular/forms';
import { Store } from '@ngrx/store';
import { AppState, Group, Task } from '../../model';
import { StopwatchService } from '../../services/stopwatch.service';
import { ErrorService } from '../../services/error.service';

@Component({
  selector: 'app-add-task-form',
  templateUrl: './add-task-form.component.html',
  styleUrls: ['./add-task-form.component.less']
})
export class AddTaskFormComponent implements OnInit, OnDestroy {
  private group$;
  public group: Group;
  public taskFG: FormGroup;
  @ViewChild(FormGroupDirective) taskForm;

  constructor(
    private stopwatch: StopwatchService,
    private fb: FormBuilder,
    private store: Store<AppState>,
    private err: ErrorService
  ) {
    this.group$ = this.store.select('selectedGroup').subscribe(g => this.group = g);
  }

  ngOnInit() {
    this.taskFG = this.fb.group({
      name: ['', Validators.required],
      costcode: ['', Validators.required]
    });
  }

  ngOnDestroy() {
    this.group$.unsubscribe();
  }

  add() {
    if (this.taskFG.valid) {
      const t: Task = {
        groupid: this.group.id,
        name: this.taskFG.get('name').value,
        costcode: this.taskFG.get('costcode').value
      };

      this.stopwatch.addTask(t).subscribe(
        () => {
          this.taskFG.reset();
          this.taskForm.resetForm();
        },
        (e: Error) => this.err.log(e)
      );
    }
  }
}
